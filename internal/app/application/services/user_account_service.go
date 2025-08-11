package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/renderview-inc/backend/internal/app/application/dtos"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
)

var (
	ErrNoAccountFound     = errors.New("account was not found")
	ErrInvalidCredentials = errors.New("invalid credentials: make sure you pass correct email/phone/tag")
)

type UserAccountRepository interface {
	Create(ctx context.Context, uacc *entities.UserAccount) error
	ReadById(ctx context.Context, accID uuid.UUID) (*entities.UserAccount, error)
	ReadByTag(ctx context.Context, accTag string) (*entities.UserAccount, error)
	ReadByEmail(ctx context.Context, email string) (*entities.UserAccount, error)
	ReadByPhone(ctx context.Context, phone string) (*entities.UserAccount, error)
	Update(ctx context.Context, uacc *entities.UserAccount) error
	Delete(ctx context.Context, accID uuid.UUID) error
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, passwordHash string) bool
}

type UserAccountService struct {
	accountRepository UserAccountRepository
	passwordHasher    PasswordHasher
}

func NewUserAccountService(accountRepository UserAccountRepository, passwordHasher PasswordHasher) UserAccountService {
	return UserAccountService{
		accountRepository: accountRepository,
		passwordHasher:    passwordHasher,
	}
}

func (uas *UserAccountService) Register(ctx context.Context, uacc *entities.UserAccount) error {
	existingAcc, err := uas.accountRepository.ReadByTag(ctx, uacc.GetTag())
	if err != nil {
		return fmt.Errorf("failed to check account existance: %w", err)
	}

	if existingAcc != nil {
		return errors.New("account with this tag already exists")
	}

	err = uas.accountRepository.Create(ctx, uacc)
	if err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	return nil
}

func (uas *UserAccountService) VerifyCredentials(ctx context.Context, credentials dtos.Credentials) (*uuid.UUID, error) {
	var acc *entities.UserAccount
	var err error

	if credentials.Email != "" {
		acc, err = uas.accountRepository.ReadByEmail(ctx, credentials.Email)
		if err != nil {
			return nil, fmt.Errorf("read by email: %w", err)
		}
	}
	if acc == nil && credentials.Phone != "" {
		acc, err = uas.accountRepository.ReadByPhone(ctx, credentials.Phone)
		if err != nil {
			return nil, fmt.Errorf("read by phone: %w", err)
		}
	}
	if acc == nil && credentials.Tag != "" {
		acc, err = uas.accountRepository.ReadByTag(ctx, credentials.Tag)
		if err != nil {
			return nil, fmt.Errorf("read by tag: %w", err)
		}
	}

	if acc == nil {
		return nil, ErrNoAccountFound
	}

	if !uas.passwordHasher.VerifyPassword(credentials.Password, acc.GetPasswordHash()) {
		return nil, nil
	}

	id := acc.GetID()
	return &id, nil
}
