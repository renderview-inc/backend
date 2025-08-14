package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/renderview-inc/backend/internal/app/application/dtos"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
)

type ChatRepository interface {
	Create(ctx context.Context, chat entities.Chat) error
	ReadByTag(ctx context.Context, id string) (*entities.Chat, error)
	Update(ctx context.Context, chat entities.Chat) error
	Delete(ctx context.Context, tag string) error
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
	foundChat, err := cr.chatRepo.ReadByTag(ctx, chat.Tag)
	if err != nil {
		return fmt.Errorf("failed to check existence of chat: %w", err)
	}

	if foundChat != nil {
		return errors.New("chat with this already exists")
	}

	chatFinal := entities.NewChat(
		chat.Tag,
		chat.OwnerId,
		time.Now(),
		chat.Title,
	)

	return cr.chatRepo.Create(ctx, chatFinal)
}

func (cr *ChatService) GetByTag(ctx context.Context, tag string) (dtos.Chat, error) {
	foundChat, err := cr.chatRepo.ReadByTag(ctx, tag)
	if err != nil {
		return dtos.Chat{}, fmt.Errorf("failed to retrieve chat information: %w", err)
	}

	chat := dtos.Chat{
		Tag:     foundChat.Tag,
		OwnerId: foundChat.OwnerId,
		Title:   foundChat.Title,
	}

	return chat, nil
}

func (cr *ChatService) Update(ctx context.Context, chat dtos.Chat) error {
	foundChat, err := cr.chatRepo.ReadByTag(ctx, chat.Tag)
	if err != nil {
		return fmt.Errorf("failed to check existence of chat: %w", err)
	}

	if foundChat == nil {
		return errors.New("chat doesn't exist")
	}

	chatFinal := entities.NewChat(
		chat.Tag,
		chat.OwnerId,
		time.Now(),
		chat.Title,
	)

	return cr.chatRepo.Update(ctx, chatFinal)
}

func (cr *ChatService) Delete(ctx context.Context, tag string) error {
	foundChat, err := cr.chatRepo.ReadByTag(ctx, tag)
	if err != nil {
		return fmt.Errorf("failed to check existence of chat: %w", err)
	}

	if foundChat == nil {
		return errors.New("chat doesn't exist")
	}

	return cr.chatRepo.Delete(ctx, tag)
}
