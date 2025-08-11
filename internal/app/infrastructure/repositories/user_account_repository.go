package repositories

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
	postgres "github.com/renderview-inc/backend/pkg/connections"
)

type UserAccountRepository struct {
	pool *pgxpool.Pool
}

func NewUserAccountRepository(pool *pgxpool.Pool) *UserAccountRepository {
	return &UserAccountRepository{pool}
}

func (uar *UserAccountRepository) Create(ctx context.Context, uacc *entities.UserAccount) error {
	sql, args, err := postgres.Psql.Insert("user_accounts").
		Columns("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
		Values(uacc.GetID(), uacc.GetTag(), uacc.GetName(), uacc.GetDesc(), uacc.GetPasswordHash(), uacc.GetEmail(), uacc.GetPhone()).
		ToSql()

	if err != nil {
		return err
	}

	_, err = uar.pool.Exec(ctx, sql, args...)
	log.Printf("UserAccount created: %v", uacc.GetID())
	return err
}

func (uar *UserAccountRepository) ReadById(ctx context.Context, accID uuid.UUID) (*entities.UserAccount, error) {
	sql, args, err := postgres.Psql.Select("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
		From("user_accounts").Where(sq.Eq{"id": accID}).ToSql()

	if err != nil {
		return nil, err
	}

	var uacc entities.UserAccount
	err = uar.pool.QueryRow(ctx, sql, args...).Scan(&uacc.Id, &uacc.Tag, &uacc.Name, &uacc.Desc, &uacc.PasswordHash, &uacc.Email, &uacc.Phone)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &uacc, nil
}

func (uar *UserAccountRepository) ReadByTag(ctx context.Context, accTag string) (*entities.UserAccount, error) {
	sql, args, err := postgres.Psql.Select("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
		From("user_accounts").Where(sq.Eq{"tag": accTag}).ToSql()

	if err != nil {
		return nil, err
	}

	var uacc entities.UserAccount
	err = uar.pool.QueryRow(ctx, sql, args...).Scan(&uacc.Id, &uacc.Tag, &uacc.Name, &uacc.Desc, &uacc.PasswordHash, &uacc.Email, &uacc.Phone)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &uacc, nil
}

func (uar *UserAccountRepository) ReadByEmail(ctx context.Context, email string) (*entities.UserAccount, error) {
	sql, args, err := postgres.Psql.Select("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
		From("user_accounts").Where(sq.Eq{"email": email}).ToSql()

	if err != nil {
		return nil, err
	}

	var uacc entities.UserAccount
	err = uar.pool.QueryRow(ctx, sql, args...).Scan(&uacc.Id, &uacc.Tag, &uacc.Name, &uacc.Desc, &uacc.PasswordHash, &uacc.Email, &uacc.Phone)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &uacc, nil
}

func (uar *UserAccountRepository) ReadByPhone(ctx context.Context, phone string) (*entities.UserAccount, error) {
	sql, args, err := postgres.Psql.Select("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
		From("user_accounts").Where(sq.Eq{"phone": phone}).ToSql()

	if err != nil {
		return nil, err
	}

	var uacc entities.UserAccount
	err = uar.pool.QueryRow(ctx, sql, args...).Scan(&uacc.Id, &uacc.Tag, &uacc.Name, &uacc.Desc, &uacc.PasswordHash, &uacc.Email, &uacc.Phone)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &uacc, nil
}

func (uar *UserAccountRepository) Update(ctx context.Context, uacc *entities.UserAccount) error {
	sql, args, err :=
		postgres.Psql.Update("user_accounts").
			Set("tag", uacc.GetTag()).
			Set("name", uacc.GetName()).
			Set("\"desc\"", uacc.GetDesc()).
			Set("password_hash", uacc.GetPasswordHash()).
			Set("email", uacc.GetEmail()).
			Set("phone", uacc.GetPhone()).
			Where(sq.Eq{"id": uacc.GetID()}).
			ToSql()

	if err != nil {
		return err
	}

	_, err = uar.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (uar *UserAccountRepository) Delete(ctx context.Context, accID uuid.UUID) error {
	sql, args, err := postgres.Psql.Delete("user_accounts").Where(sq.Eq{"id": accID}).ToSql()

	if err != nil {
		return err
	}

	_, err = uar.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
