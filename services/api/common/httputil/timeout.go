// Package httputil 提供HTTP服务器相关的实用工具
package httputil

import "time"

// DefaultTimeouts HTTP服务器的默认超时配置
// 这些超时值用于防止慢速客户端攻击和资源泄漏
var DefaultTimeouts = HTTPTimeouts{
	ReadTimeout:       30 * time.Second,  // 读取整个请求（包括body）的最大时间
	ReadHeaderTimeout: 30 * time.Second,  // 读取请求头的最大时间
	WriteTimeout:      30 * time.Second,  // 写入响应的最大时间
	IdleTimeout:       120 * time.Second, // Keep-Alive连接的最大空闲时间
}

// HTTPTimeouts HTTP服务器超时配置结构体
// 包含HTTP服务器的各种超时参数
type HTTPTimeouts struct {
	ReadTimeout       time.Duration // 读取完整请求的超时时间（从连接建立到body读取完毕）
	ReadHeaderTimeout time.Duration // 读取请求头的超时时间（更短，防止慢速攻击）
	WriteTimeout      time.Duration // 写入响应的超时时间（从请求头读取完毕到响应写入完毕）
	IdleTimeout       time.Duration // Keep-Alive连接的空闲超时时间
}

// withTimeouts 创建一个应用自定义超时配置的HTTPOption
// 这是一个工厂函数，返回一个配置函数
//
// 参数:
//   - timeouts: 自定义的超时配置
//
// 返回:
//   - HTTPOption: 可以传递给StartHTTPServer的配置函数
func withTimeouts(timeouts HTTPTimeouts) HTTPOption {
	return func(s *HTTPServer) error {
		s.srv.ReadTimeout = timeouts.ReadTimeout
		s.srv.ReadHeaderTimeout = timeouts.ReadHeaderTimeout
		s.srv.WriteTimeout = timeouts.WriteTimeout
		s.srv.IdleTimeout = timeouts.IdleTimeout
		return nil
	}
}
