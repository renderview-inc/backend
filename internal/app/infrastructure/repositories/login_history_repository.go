package repositories

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
	postgres "github.com/renderview-inc/backend/pkg/connections"
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
			Values(loginInfo.GetID(), loginInfo.UserGetID(), loginInfo.GetLoginTime(), loginInfo.GetUserAgent(),
				loginInfo.GetIpAddr().String(), loginInfo.GetSuccess()).ToSql()
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
	err = lhr.pool.QueryRow(ctx, sql, args...).Scan(&li.Id, &li.UserID, &li.LoginTime, &li.UserAgent, &li.IpAddr, &li.Success)

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
			Set("user_id", loginInfo.UserGetID()).
			Set("login_time", loginInfo.GetLoginTime()).
			Set("user_agent", loginInfo.GetUserAgent()).
			Set("ip_address", loginInfo.GetIpAddr()).
			Set("success", loginInfo.GetSuccess()).
			Where(sq.Eq{"login_id": loginInfo.GetID()}).
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
