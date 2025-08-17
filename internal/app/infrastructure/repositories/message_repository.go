package repositories

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
)

type MessageRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewMessageRepository(pool *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (mr *MessageRepository) Create(ctx context.Context, msg *entities.Message) error {
	builder := mr.builder.Insert("messages").
		Columns("id", "user_id", "chat_tag", "content", "created_at").
		Values(msg.ID, msg.UserID, msg.ChatTag, msg.Content, msg.CreatedAt)

	if msg.ReplyToID != uuid.Nil {
		builder.Columns("reply_to").Values(msg.ReplyToID)
	}

	sql, args, err := builder.ToSql()

	if err != nil {
		return err
	}

	_, err = mr.pool.Exec(ctx, sql, args...)
	return err
}

func (mr *MessageRepository) ReadByID(ctx context.Context, id uuid.UUID) (*entities.Message, error) {
	sql, args, err := mr.builder.Select("id", "reply_to", "user_id", "chat_tag", "content", "created_at").
		From("messages").Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return nil, err
	}

	var msg entities.Message
	err = mr.pool.QueryRow(ctx, sql, args...).Scan(&msg.ID, &msg.UserID, &msg.ChatTag, &msg.Content, &msg.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &msg, nil
}

func (mr *MessageRepository) GetLastByChatTag(ctx context.Context, chatTag string) (*entities.Message, error) {
	sql, args, err := mr.builder.
		Select("id", "reply_to", "user_id", "chat_tag", "content", "created_at").
		From("messages").
		Where(sq.Eq{"chat_tag": chatTag}).
		OrderBy("created_at DESC").
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	var msg entities.Message
	err = mr.pool.QueryRow(ctx, sql, args...).Scan(&msg.ID, &msg.UserID, &msg.ChatTag, &msg.Content, &msg.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &msg, nil
}

func (mr *MessageRepository) Update(ctx context.Context, msg *entities.Message) error {
	sql, args, err := mr.builder.Update("messages").
		Set("user_id", msg.UserID).
		Set("chat_tag", msg.ChatTag).
		Set("content", msg.Content).
		Set("created_at", msg.CreatedAt).
		Where(sq.Eq{"id": msg.ID}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = mr.pool.Exec(ctx, sql, args...)
	return err
}

func (mr *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := mr.builder.Delete("messages").Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return err
	}

	_, err = mr.pool.Exec(ctx, sql, args...)
	return err
}
