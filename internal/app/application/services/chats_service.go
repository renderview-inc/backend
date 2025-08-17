package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/renderview-inc/backend/internal/app/application/dtos"
	"github.com/renderview-inc/backend/internal/app/domain/entities"
)

type ChatRepository interface {
	Create(ctx context.Context, chat entities.Chat) error
	AddParticipant(ctx context.Context, chatID, userID uuid.UUID) error
	ReadByTag(ctx context.Context, tag string) (*entities.Chat, error)
	ReadByID(ctx context.Context, id uuid.UUID) (*entities.Chat, error)
	Update(ctx context.Context, chat entities.Chat) error
	Delete(ctx context.Context, id uuid.UUID) error
	RemoveParticipant(ctx context.Context, chatID, userID uuid.UUID) error
	RemoveAllParticipants(ctx context.Context, chatID uuid.UUID) error
}

type ChatService struct {
	chatRepo ChatRepository
}

func NewChatService(chatRepo ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

func (cr *ChatService) Create(ctx context.Context, chat dtos.ChatRequest) (dtos.ChatResponse, error) {
    foundChat, err := cr.chatRepo.ReadByTag(ctx, chat.Tag)
    if err != nil {
        return dtos.ChatResponse{}, fmt.Errorf("failed to check existence of chat: %w", err)
    }

    if foundChat != nil {
        return dtos.ChatResponse{}, errors.New("chat with this tag already exists")
    }

    chatFinal := entities.NewChat(
        chat.Tag,
        chat.OwnerId,
        chat.Title,
    )

    err = cr.chatRepo.Create(ctx, chatFinal)
    if err != nil {
        return dtos.ChatResponse{}, fmt.Errorf("failed to create the chat: %w", err)
    }

    if err := cr.chatRepo.AddParticipant(ctx, chatFinal.Id, chatFinal.OwnerId); err != nil {
        return dtos.ChatResponse{}, fmt.Errorf("failed to add owner as participant: %w", err)
    }

    return dtos.ChatResponse{Id: chatFinal.Id.String(), Tag: chatFinal.Tag}, nil
}

func (cr *ChatService) AddParticipant(ctx context.Context, participation dtos.ChatParticipation) error {
    foundChat, err := cr.chatRepo.ReadByID(ctx, participation.ChatID)
    if err != nil {
        return fmt.Errorf("failed to check existence of chat: %w", err)
    }
    if foundChat == nil {
        return errors.New("chat doesn't exist")
    }

    return cr.chatRepo.AddParticipant(ctx, participation.ChatID, participation.UserID)
}

func (cr *ChatService) GetByTag(ctx context.Context, tag string) (dtos.ChatRequest, error) {
	foundChat, err := cr.chatRepo.ReadByTag(ctx, tag)
	if err != nil {
		return dtos.ChatRequest{}, fmt.Errorf("failed to retrieve chat information: %w", err)
	}

	chat := dtos.ChatRequest{
		Tag:     foundChat.Tag,
		OwnerId: foundChat.OwnerId,
		Title:   foundChat.Title,
	}

	return chat, nil
}

func (cr *ChatService) GetByID(ctx context.Context, id uuid.UUID) (dtos.ChatRequest, error) {
	foundChat, err := cr.chatRepo.ReadByID(ctx, id)
	if err != nil {
		return dtos.ChatRequest{}, fmt.Errorf("failed to retrieve chat information: %w", err)
	}

	chat := dtos.ChatRequest{
		Tag:     foundChat.Tag,
		OwnerId: foundChat.OwnerId,
		Title:   foundChat.Title,
	}

	return chat, nil
}

func (cr *ChatService) Update(ctx context.Context, chat dtos.ChatRequest) error {
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

	if err := cr.chatRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return cr.chatRepo.RemoveAllParticipants(ctx, id)
}

func (cr *ChatService) RemoveParticipant(ctx context.Context, participation dtos.ChatParticipation) error {
    foundChat, err := cr.chatRepo.ReadByID(ctx, participation.ChatID)
    if err != nil {
        return fmt.Errorf("failed to check existence of chat: %w", err)
    }
    if foundChat == nil {
        return errors.New("chat doesn't exist")
    }

    return cr.chatRepo.RemoveParticipant(ctx, participation.ChatID, participation.UserID)
}