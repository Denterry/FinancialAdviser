package handler

import (
	"net/http"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/service"
	"github.com/gin-gonic/gin"
)

type MLHandler struct{ svc *service.MLService }

func NewMLHandler(s *service.MLService) *MLHandler { return &MLHandler{svc: s} }

func (h *MLHandler) Analyze(c *gin.Context) {
	var req struct {
		Data []byte `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.Analyze(c, req.Data)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *MLHandler) GetRecommendations(c *gin.Context) {
	recs, err := h.svc.Recommend(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recs)
}
