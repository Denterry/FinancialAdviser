package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/todorov_want/Desktop/HSE/FinancialAdviser/Backend/brain-service/internal/usecase"
)

func NewRouter(conversationUC usecase.ConversationUseCase) chi.Router {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Conversation routes
		newConversationRoutes(r, conversationUC)
	})

	return router
}