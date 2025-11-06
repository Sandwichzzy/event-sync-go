// Package httputil 提供HTTP服务器创建和配置的实用工具
package httputil

import (
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

// timeouts 使用go-ethereum RPC包提供的默认HTTP超时配置
// 这些配置值已经过生产环境验证，适合大多数场景
var timeouts = rpc.DefaultHTTPTimeouts

// NewHttpServer 创建一个配置了标准超时的HTTP服务器
// 这是一个简单的工厂函数，用于快速创建HTTP服务器
//
// 参数:
//   - handler: HTTP请求处理器（通常是路由器或处理器链）
//
// 返回:
//   - *http.Server: 配置好超时参数的HTTP服务器实例
//
// 注意: 此函数只创建服务器对象，不启动监听。
// 建议使用StartHTTPServer函数来启动服务器。
func NewHttpServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler:           handler,                  // HTTP请求处理器
		ReadTimeout:       timeouts.ReadTimeout,     // 读取请求超时
		ReadHeaderTimeout: timeouts.ReadHeaderTimeout, // 读取请求头超时
		WriteTimeout:      timeouts.WriteTimeout,    // 写入响应超时
		IdleTimeout:       timeouts.IdleTimeout,     // 空闲连接超时
	}
}
