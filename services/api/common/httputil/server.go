// Package httputil 提供HTTP服务器的封装和管理功能
package httputil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
)

// HTTPServer HTTP服务器封装结构体
// 提供对标准库http.Server的增强封装，添加优雅关闭和状态管理功能
type HTTPServer struct {
	listener net.Listener  // TCP监听器，用于接收连接
	srv      *http.Server  // 底层的标准库HTTP服务器
	closed   atomic.Bool   // 原子布尔值，标记服务器是否已关闭
}

// HTTPOption HTTP服务器配置选项函数类型
// 使用函数式选项模式，允许灵活配置服务器
type HTTPOption func(srv *HTTPServer) error

// StartHTTPServer 启动一个新的HTTP服务器
// 功能:
//   1. 创建TCP监听器并绑定到指定地址
//   2. 配置HTTP服务器（超时、处理器等）
//   3. 在后台goroutine中启动服务器
//   4. 应用可选的配置选项
//
// 参数:
//   - addr: 监听地址，格式为"host:port"（如"0.0.0.0:8080"）
//   - handler: HTTP请求处理器（通常是路由器）
//   - opts: 可选的配置函数，用于自定义服务器设置
//
// 返回:
//   - *HTTPServer: 启动成功的服务器实例
//   - error: 如果启动失败（端口被占用、配置错误等），返回错误
func StartHTTPServer(addr string, handler http.Handler, opts ...HTTPOption) (*HTTPServer, error) {
	// 步骤1: 创建TCP监听器
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to address %q: %w", addr, err)
	}

	// 步骤2: 创建可取消的上下文，用于优雅关闭
	srvCtx, srvCancel := context.WithCancel(context.Background())

	// 步骤3: 配置HTTP服务器
	srv := &http.Server{
		Handler:           handler,
		ReadTimeout:       DefaultTimeouts.ReadTimeout,       // 读超时
		ReadHeaderTimeout: DefaultTimeouts.ReadHeaderTimeout, // 读头超时
		WriteTimeout:      DefaultTimeouts.WriteTimeout,      // 写超时
		IdleTimeout:       DefaultTimeouts.IdleTimeout,       // 空闲超时
		BaseContext: func(listener net.Listener) context.Context {
			return srvCtx  // 为所有请求提供基础上下文
		},
	}

	// 步骤4: 创建HTTPServer封装对象
	out := &HTTPServer{listener: listener, srv: srv}

	// 步骤5: 应用可选配置
	for _, opt := range opts {
		if err := opt(out); err != nil {
			srvCancel()
			return nil, errors.Join(fmt.Errorf("failed to apply HTTP option: %w", err), listener.Close())
		}
	}

	// 步骤6: 在后台goroutine中启动服务器
	go func() {
		err := out.srv.Serve(listener)
		srvCancel()  // 服务器停止后取消上下文
		if errors.Is(err, http.ErrServerClosed) {
			// 正常关闭
			out.closed.Store(true)
		} else {
			// 异常退出，触发panic
			panic(fmt.Errorf("unexpected serve error: %w", err))
		}
	}()

	return out, nil
}

// Closed 检查服务器是否已关闭
// 返回: true表示已关闭，false表示仍在运行
func (s *HTTPServer) Closed() bool {
	return s.closed.Load()
}

// Stop 停止HTTP服务器（优雅关闭）
// 功能:
//   1. 尝试优雅关闭（等待现有请求完成）
//   2. 如果上下文超时，则强制关闭所有连接
//
// 参数:
//   - ctx: 控制关闭超时的上下文
//
// 返回:
//   - error: 如果关闭失败，返回错误
func (s *HTTPServer) Stop(ctx context.Context) error {
	if err := s.Shutdown(ctx); err != nil {
		if errors.Is(err, ctx.Err()) {
			// 上下文超时/取消，强制关闭所有连接
			return s.Close()
		}
		return err
	}
	return nil
}

// Shutdown 优雅关闭HTTP服务器
// 不接受新连接，等待现有连接处理完毕
// 也会关闭底层的监听器
//
// 参数:
//   - ctx: 控制关闭超时的上下文
//
// 返回:
//   - error: 如果关闭失败，返回错误
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	// 关闭HTTP服务器（也会关闭底层监听器）
	return s.srv.Shutdown(ctx)
}

// Close 立即关闭HTTP服务器
// 不等待现有连接，强制关闭所有连接和监听器
//
// 返回:
//   - error: 如果关闭失败，返回错误
func (s *HTTPServer) Close() error {
	// 强制关闭（也会关闭底层监听器）
	return s.srv.Close()
}

// Addr 获取服务器监听的网络地址
// 返回: 监听地址（如"127.0.0.1:8080"）
func (s *HTTPServer) Addr() net.Addr {
	return s.listener.Addr()
}

// WithMaxHeaderBytes 创建一个设置最大请求头大小的配置选项
// 用于防止恶意客户端发送超大请求头
//
// 参数:
//   - max: 最大请求头字节数
//
// 返回:
//   - HTTPOption: 配置函数
func WithMaxHeaderBytes(max int) HTTPOption {
	return func(srv *HTTPServer) error {
		srv.srv.MaxHeaderBytes = max
		return nil
	}
}
