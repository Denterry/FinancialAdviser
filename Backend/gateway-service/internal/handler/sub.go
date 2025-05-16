package handler

import (
	"net/http"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct{ svc *service.SubscriptionService }

func NewSubscriptionHandler(s *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: s}
}

func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	plans, err := h.svc.GetPlans(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var body struct {
		PlanID string `json:"plan_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Subscribe(c, body.PlanID); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *SubscriptionHandler) GetStatus(c *gin.Context) {
	status, err := h.svc.Status(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *SubscriptionHandler) Cancel(c *gin.Context) {
	if err := h.svc.Cancel(c); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}
