package v1

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/renderview-inc/backend/internal/app/application/dtos"
	"github.com/renderview-inc/backend/internal/app/application/services"
)

type MessageHandler struct {
	messageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) MessageHandler {
	return MessageHandler{
		messageService: messageService,
	}
}

func (mh *MessageHandler) HandleCreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg dtos.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := mh.messageService.Create(r.Context(), msg); err != nil {
		http.Error(w, "failed to create message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (mh *MessageHandler) HandleGetMessage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid ID; must be UUID", http.StatusBadRequest)
		return
	}

	msg, err := mh.messageService.GetByID(r.Context(), parsedUUID)
	if err != nil {
		http.Error(w, "failed to retrieve message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(msg); err != nil {
		http.Error(w, "failed to encode message", http.StatusInternalServerError)
		return
	}
}

func (mh *MessageHandler) HandleGetLastMessageByChatTag(w http.ResponseWriter, r *http.Request) {
	chatTag := r.URL.Query().Get("chat_tag")
	if chatTag == "" {
		http.Error(w, "chat_tag is required", http.StatusBadRequest)
		return
	}

	msg, err := mh.messageService.GetLastByChatTag(r.Context(), chatTag)
	if err != nil {
		http.Error(w, "failed to get last message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		http.Error(w, "failed to encode message", http.StatusInternalServerError)
		return
	}
}

func (mh *MessageHandler) HandleUpdateMessage(w http.ResponseWriter, r *http.Request) {
	var msg dtos.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := mh.messageService.Update(r.Context(), msg); err != nil {
		http.Error(w, "failed to update message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (mh *MessageHandler) HandleDeleteMessage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid ID; must be UUID", http.StatusBadRequest)
		return
	}

	if err := mh.messageService.Delete(r.Context(), parsedUUID); err != nil {
		http.Error(w, "failed to delete message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
