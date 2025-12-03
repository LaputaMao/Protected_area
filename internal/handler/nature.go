package handler

import (
	"ProtectedArea/internal/model"
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

// GetProtectedAreaStats 接口1: 保护地统计
func (h *NatureHandler) GetProtectedAreaStats(c *gin.Context) {
	var req model.NatureQueryRequest
	// ShouldBindQuery 自动把 URL 参数绑定到结构体
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.srv.GetProtectedAreaStats(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetSpotList 接口2: 图斑明细
func (h *NatureHandler) GetSpotList(c *gin.Context) {
	var req model.NatureQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.srv.GetSpotList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetTransitionStats 接口3: 流向分析 (饼图)
func (h *NatureHandler) GetTransitionStats(c *gin.Context) {
	var req model.NatureQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 接口3 必填 qlx
	if req.QLX == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "前地类(qlx)参数必填"})
		return
	}

	data, err := h.srv.GetTransitionStats(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetLargeSpots 大面积图斑预警接口
func (h *NatureHandler) GetLargeSpots(c *gin.Context) {
	var req model.AlertQueryRequest
	// 绑定参数
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 简单的业务校验 (可选)
	if req.AlertArea < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "预警面积必须大于等于0"})
		return
	}

	data, err := h.srv.GetLargeSpots(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, data)
}
