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
	}

	return r
}
