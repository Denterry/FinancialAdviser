package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/notification-service/internal/mailer"
	"github.com/Denterry/FinancialAdviser/Backend/notification-service/internal/storage"
	"github.com/Denterry/FinancialAdviser/Backend/notification-service/internal/types"
)

type NotificationService struct {
	Storage *storage.Repository
	Mailer  mailer.Mailer
}

func NewNotificationService(repo *storage.Repository, mailer mailer.Mailer) *NotificationService {
	return &NotificationService{
		Storage: repo,
		Mailer:  mailer,
	}
}

// ProcessNewPost filters users and sends post to matching ones
func (s *NotificationService) ProcessNewPost(ctx context.Context, post types.FinancialPost) {
	userIDs, err := s.Storage.FindInterestedUsers(ctx, post.Symbols)
	if err != nil {
		fmt.Printf("error fetching users for symbols: %v\n", err)
		return
	}

	for _, uid := range userIDs {
		if s.Storage.HasUserSeenPost(ctx, uid, post.ID) {
			continue
		}

		_ = s.Storage.MarkPostAsSeen(ctx, uid, post.ID)

		email, err := s.Storage.GetUserEmail(ctx, uid)
		if err != nil {
			fmt.Printf("could not fetch email for user %s: %v\n", uid, err)
			continue
		}

		subject := "📈 New Financial Post for You"
		body := fmt.Sprintf("Hello,\n\nA new post about %v has been published:\n\nTitle: %s\n\n%s\n\nPosted at: %s",
			post.Symbols, post.Title, post.Content, post.Timestamp.Format(time.RFC1123))

		err = s.Mailer.SendEmail(email, subject, body)
		if err != nil {
			fmt.Printf("failed to send email to %s: %v\n", email, err)
		} else {
			fmt.Printf("Email sent to %s about post %s\n", email, post.ID)
		}
	}
}
