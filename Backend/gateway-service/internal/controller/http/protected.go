package http

import (
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterProtected - mounts authenticated routes under /api
func RegisterProtected(r *gin.RouterGroup, h *handler.Handlers) {
	// chat := r.Group("/chat")
	// {
	// 	chat.GET("/", h.Chat.ListChats)
	// 	chat.POST("/", h.Chat.CreateChat)
	// 	chat.GET("/:id/messages", h.Chat.GetMessages)
	// 	chat.POST("/:id/messages", h.Chat.SendMessage)
	// }
}
