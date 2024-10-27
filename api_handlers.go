package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projectgo/internal/messagesService"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service *messagesService.MessageService
}

func NewHandler(service *messagesService.MessageService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := h.Service.GetAllMessages()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get messages: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) PostMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var body messagesService.Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	if body.Text == "" {
		http.Error(w, "Message text is required", http.StatusBadRequest)
		return
	}

	createdMessage, err := h.Service.CreateMessage(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create message: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(createdMessage)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) PatchMessagesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("Invalid ID: %v", idStr), http.StatusBadRequest)
		return
	}

	var body messagesService.Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	updatedMessage, err := h.Service.UpdateMessageByID(id, body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update message: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(updatedMessage)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) DeleteMessagesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("Invalid ID: %v", idStr), http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteMessageByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Message was successfully deleted")
}
