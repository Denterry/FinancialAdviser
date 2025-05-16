package repo

import (
	"context"

	"github.com/Denterry/FinancialAdviser/Backend/brain-service/internal/entity"
)

// ConversationFilter defines filters for querying conversations
type ConversationFilter struct {
	UserID    string
	StartTime int64
	EndTime   int64
	Limit     int32
	Offset    int32
}

// ConversationRepo defines the interface for conversation storage operations
type ConversationRepo interface {
	// Create stores a new conversation
	Create(ctx context.Context, conv *entity.Conversation) error

	// Get retrieves a conversation by ID
	Get(ctx context.Context, id string) (*entity.Conversation, error)

	// List retrieves conversations based on filter criteria
	List(ctx context.Context, filter ConversationFilter) ([]*entity.Conversation, error)

	// Update modifies an existing conversation
	Update(ctx context.Context, conv *entity.Conversation) error

	// Delete removes a conversation by ID
	Delete(ctx context.Context, id string) error

	// AddMessage appends a new message to a conversation
	AddMessage(ctx context.Context, convID string, msg *entity.Message) error

	// UpdateMessage modifies an existing message
	UpdateMessage(ctx context.Context, convID string, msg *entity.Message) error

	// GetMessage retrieves a specific message from a conversation
	GetMessage(ctx context.Context, convID string, msgID string) (*entity.Message, error)
}