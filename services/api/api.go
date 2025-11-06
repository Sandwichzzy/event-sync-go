// Package api 提供HTTP API服务，用于对外暴露数据查询接口
package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Sandwichzzy/event-sync-go/config"
	"github.com/Sandwichzzy/event-sync-go/database"
	"github.com/Sandwichzzy/event-sync-go/services/api/common/httputil"
	"github.com/Sandwichzzy/event-sync-go/services/api/routes"
	"github.com/Sandwichzzy/event-sync-go/services/api/service"
)

const (
	// HealthPath 健康检查端点路径
	HealthPath          = "/healthz"
	// DepositTokensV1Path 充值代币查询API v1版本路径
	DepositTokensV1Path = "/api/v1/deposit/tokens"
)

// APIConfig API服务配置
type APIConfig struct {
	HTTPServer    config.ServerConfig // HTTP服务器配置
	MetricsServer config.ServerConfig // 指标服务器配置
}

// API 主API服务结构体，管理HTTP服务器、路由和数据库连接
type API struct {
	router    *chi.Mux              // Chi路由器，处理HTTP路由
	apiServer *httputil.HTTPServer  // HTTP服务器实例
	db        *database.DB          // 数据库连接
	stopped   atomic.Bool           // 原子布尔值，标记服务是否已停止
}

// NewApi 创建并初始化一个新的API服务实例
// 参数:
//   - ctx: 上下文对象，用于控制初始化过程
//   - cfg: 配置对象，包含服务所需的所有配置信息
// 返回:
//   - *API: 初始化成功的API实例
//   - error: 如果初始化失败，返回错误信息
func NewApi(ctx context.Context, cfg *config.Config) (*API, error) {
	out := &API{}
	if err := out.initFromConfig(ctx, cfg); err != nil {
		// 如果初始化失败，尝试清理资源并返回错误
		return nil, errors.Join(err, out.Stop(ctx))
	}
	return out, nil
}

// initFromConfig 根据配置文件初始化API服务的所有组件
// 执行步骤:
//   1. 初始化数据库连接
//   2. 初始化路由和中间件
//   3. 启动HTTP服务器
func (a *API) initFromConfig(ctx context.Context, cfg *config.Config) error {
	// 步骤1: 初始化数据库连接
	if err := a.initDB(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}
	// 步骤2: 初始化路由器
	a.initRouter(cfg.HTTPServer, cfg)
	// 步骤3: 启动HTTP服务器
	if err := a.startServer(cfg.HTTPServer); err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	return nil
}

// initRouter 初始化HTTP路由器和中间件
// 功能:
//   1. 创建服务层实例（包含验证器和数据库访问层）
//   2. 注册中间件（超时控制、错误恢复、健康检查）
//   3. 注册API路由端点
func (a *API) initRouter(conf config.ServerConfig, cfg *config.Config) {
	// 创建请求参数验证器
	v := new(service.Validator)

	// 创建服务层实例，连接验证器和数据库视图
	svc := service.New(v, a.db.DepositTokens)
	apiRouter := chi.NewRouter()
	// 创建路由处理器实例
	h := routes.NewRoutes(apiRouter, svc)

	// 中间件1: 请求超时控制（12秒）
	apiRouter.Use(middleware.Timeout(time.Second * 12))
	// 中间件2: Panic恢复中间件，防止单个请求崩溃导致整个服务down掉
	apiRouter.Use(middleware.Recoverer)

	// 中间件3: 健康检查心跳端点，用于负载均衡器探测服务状态
	apiRouter.Use(middleware.Heartbeat(HealthPath))

	// 注册API路由: GET /api/v1/deposit/tokens - 查询充值代币列表
	apiRouter.Get(fmt.Sprintf(DepositTokensV1Path), h.DepositTokensHandler)

	a.router = apiRouter
}

// initDB 初始化数据库连接
// 功能:
//   - 根据配置决定连接主库还是从库
//   - 如果启用了从库，则连接从库（用于读操作，减轻主库压力）
//   - 否则连接主库
func (a *API) initDB(ctx context.Context, cfg *config.Config) error {
	var initDb *database.DB
	var err error
	if !cfg.SlaveDbEnable {
		// 从库未启用，连接主数据库
		initDb, err = database.NewDB(ctx, cfg.MasterDB)
		if err != nil {
			log.Error("failed to connect to master database", "err", err)
			return err
		}
	} else {
		// 从库已启用，连接从数据库（读写分离架构）
		initDb, err = database.NewDB(ctx, cfg.SlaveDB)
		if err != nil {
			log.Error("failed to connect to slave database", "err", err)
			return err
		}
	}
	a.db = initDb
	return nil
}

// Start 启动API服务（当前为空实现，因为服务器在初始化时已启动）
func (a *API) Start(ctx context.Context) error {
	return nil
}

// Stop 优雅关闭API服务
// 执行步骤:
//   1. 停止HTTP服务器（等待现有请求完成）
//   2. 关闭数据库连接
//   3. 标记服务为已停止状态
// 返回所有关闭过程中产生的错误（如果有）
func (a *API) Stop(ctx context.Context) error {
	var result error
	// 步骤1: 停止HTTP服务器
	if a.apiServer != nil {
		if err := a.apiServer.Stop(ctx); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop API server: %w", err))
		}
	}
	// 步骤2: 关闭数据库连接
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to close DB: %w", err))
		}
	}
	// 步骤3: 标记为已停止
	a.stopped.Store(true)
	log.Info("API service shutdown complete")
	return result
}

// startServer 启动HTTP服务器
// 功能:
//   1. 构建监听地址（host:port）
//   2. 启动HTTP服务器并绑定路由器
//   3. 服务器在独立goroutine中运行，不阻塞主流程
func (a *API) startServer(serverConfig config.ServerConfig) error {
	log.Debug("API server listening...", "port", serverConfig.Port)
	// 构建监听地址，格式为 "host:port"
	addr := net.JoinHostPort(serverConfig.Host, strconv.Itoa(serverConfig.Port))
	// 启动HTTP服务器
	srv, err := httputil.StartHTTPServer(addr, a.router)
	if err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	log.Info("API server started", "addr", srv.Addr().String())
	a.apiServer = srv
	return nil
}

// Stopped 检查API服务是否已停止
// 返回: true表示已停止，false表示仍在运行
func (a *API) Stopped() bool {
	return a.stopped.Load()
}
