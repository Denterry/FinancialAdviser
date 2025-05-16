package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Denterry/FinancialAdviser/Backend/brain-service/internal/entity"
	"github.com/Denterry/FinancialAdviser/Backend/brain-service/internal/repo"
	"github.com/lib/pq"
)

type ConversationRepo struct {
	db *sql.DB
}

func NewConversationRepo(db *sql.DB) *ConversationRepo {
	return &ConversationRepo{db: db}
}

func (r *ConversationRepo) Get(ctx context.Context, id string) (*entity.Conversation, error) {
	conv := &entity.Conversation{}
	query := `
		SELECT id, user_id, title, created_at, last_activity
		FROM conversations
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conv.ID, &conv.UserID, &conv.Title, &conv.CreatedAt, &conv.LastActivity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found: %v", err)
		}
		return nil, fmt.Errorf("error getting conversation: %v", err)
	}

	// Get messages
	query = `
		SELECT id, role, content, status, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error getting messages: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		msg := &entity.Message{}
		err := rows.Scan(&msg.ID, &msg.Role, &msg.Content, &msg.Status, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %v", err)
		}
		conv.Messages = append(conv.Messages, *msg)
	}

	return conv, nil
}

func (r *ConversationRepo) List(ctx context.Context, filter repo.ConversationFilter) ([]*entity.Conversation, error) {
	query := `
		SELECT id, user_id, title, created_at, last_activity
		FROM conversations
		WHERE user_id = $1
		AND created_at >= $2
		AND created_at <= $3
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`
	rows, err := r.db.QueryContext(ctx, query,
		filter.UserID, filter.StartTime, filter.EndTime, filter.Limit, filter.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error listing conversations: %v", err)
	}
	defer rows.Close()

	var conversations []*entity.Conversation
	for rows.Next() {
		conv := &entity.Conversation{}
		err := rows.Scan(
			&conv.ID, &conv.UserID, &conv.Title, &conv.CreatedAt, &conv.LastActivity,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning conversation: %v", err)
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (r *ConversationRepo) Update(ctx context.Context, conv *entity.Conversation) error {
	query := `
		UPDATE conversations
		SET title = $1, last_activity = $2
		WHERE id = $3
	`
	result, err := r.db.ExecContext(ctx, query, conv.Title, conv.LastActivity, conv.ID)
	if err != nil {
		return fmt.Errorf("error updating conversation: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

func (r *ConversationRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM conversations WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting conversation: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

func (r *ConversationRepo) AddMessage(ctx context.Context, convID string, msg *entity.Message) error {
	query := `
		INSERT INTO messages (id, conversation_id, role, content, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		msg.ID, convID, msg.Role, msg.Content, msg.Status, msg.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error adding message: %v", err)
	}

	// Update conversation's last_activity
	query = `UPDATE conversations SET last_activity = $1 WHERE id = $2`
	_, err = r.db.ExecContext(ctx, query, msg.CreatedAt, convID)
	if err != nil {
		return fmt.Errorf("error updating conversation last_activity: %v", err)
	}

	return nil
}

func (r *ConversationRepo) UpdateMessage(ctx context.Context, convID string, msg *entity.Message) error {
	query := `
		UPDATE messages
		SET content = $1, status = $2
		WHERE id = $3 AND conversation_id = $4
	`
	result, err := r.db.ExecContext(ctx, query, msg.Content, msg.Status, msg.ID, convID)
	if err != nil {
		return fmt.Errorf("error updating message: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

func (r *ConversationRepo) GetMessage(ctx context.Context, convID string, msgID string) (*entity.Message, error) {
	msg := &entity.Message{}
	query := `
		SELECT id, role, content, status, created_at
		FROM messages
		WHERE id = $1 AND conversation_id = $2
	`
	err := r.db.QueryRowContext(ctx, query, msgID, convID).Scan(
		&msg.ID, &msg.Role, &msg.Content, &msg.Status, &msg.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found: %v", err)
		}
		return nil, fmt.Errorf("error getting message: %v", err)
	}

	return msg, nil
}

func (r *ConversationRepo) Create(ctx context.Context, conv *entity.Conversation) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert conversation
	query := `
		INSERT INTO conversations (id, user_id, title, created_at, last_activity)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(ctx, query, conv.ID, conv.UserID, conv.Title, conv.CreatedAt, conv.LastActivity)
	if err != nil {
		return err
	}

	// Insert messages if any
	if len(conv.Messages) > 0 {
		for _, msg := range conv.Messages {
			query = `
				INSERT INTO messages (id, conversation_id, role, content, status, created_at)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
			_, err = tx.ExecContext(ctx, query, msg.ID, conv.ID, msg.Role, msg.Content, msg.Status, msg.CreatedAt)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}