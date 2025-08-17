package repositories

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
)

type UserSessionRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewUserSessionRepository(pool *pgxpool.Pool) *UserSessionRepository {
	return &UserSessionRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (usr *UserSessionRepository) Create(ctx context.Context, tx pgx.Tx, session entities.UserSession) error {
	sql, args, err :=
		usr.builder.Insert("user_sessions").
			Columns("id", "user_id", "refresh_token_hash", "created_at", "updated_at", "refresh_expires_at", "last_used_at",
				"revoked", "rotated_from_session_id").
			Values(session.Id, session.UserID, session.RefreshTokenHash, session.CreatedAt, session.UpdatedAt,
				session.RefreshExpiresAt, session.LastUsedAt, session.Revoked, session.RotatedFromSessionID).
			ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (usr *UserSessionRepository) CreateStandalone(ctx context.Context, session entities.UserSession) error {
	sql, args, err :=
		usr.builder.Insert("user_sessions").
			Columns("id", "user_id", "refresh_token_hash", "created_at", "updated_at", "refresh_expires_at", "last_used_at",
				"revoked", "rotated_from_session_id").
			Values(session.Id, session.UserID, session.RefreshTokenHash, session.CreatedAt,
				session.RefreshExpiresAt, session.LastUsedAt, session.Revoked, session.RotatedFromSessionID).
			ToSql()
	if err != nil {
		return err
	}

	_, err = usr.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (usr *UserSessionRepository) ReadById(ctx context.Context, id uuid.UUID) (*entities.UserSession, error) {
	sql, args, err := usr.builder.Select("id", "user_id", "refresh_token_hash", "created_at", "updated_at", "refresh_expires_at", "last_used_at",
		"revoked", "rotated_from_session_id").
		From("user_sessions").Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return nil, err
	}

	var us entities.UserSession
	err = usr.pool.QueryRow(ctx, sql, args...).Scan(
		&us.Id,
		&us.UserID,
		&us.RefreshTokenHash,
		&us.CreatedAt,
		&us.UpdatedAt,
		&us.RefreshExpiresAt,
		&us.LastUsedAt,
		&us.Revoked,
		&us.RotatedFromSessionID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &us, nil
}

func (usr *UserSessionRepository) Update(ctx context.Context, session entities.UserSession) error {
	sql, args, err :=
		usr.builder.Update("user_sessions").
			Set("user_id", session.UserID).
			Set("refresh_token_hash", session.RefreshTokenHash).
			Set("created_at", session.CreatedAt).
			Set("updated_at", session.UpdatedAt).
			Set("refresh_expires_at", session.RefreshExpiresAt).
			Set("last_used_at", session.LastUsedAt).
			Set("revoked", session.Revoked).
			Set("rotated_from_session_id", session.RotatedFromSessionID).
			Where(sq.Eq{"id": session.Id}).
			ToSql()

	if err != nil {
		return err
	}

	_, err = usr.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (usr *UserSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := usr.builder.Delete("user_sessions").Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return err
	}

	_, err = usr.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
