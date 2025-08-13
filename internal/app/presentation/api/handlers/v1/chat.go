package v1

import (
	"encoding/json"
	"net/http"

	"github.com/renderview-inc/backend/internal/app/application/dtos"
	"github.com/renderview-inc/backend/internal/app/application/services"
)

type ChatHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(chatService *services.ChatService) ChatHandler {
	return ChatHandler{
		chatService: chatService,
	}
}

func (ch *ChatHandler) HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	var chat dtos.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := ch.chatService.Create(r.Context(), chat); err != nil {
		http.Error(w, "failed to create chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ch *ChatHandler) HandleGetChatInfo(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")

	chat, err := ch.chatService.GetByTag(r.Context(), tag)
	if err != nil {
		http.Error(w, "failed to retreive chat info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, "failed to encode chat info", http.StatusInternalServerError)
		return
	}
}

func (ch *ChatHandler) HandleUpdateChat(w http.ResponseWriter, r *http.Request) {
	var chat dtos.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := ch.chatService.Update(r.Context(), chat); err != nil {
		http.Error(w, "failed to update chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ch *ChatHandler) HandleDeleteChat(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")

	if err := ch.chatService.Delete(r.Context(), tag); err != nil {
		http.Error(w, "failed to delete chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}