package repositories

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
	"github.com/renderview-inc/backend/pkg/postgres"
)

type LoginHistoryRepository struct {
	pool *pgxpool.Pool
}

func NewLoginHistoryRepository(pool *pgxpool.Pool) *LoginHistoryRepository {
	return &LoginHistoryRepository{pool}
}

func (lhr *LoginHistoryRepository) Create(ctx context.Context, tx pgx.Tx, loginInfo entities.LoginInfo) error {
	sql, args, err :=
		postgres.Psql.Insert("user_login_histories").
			Columns("login_id", "user_id", "login_time", "user_agent", "ip_address", "success").
			Values(loginInfo.ID(), loginInfo.UserID(), loginInfo.LoginTime(), loginInfo.UserAgent(),
				loginInfo.IpAddr().String(), loginInfo.Success()).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (lhr *LoginHistoryRepository) ReadById(ctx context.Context, id uuid.UUID) (*entities.LoginInfo, error) {
	sql, args, err := postgres.Psql.Select("login_id", "user_id", "login_time", "user_agent", "ip_address", "success").
		From("user_login_histories").Where(sq.Eq{"login_id": id}).ToSql()

	if err != nil {
		return nil, err
	}

	var li entities.LoginInfo
	err = lhr.pool.QueryRow(ctx, sql, args...).Scan(&li)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &li, nil
}

func (lhr *LoginHistoryRepository) Update(ctx context.Context, loginInfo entities.LoginInfo) error {
	sql, args, err :=
		postgres.Psql.Update("user_login_histories").
			Set("user_id", loginInfo.UserID()).
			Set("login_time", loginInfo.LoginTime()).
			Set("user_agent", loginInfo.UserAgent()).
			Set("ip_address", loginInfo.IpAddr()).
			Set("success", loginInfo.Success()).
			Where(sq.Eq{"login_id": loginInfo.ID()}).
			ToSql()

	if err != nil {
		return err
	}

	_, err = lhr.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (lhr *LoginHistoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := postgres.Psql.Delete("user_login_histories").Where(sq.Eq{"login_id": id}).ToSql()

	if err != nil {
		return err
	}

	_, err = lhr.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
