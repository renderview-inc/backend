package repositories

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
)

type UserAccountRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewUserAccountRepository(pool *pgxpool.Pool) *UserAccountRepository {
	return &UserAccountRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (uar *UserAccountRepository) Create(ctx context.Context, uacc *entities.UserAccount) error {
	sql, args, err := uar.builder.Insert("user_accounts").
		Columns("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
		Values(uacc.Id, uacc.Tag, uacc.Name, uacc.Desc, uacc.PasswordHash, uacc.Email, uacc.Phone).
		ToSql()

	if err != nil {
		return err
	}

	_, err = uar.pool.Exec(ctx, sql, args...)
	log.Printf("UserAccount created: %v", uacc.Id)
	return err
}

func (uar *UserAccountRepository) ReadById(ctx context.Context, accID uuid.UUID) (*entities.UserAccount, error) {
	sql, args, err := uar.builder.Select("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
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
	sql, args, err := uar.builder.Select("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
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
	sql, args, err := uar.builder.Select("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
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
	sql, args, err := uar.builder.Select("id", "tag", "name", "\"desc\"", "password_hash", "email", "phone").
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
		uar.builder.Update("user_accounts").
			Set("tag", uacc.Tag).
			Set("name", uacc.Name).
			Set("\"desc\"", uacc.Desc).
			Set("password_hash", uacc.PasswordHash).
			Set("email", uacc.Email).
			Set("phone", uacc.Phone).
			Where(sq.Eq{"id": uacc.Id}).
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
	sql, args, err := uar.builder.Delete("user_accounts").Where(sq.Eq{"id": accID}).ToSql()

	if err != nil {
		return err
	}

	_, err = uar.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
