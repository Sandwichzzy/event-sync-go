// Package routes 定义HTTP路由处理器，负责处理具体的HTTP请求
package routes

import (
	"github.com/Sandwichzzy/event-sync-go/services/api/service"
	"github.com/go-chi/chi/v5"
)

// Routes 路由处理器结构体
// 它连接HTTP路由器和业务服务层，负责：
//   1. 接收HTTP请求
//   2. 调用服务层处理业务逻辑
//   3. 返回HTTP响应
type Routes struct {
	router *chi.Mux           // Chi路由器实例，用于路由匹配
	svc    service.Service    // 业务服务层接口，包含实际业务逻辑
}

// NewRoutes 创建一个新的路由处理器实例
// 参数:
//   - r: Chi路由器实例
//   - svc: 业务服务层实例
// 返回:
//   - Routes: 初始化完成的路由处理器
func NewRoutes(r *chi.Mux, svc service.Service) Routes {
	return Routes{
		router: r,
		svc:    svc,
	}
}
