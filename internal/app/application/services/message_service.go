package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/renderview-inc/backend/internal/app/application/dtos"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
)

type MessageRepository interface {
    Create(ctx context.Context, msg *entities.Message) error
    ReadByID(ctx context.Context, id uuid.UUID) (*entities.Message, error)
    GetLastByChatTag(ctx context.Context, chatTag string) (*entities.Message, error)
    Update(ctx context.Context, msg *entities.Message) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type MessageService struct {
    msgRepo MessageRepository
}

func NewMessageService(msgRepo MessageRepository) *MessageService {
    return &MessageService{
        msgRepo: msgRepo,
    }
}

func (ms *MessageService) Create(ctx context.Context, msg dtos.Message) error {
    msgEntity := entities.NewMessage(
        uuid.New(),
        msg.ReplyToID,
        msg.UserID,
        msg.ChatTag,
        msg.Content,
        time.Now(),
    )
    return ms.msgRepo.Create(ctx, msgEntity)
}

func (ms *MessageService) GetByID(ctx context.Context, id uuid.UUID) (dtos.Message, error) {
    msgEntity, err := ms.msgRepo.ReadByID(ctx, id)
    if err != nil {
        return dtos.Message{}, fmt.Errorf("failed to retrieve message: %w", err)
    }
    if msgEntity == nil {
        return dtos.Message{}, errors.New("message not found")
    }
    return dtos.Message{
        ID:      msgEntity.ID,
        UserID:  msgEntity.UserID,
        ChatTag:  msgEntity.ChatTag,
        Content: msgEntity.Content,
    }, nil
}

func (ms *MessageService) GetLastByChatTag(ctx context.Context, chatTag string) (dtos.Message, error) {
    msgEntity, err := ms.msgRepo.GetLastByChatTag(ctx, chatTag)
    if err != nil {
        return dtos.Message{}, fmt.Errorf("failed to get last message: %w", err)
    }
    if msgEntity == nil {
        return dtos.Message{}, errors.New("no messages found")
    }
    
    msg := dtos.Message{
        ID:     msgEntity.ID,
        UserID: msgEntity.UserID,
        ChatTag: msgEntity.ChatTag,
        Content: msgEntity.Content,
    }

    return msg, nil
}

func (ms *MessageService) Update(ctx context.Context, msg dtos.Message) error {
    msgEntity, err := ms.msgRepo.ReadByID(ctx, msg.ID)
    if err != nil {
        return fmt.Errorf("failed to check existence of message: %w", err)
    }
    if msgEntity == nil {
        return errors.New("message doesn't exist")
    }

    return ms.msgRepo.Update(ctx, msgEntity)
}

func (ms *MessageService) Delete(ctx context.Context, id uuid.UUID) error {
    msgEntity, err := ms.msgRepo.ReadByID(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to check existence of message: %w", err)
    }
    if msgEntity == nil {
        return errors.New("message doesn't exist")
    }
    return ms.msgRepo.Delete(ctx, id)
}