package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/brain-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/brain-service/internal/repo"
	"github.com/google/uuid"
)

type ConversationService struct {
	repo    repo.ConversationRepo
	llm     LLMProvider
	contextEnricher ContextEnricher
}

type ContextEnricher interface {
	// EnrichContext adds relevant financial context to the conversation
	EnrichContext(ctx context.Context, content string) (string, error)
}

func NewConversationService(repo repo.ConversationRepo, llm LLMProvider, enricher ContextEnricher) *ConversationService {
	return &ConversationService{
		repo:    repo,
		llm:     llm,
		contextEnricher: enricher,
	}
}

func (s *ConversationService) CreateConversation(ctx context.Context, userID, title string) (*entity.Conversation, error) {
	conv := &entity.Conversation{
		ID:           uuid.New().String(),
		UserID:       userID,
		Title:        title,
		CreatedAt:    time.Now().Unix(),
		LastActivity: time.Now().Unix(),
		Messages:     []entity.Message{},
	}

	err := s.repo.Create(ctx, conv)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %v", err)
	}

	return conv, nil
}

func (s *ConversationService) GetConversation(ctx context.Context, id string) (*entity.Conversation, error) {
	return s.repo.Get(ctx, id)
}

func (s *ConversationService) ListConversations(ctx context.Context, filter repo.ConversationFilter) ([]*entity.Conversation, error) {
	return s.repo.List(ctx, filter)
}

func (s *ConversationService) DeleteConversation(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ConversationService) SendMessage(ctx context.Context, convID string, content string) (*entity.Message, error) {
	// Get the conversation
	conv, err := s.repo.Get(ctx, convID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}

	// Create user message
	userMsg := &entity.Message{
		ID:        uuid.New().String(),
		Role:      entity.UserRole,
		Content:   content,
		Status:    entity.MessageStatusComplete,
		CreatedAt: time.Now().Unix(),
	}

	// Add user message to conversation
	err = s.repo.AddMessage(ctx, convID, userMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to add user message: %v", err)
	}

	// Enrich context with financial information
	enrichedContent, err := s.contextEnricher.EnrichContext(ctx, content)
	if err != nil {
		// Log the error but continue with original content
		fmt.Printf("failed to enrich context: %v\n", err)
		enrichedContent = content
	}

	// Prepare assistant message
	assistantMsg := &entity.Message{
		ID:        uuid.New().String(),
		Role:      entity.AssistantRole,
		Status:    entity.MessageStatusPending,
		CreatedAt: time.Now().Unix(),
	}

	// Add pending assistant message
	err = s.repo.AddMessage(ctx, convID, assistantMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to add assistant message: %v", err)
	}

	// Generate response using LLM
	response, err := s.llm.GenerateResponse(ctx, append(conv.Messages, *userMsg))
	if err != nil {
		// Update message status to failed
		assistantMsg.Status = entity.MessageStatusFailed
		assistantMsg.Content = fmt.Sprintf("Failed to generate response: %v", err)
		s.repo.UpdateMessage(ctx, convID, assistantMsg)
		return nil, fmt.Errorf("failed to generate response: %v", err)
	}

	// Update assistant message with response
	assistantMsg.Content = response
	assistantMsg.Status = entity.MessageStatusComplete
	err = s.repo.UpdateMessage(ctx, convID, assistantMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to update assistant message: %v", err)
	}

	return assistantMsg, nil
}

func (s *ConversationService) GetMessage(ctx context.Context, convID, msgID string) (*entity.Message, error) {
	return s.repo.GetMessage(ctx, convID, msgID)
}