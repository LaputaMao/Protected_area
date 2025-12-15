package router

import (
	"ProtectedArea/internal/handler"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(natureHandler *handler.NatureHandler) *gin.Engine {
	r := gin.Default()

	// 可以在这里加跨域中间件等

	api := r.Group("/api")
	{
		api.GET("/stats/trend", natureHandler.GetTrendStats)

		// 1. 年度概况: /api/stats/overview?year=2023
		api.GET("/stats/overview", natureHandler.GetYearlyOverview)

		// 2. 分批次损毁统计: /api/stats/damage-batch?year=2023
		api.GET("/stats/damage-batch", natureHandler.GetDamageBatchStats)

		// 3. 行政区划统计: /api/stats/region?year=2025&scope=province&name=河北省
		api.GET("/stats/region", natureHandler.GetRegionStats)

		// 4. 保护地统计
		api.GET("/stats/protected-area", natureHandler.GetProtectedAreaStats)

		// 5. 图斑明细
		api.GET("/stats/spot-list", natureHandler.GetSpotList)

		// 6. 流向分析 (饼图)
		api.GET("/stats/transition", natureHandler.GetTransitionStats)

		// 7. 大图斑预警 : /api/stats/alert/large-spots?year=2022&alert_area=500&page=1
		api.GET("/stats/alert/large-spots", natureHandler.GetLargeSpots)

		// 8. 获取图斑图片: /api/image?tbbh=110109202202NR001
		api.GET("/image", natureHandler.GetPatchImage)
	}

	return r
}
