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
		Columns("id", "tag", "owner_id", "created_at", "title").
		Values(chat.Id, chat.Tag, chat.OwnerId, chat.CreatedAt, chat.Title).
		ToSql()

	if err != nil {
		return err
	}

	_, err = cr.pool.Exec(ctx, sql, args...)
	log.Printf("Chat created: %v", chat.Tag)

	return err
}

func (cr *ChatRepository) AddParticipant(ctx context.Context, chatID, userID uuid.UUID) error {
    sql, args, err := cr.builder.Insert("chat_participants").
        Columns("chat_id", "user_id").
        Values(chatID, userID).
        ToSql()

    if err != nil {
        return err
    }

    _, err = cr.pool.Exec(ctx, sql, args...)
    return err
}

func (cr *ChatRepository) ReadByTag(ctx context.Context, tag string) (*entities.Chat, error) {
	sql, args, err := cr.builder.Select("id", "owner_id", "created_at", "title").
		From("chats").Where(sq.Eq{"tag": tag}).ToSql()

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

	chat.Tag = tag

	return &chat, nil
}

func (cr *ChatRepository) ReadByID(ctx context.Context, id uuid.UUID) (*entities.Chat, error) {
	sql, args, err := cr.builder.Select("tag", "owner_id", "created_at", "title").
		From("chats").Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return nil, err
	}

	var chat entities.Chat
	err = cr.pool.QueryRow(ctx, sql, args...).Scan(&chat.Tag, &chat.OwnerId, &chat.CreatedAt, &chat.Title)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	chat.Id = id

	return &chat, nil
}

func (cr *ChatRepository) Update(ctx context.Context, chat entities.Chat) error {
	sql, args, err :=
		cr.builder.Update("chats").
			Set("tag", chat.Tag).
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

func (cr *ChatRepository) RemoveParticipant(ctx context.Context, chatID, userID uuid.UUID) error {
    sql, args, err := cr.builder.Delete("chat_participants").
        Where(sq.Eq{"chat_id": chatID, "user_id": userID}).
        ToSql()

    if err != nil {
        return err
    }

    _, err = cr.pool.Exec(ctx, sql, args...)
    return err
}

func (cr *ChatRepository) RemoveAllParticipants(ctx context.Context, chatID uuid.UUID) error {
    sql, args, err := cr.builder.Delete("chat_participants").
        Where(sq.Eq{"chat_id": chatID}).
        ToSql()

    if err != nil {
        return err
    }

    _, err = cr.pool.Exec(ctx, sql, args...)
    return err
}