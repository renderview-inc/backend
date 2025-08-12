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

type ChatRepository interface {
	Create(ctx context.Context, chat entities.Chat) error
	ReadByID(ctx context.Context, id uuid.UUID) (*entities.Chat, error)
	Update(ctx context.Context, chat entities.Chat) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ChatService struct {
	chatRepo ChatRepository
}

func NewChatService(chatRepo ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

func (cr *ChatService) Create(ctx context.Context, chat dtos.Chat) error {
	foundChat, err := cr.chatRepo.ReadByID(ctx, chat.Id)
	if err != nil {
		return fmt.Errorf("failed to check existence of chat: %w", err)
	}

	if foundChat != nil {
		return errors.New("chat with this already exists")
	}

	chatFinal := entities.NewChat(
		uuid.New(),
		chat.OwnerId,
		time.Now(),
		chat.Title,
	)

	return cr.chatRepo.Create(ctx, chatFinal)
}

func (cr *ChatService) GetByID(ctx context.Context, id uuid.UUID) (dtos.Chat, error) {
	foundChat, err := cr.chatRepo.ReadByID(ctx, id)
	if err != nil {
		return dtos.Chat{}, fmt.Errorf("failed to retreive chat information: %w", err)
	}

	chat := dtos.Chat{
		Id: foundChat.Id,
		OwnerId: foundChat.OwnerId,
		Title: foundChat.Title,
	}

	return chat, nil
}

func (cr *ChatService) Update(ctx context.Context, chat dtos.Chat) error {
	foundChat, err := cr.chatRepo.ReadByID(ctx, chat.Id)
	if err != nil {
		return fmt.Errorf("failed to check existence of chat: %w", err)
	}

	if foundChat == nil {
		return errors.New("chat doesn't exist")
	}

	chatFinal := entities.NewChat(
		chat.Id,
		chat.OwnerId,
		time.Now(),
		chat.Title,
	)

	return cr.chatRepo.Update(ctx, chatFinal)
}

func (cr *ChatService) Delete(ctx context.Context, id uuid.UUID) error {
	foundChat, err := cr.chatRepo.ReadByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check existence of chat: %w", err)
	}

	if foundChat == nil {
		return errors.New("chat doesn't exist")
	}

	return cr.chatRepo.Delete(ctx, id)
}