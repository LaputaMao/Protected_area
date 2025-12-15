package handler

import (
	"ProtectedArea/internal/model"
	"ProtectedArea/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProtectedTypeMap 定义了中文保护地类型到英文缩写的映射
var ProtectedTypeMap = map[string]string{
	"NP":               "NP", // 允许直接传英文
	"NR":               "NR",
	"FP":               "FP",
	"WP":               "WP",
	"GP":               "GP",
	"DP":               "DP",
	"SH":               "SH",
	"国家公园":         "NP",
	"国家级自然保护区": "NR",
	"森林公园":         "FP",
	"湿地公园":         "WP",
	"地质公园":         "GP",
	"荒漠公园":         "DP",
	"风景名胜区":       "SH",
}

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

	// 1. 预处理 ProtectedType 字段
	if req.ProtectedType != "" {
		// 建议先进行 trim 和 to upper/lower，以防前端传入空格或大小写不一致
		protectedTypeKey := strings.TrimSpace(req.ProtectedType)

		// 查找映射表并更新请求结构体字段
		// 注意：这里调用的是前面定义的 util.MapProtectedType 函数
		req.ProtectedType = MapProtectedType(protectedTypeKey)
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

	// 1. 预处理 ProtectedType 字段
	if req.ProtectedType != "" {
		// 建议先进行 trim 和 to upper/lower，以防前端传入空格或大小写不一致
		protectedTypeKey := strings.TrimSpace(req.ProtectedType)

		// 查找映射表并更新请求结构体字段
		// 注意：这里调用的是前面定义的 util.MapProtectedType 函数
		req.ProtectedType = MapProtectedType(protectedTypeKey)
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

	// 1. 预处理 ProtectedType 字段
	if req.ProtectedType != "" {
		// 建议先进行 trim 和 to upper/lower，以防前端传入空格或大小写不一致
		protectedTypeKey := strings.TrimSpace(req.ProtectedType)

		// 查找映射表并更新请求结构体字段
		// 注意：这里调用的是前面定义的 util.MapProtectedType 函数
		req.ProtectedType = MapProtectedType(protectedTypeKey)
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

func (h *NatureHandler) GetPatchImage(c *gin.Context) {
	// 1. 获取参数
	tbbh := c.Query("tbbh")
	if tbbh == "" {
		c.String(http.StatusBadRequest, "图斑编号不能为空")
		return
	}

	// 2. 调用 Service 查找文件
	filePath, exists := h.srv.GetImagePath(tbbh)

	// 3. 根据结果返回
	if !exists {
		// 按照你的要求，返回纯文本
		c.String(http.StatusOK, "暂无图片")
		return
	}

	// Gin 自带的方法，会自动设置 Content-Type 并流式传输文件
	c.File(filePath)
}

func MapProtectedType(chineseName string) string {
	// 检查映射表，如果找到则返回英文缩写，否则返回原始输入
	if abbr, ok := ProtectedTypeMap[chineseName]; ok {
		return abbr
	}
	return chineseName // 如果找不到，返回原始值，让后续查询逻辑处理错误
}
