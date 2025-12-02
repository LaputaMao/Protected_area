package handler

import (
	"ProtectedArea/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NatureHandler struct {
	srv service.NatureService
}

func NewNatureHandler(srv service.NatureService) *NatureHandler {
	return &NatureHandler{srv: srv}
}

// GetTrendStats 接口入口
func (h *NatureHandler) GetTrendStats(c *gin.Context) {
	data, err := h.srv.GetTrendAnalysis()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	c.JSON(http.StatusOK, data)
}
