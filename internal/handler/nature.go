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

// GetYearlyOverview 1. 接口：获取年度概况
func (h *NatureHandler) GetYearlyOverview(c *gin.Context) {
	year := c.Query("year")
	if year == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "年份参数(year)不能为空"})
		return
	}

	data, err := h.srv.GetYearlyOverview(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetDamageBatchStats 2. 接口：获取资源损毁分批次统计
func (h *NatureHandler) GetDamageBatchStats(c *gin.Context) {
	year := c.Query("year")
	if year == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "年份参数(year)不能为空"})
		return
	}

	data, err := h.srv.GetDamageAnalysisByBatch(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *NatureHandler) GetRegionStats(c *gin.Context) {
	// 获取参数
	year := c.Query("year")
	scope := c.Query("scope")
	name := c.Query("name") // 可选，没传就是空字符串

	// 必填校验
	if year == "" || scope == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数 year 和 scope 不能为空"})
		return
	}

	// 调用 Service
	data, err := h.srv.GetAdministrativeStats(year, scope, name)
	if err != nil {
		// 如果是业务逻辑报错（比如县级查下级），返回 400
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
