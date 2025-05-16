package http

import (
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterPublic - mounts unauthenticated routes: /api/auth/…
func RegisterPublic(r *gin.RouterGroup, h *handler.Handlers) {
	auth := r.Group("/auth")
	{
		auth.POST("/signup", h.Auth.Register)
		auth.POST("/signin", h.Auth.Login)
		// auth.POST("/signout", h.Auth.Logout)
	}
}
