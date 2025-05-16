package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/entity"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/usecase"
)

type conversationRoutes struct {
	convUseCase usecase.ConversationUseCase
}

func newConversationRoutes(handler chi.Router, uc usecase.ConversationUseCase) {
	r := &conversationRoutes{convUseCase: uc}

	handler.Route("/conversations", func(router chi.Router) {
		router.Post("/", r.createConversation)
		router.Get("/", r.listConversations)
		router.Route("/{id}", func(router chi.Router) {
			router.Get("/", r.getConversation)
			router.Post("/messages", r.addMessage)
			router.Delete("/", r.deleteConversation)
		})
	})
}

type createConversationRequest struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
}

func (r *conversationRoutes) createConversation(w http.ResponseWriter, req *http.Request) {
	var request createConversationRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	conv := &entity.Conversation{
		UserID: request.UserID,
		Title:  request.Title,
	}

	result, err := r.convUseCase.Create(req.Context(), conv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (r *conversationRoutes) listConversations(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	convs, err := r.convUseCase.List(req.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convs)
}

func (r *conversationRoutes) getConversation(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	convID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	conv, err := r.convUseCase.Get(req.Context(), convID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conv)
}

type addMessageRequest struct {
	Content string `json:"content"`
}

func (r *conversationRoutes) addMessage(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	convID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	var request addMessageRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	msg := &entity.Message{
		Role:    entity.UserRole,
		Content: request.Content,
	}

	response, err := r.convUseCase.AddMessage(req.Context(), convID.String(), msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (r *conversationRoutes) deleteConversation(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	convID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	if err := r.convUseCase.Delete(req.Context(), convID.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}