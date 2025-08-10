package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/renderview-inc/backend/internal/app/application/dtos"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
	"github.com/renderview-inc/backend/internal/pkg/txhelper"
)

var (
	ErrAccessTokenInvalid  = errors.New("access token expired")
	ErrInvalidSessionID    = errors.New("invalid session id")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrSessionExpired      = errors.New("session expired")
)

type LoginHistoryRepository interface {
	Create(ctx context.Context, tx pgx.Tx, loginInfo entities.LoginInfo) error
	ReadById(ctx context.Context, id uuid.UUID) (*entities.LoginInfo, error)
	Update(ctx context.Context, loginInfo entities.LoginInfo) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserSessionRepository interface {
	Create(ctx context.Context, tx pgx.Tx, session entities.UserSession) error
	CreateStandalone(ctx context.Context, session entities.UserSession) error
	ReadById(ctx context.Context, id uuid.UUID) (*entities.UserSession, error)
	Update(ctx context.Context, session entities.UserSession) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserSessionCache interface {
	SaveToken(ctx context.Context, token string, ttl time.Duration) error
	CheckToken(ctx context.Context, token string) (bool, error)
	RevokeToken(ctx context.Context, token string) error
}

type TokenIssuer interface {
	IssueAccessToken() (string, time.Duration, error)
	IssueRefreshToken(sessionID string) (string, time.Duration, error)
}

type TokenHasher interface {
	HashToken(token string) (string, error)
}

type AuthService struct {
	loginHistoryRepository LoginHistoryRepository
	sessionRepository      UserSessionRepository
	sessionCache           UserSessionCache
	accountService         UserAccountService
	txHelper               *txhelper.TxHelper
	tokenIssuer            TokenIssuer
	tokenHasher            TokenHasher
}

func NewAuthService(loginHistoryRepository LoginHistoryRepository,
	sessionRepository UserSessionRepository, sessionCache UserSessionCache,
	accountService UserAccountService, txHelper *txhelper.TxHelper, tokenIssuer TokenIssuer,
	tokenHasher TokenHasher) *AuthService {
	return &AuthService{
		loginHistoryRepository, sessionRepository, sessionCache, accountService, txHelper, tokenIssuer,
		tokenHasher,
	}
}

func (as *AuthService) Login(ctx context.Context, loginDto dtos.LoginDto) (dtos.TokensDto, error) {
	userID, err := as.accountService.VerifyCredentials(ctx, loginDto.Credentials)

	if err != nil {
		return dtos.TokensDto{}, fmt.Errorf("verify credentials: %w", err)
	}
	if userID == nil {
		return dtos.TokensDto{}, ErrInvalidCredentials
	}
	sessionID := uuid.New()

	accessToken, accessLifeTime,
		refreshToken, refreshLifeTime, err := as.issueTokens(sessionID)
	if err != nil {
		return dtos.TokensDto{}, fmt.Errorf("issue tokens: %w", err)
	}

	refreshTokenHash, err := as.tokenHasher.HashToken(refreshToken)
	if err != nil {
		return dtos.TokensDto{}, fmt.Errorf("hash refresh token: %w", err)
	}

	now := time.Now()

	userSession := entities.NewUserSession(
		sessionID,
		*userID,
		refreshTokenHash,
		now,
		now,
		now.Add(refreshLifeTime),
		now,
		false,
		nil,
	)
	userLogin := entities.NewLoginInfo(
		uuid.New(),
		*userID,
		time.Now(),
		loginDto.LoginMeta.UserAgent,
		loginDto.LoginMeta.IpAddr,
		true,
	)

	repoTx, err := as.txHelper.Begin(ctx)
	if err != nil {
		return dtos.TokensDto{}, fmt.Errorf("save session and login info: %w", err)
	}

	if err := as.sessionRepository.Create(ctx, repoTx, userSession); err != nil {
		if err := repoTx.Rollback(ctx); err != nil {
			return dtos.TokensDto{}, fmt.Errorf("rollback save session db transaction: %w", err)
		}
		return dtos.TokensDto{}, fmt.Errorf("save session: %w", err)
	}
	if err := as.loginHistoryRepository.Create(ctx, repoTx, userLogin); err != nil {
		if err := repoTx.Rollback(ctx); err != nil {
			return dtos.TokensDto{}, fmt.Errorf("rollback save login info db transaction: %w", err)
		}
		return dtos.TokensDto{}, fmt.Errorf("save login info: %w", err)
	}

	if err = repoTx.Commit(ctx); err != nil {
		return dtos.TokensDto{}, fmt.Errorf("commit transaction to save session and login info: %w", err)
	}

	if err = as.sessionCache.SaveToken(ctx, accessToken, accessLifeTime); err != nil {
		return dtos.TokensDto{}, fmt.Errorf("save access token to cache: %w", err)
	}
	if err = as.sessionCache.SaveToken(ctx, refreshToken, refreshLifeTime); err != nil {
		// TODO just log that cache failed
	}

	return dtos.TokensDto{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (as *AuthService) Authorize(ctx context.Context, accessToken string) error {
	ok, err := as.sessionCache.CheckToken(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("check access token in cache: %w", err)
	}

	if !ok {
		return ErrAccessTokenInvalid
	}

	return nil
}

func (as *AuthService) Refresh(ctx context.Context, refreshToken string) (dtos.TokensDto, error) {
	var session *entities.UserSession
	tokenInfo := strings.Split(refreshToken, ".")
	sessionID, err := uuid.Parse(tokenInfo[0])

	if err != nil {
		return dtos.TokensDto{}, ErrInvalidSessionID
	}

	refreshTokenHash, err := as.tokenHasher.HashToken(refreshToken)

	if err != nil {
		return dtos.TokensDto{}, fmt.Errorf("hash refresh token: %w", err)
	}

	ok, err := as.sessionCache.CheckToken(ctx, refreshToken)
	if err != nil {
		return dtos.TokensDto{}, fmt.Errorf("check refresh token in cache: %w", err)
	}

	if !ok {
		session, err = as.sessionRepository.ReadById(ctx, sessionID)
		if err != nil {
			return dtos.TokensDto{}, fmt.Errorf("read session by ID: %w", err)
		}
		if session == nil {
			return dtos.TokensDto{}, ErrInvalidSessionID
		}

		if session.RefreshTokenHash() != refreshTokenHash {
			return dtos.TokensDto{}, ErrInvalidRefreshToken
		}
		if session.RefreshExpiresAt().Before(time.Now()) {
			return dtos.TokensDto{}, ErrSessionExpired
		}
	}

	newAccessToken, newAccessLifeTime,
		newRefreshToken, newRefreshLifeTime, err := as.issueTokens(session.ID())
	if err != nil {
		return dtos.TokensDto{}, fmt.Errorf("issue new tokens: %w", err)
	}

	newRefreshTokenHash, err := as.tokenHasher.HashToken(newRefreshToken)
	if err != nil {
		return dtos.TokensDto{}, fmt.Errorf("hash new refresh token: %w", err)
	}

	now := time.Now()

	newSession := entities.NewUserSession(
		uuid.New(),
		session.UserID(),
		newRefreshTokenHash,
		session.CreatedAt(),
		now,
		now.Add(newRefreshLifeTime),
		now,
		false,
		&sessionID,
	)
	updatedOldSession := entities.NewUserSession(
		sessionID,
		session.UserID(),
		session.RefreshTokenHash(),
		session.CreatedAt(),
		now,
		session.RefreshExpiresAt(),
		session.LastUsedAt(),
		true,
		session.RotatedFromSessionID(),
	)

	tx, err := as.txHelper.Begin(ctx)
	if err != nil {
		return dtos.TokensDto{}, err
	}

	if err = as.sessionRepository.Create(ctx, tx, newSession); err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return dtos.TokensDto{}, fmt.Errorf("persist new session: %w", err)
		}
	}
	if err = as.sessionRepository.Update(ctx, updatedOldSession); err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return dtos.TokensDto{}, fmt.Errorf("update old session: %w", err)
		}
	}

	if err = as.sessionCache.SaveToken(ctx, newAccessToken, newAccessLifeTime); err != nil {
		return dtos.TokensDto{}, fmt.Errorf("new access token cache: %w", err)
	}
	if err = as.sessionCache.SaveToken(ctx, newRefreshToken, newRefreshLifeTime); err != nil {
		// TODO just log
	}

	if err = tx.Commit(ctx); err != nil {
		return dtos.TokensDto{}, fmt.Errorf("commit transaction to save new session and update old one: %w", err)
	}

	return dtos.TokensDto{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
}

func (as *AuthService) Logout(ctx context.Context, tokens dtos.TokensDto) error {
	if err := as.sessionCache.RevokeToken(ctx, tokens.AccessToken); err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}

	refreshTokenInfo := strings.Split(tokens.RefreshToken, ".")
	sessionID, err := uuid.Parse(refreshTokenInfo[0])
	if err != nil {
		return fmt.Errorf("parse refresh token secret: %w", err)
	}

	session, err := as.sessionRepository.ReadById(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("read session to revoke: %w", err)
	}
	if session == nil {
		return ErrSessionExpired
	}

	revokedSession := entities.NewUserSession(
		sessionID,
		session.UserID(),
		session.RefreshTokenHash(),
		session.CreatedAt(),
		time.Now(),
		session.RefreshExpiresAt(),
		session.LastUsedAt(),
		true,
		session.RotatedFromSessionID(),
	)

	if err := as.sessionRepository.Update(ctx, revokedSession); err != nil {
		return fmt.Errorf("save revoked session: %w", err)
	}

	return nil
}

func (as *AuthService) issueTokens(sessionID uuid.UUID) (string, time.Duration, string, time.Duration, error) {
	accessToken, accessLifeTime, err := as.tokenIssuer.IssueAccessToken()
	if err != nil {
		return "", time.Duration(0), "", time.Duration(0), fmt.Errorf("issue access token: %w", err)
	}
	refreshToken, refreshLifeTime, err := as.tokenIssuer.IssueRefreshToken(sessionID.String())
	if err != nil {
		return "", time.Duration(0), "", time.Duration(0), fmt.Errorf("issue refresh token: %w", err)
	}

	return accessToken, accessLifeTime, refreshToken, refreshLifeTime, nil
}
