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

type ChatRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (cr *ChatRepository) Create(ctx context.Context, chat entities.Chat) error {
	sql, args, err := cr.builder.Insert("chats").
		Columns("id", "owner_id", "created_at", "title").
		Values(chat.Id, chat.OwnerId, chat.CreatedAt, chat.Title).
		ToSql()

	if err != nil {
		return err
	}

	_, err = cr.pool.Exec(ctx, sql, args...)
	log.Printf("Chat created: %v", chat.Id)

	return err
}

func (cr *ChatRepository) ReadByID(ctx context.Context, id uuid.UUID) (*entities.Chat, error) {
	sql, args, err := cr.builder.Select("id", "owner_id", "created_at", "title").
		From("chats").Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return nil, err
	}

	var chat entities.Chat
	err = cr.pool.QueryRow(ctx, sql, args...).Scan(&chat.Id, &chat.OwnerId, &chat.CreatedAt, &chat.Title)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &chat, nil
}

func (cr *ChatRepository) Update(ctx context.Context, chat entities.Chat) error {
	sql, args, err :=
		cr.builder.Update("chats").
			Set("id", chat.Id).
			Set("title", chat.Title).
			Set("owner_id", chat.OwnerId).
			Set("created_at", chat.CreatedAt).
			Where(sq.Eq{"id": chat.Id}).
			ToSql()

	if err != nil {
		return err
	}

	_, err = cr.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ChatRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := cr.builder.Delete("chats").Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return err
	}

	_, err = cr.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}