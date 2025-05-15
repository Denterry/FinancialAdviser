package handler

import (
	"net/http"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth         *AuthHandler
	Subscription *SubscriptionHandler
	ML           *MLHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth:         NewAuthHandler(services.Auth),
		Subscription: NewSubscriptionHandler(services.Subscription),
		ML:           NewMLHandler(services.ML),
	}
}

type AuthHandler struct {
	service *service.AuthService
}

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

type MLHandler struct {
	service *service.MLService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func NewMLHandler(service *service.MLService) *MLHandler {
	return &MLHandler{service: service}
}

// Auth handlers
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement registration
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement login
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement token refresh
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// Subscription handlers
func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	plans, err := h.service.GetPlans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var req struct {
		PlanID string `json:"plan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Subscribe(c.Request.Context(), req.PlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription successful"})
}

func (h *SubscriptionHandler) GetStatus(c *gin.Context) {
	status, err := h.service.GetStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *SubscriptionHandler) Cancel(c *gin.Context) {
	err := h.service.Cancel(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled"})
}

// ML handlers
func (h *MLHandler) Analyze(c *gin.Context) {
	var req struct {
		Data []byte `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analysis, err := h.service.Analyze(c.Request.Context(), req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

func (h *MLHandler) GetRecommendations(c *gin.Context) {
	recommendations, err := h.service.GetRecommendations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}
