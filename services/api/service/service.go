// Package service 业务服务层，负责业务逻辑处理和数据访问
package service

import (
	"strconv"

	"github.com/Sandwichzzy/event-sync-go/database/worker"
	"github.com/Sandwichzzy/event-sync-go/services/api/models"
)

// Service 业务服务接口，定义所有API的业务逻辑方法
// 该接口将路由处理器与具体实现解耦，便于测试和扩展
type Service interface {
	// GetDepositTokensList 获取充值代币分页列表
	// 参数: 查询参数（页码、每页条数、排序方式）
	// 返回: 分页响应数据和可能的错误
	GetDepositTokensList(*models.QueryDTParams) (*models.DepositTokensResponse, error)

	// QueryDTListParams 验证并构建查询参数对象
	// 参数: 页码字符串、每页条数字符串、排序方式字符串
	// 返回: 验证后的查询参数对象和可能的错误
	QueryDTListParams(page string, pageSize string, order string) (*models.QueryDTParams, error)
}

// HandlerSvc 业务服务实现结构体
// 它组合了验证器和数据访问层，实现Service接口
type HandlerSvc struct {
	v                 *Validator                // 参数验证器
	depositTokensView worker.DepositTokensView  // 充值代币数据访问层
}

// GetDepositTokensList 获取充值代币分页列表
// 功能:
//   1. 调用数据访问层查询数据库
//   2. 构建分页响应对象
//
// 参数:
//   - params: 已验证的查询参数（页码、每页条数、排序方式）
//
// 返回:
//   - *models.DepositTokensResponse: 包含当前页码、每页条数、总记录数和结果列表的响应对象
//   - error: 如果查询失败，返回错误（当前实现总是返回nil）
func (h HandlerSvc) GetDepositTokensList(params *models.QueryDTParams) (*models.DepositTokensResponse, error) {
	// 调用数据访问层查询数据库
	dtList, totalCount := h.depositTokensView.QueryDepositTokensList(params.Page, params.PageSize)

	// 构建并返回分页响应对象
	return &models.DepositTokensResponse{
		Current: params.Page,         // 当前页码
		Size:    params.PageSize,     // 每页条数
		Total:   int64(totalCount),   // 总记录数
		Result:  dtList,              // 当前页的数据列表
	}, nil
}

// QueryDTListParams 验证并构建充值代币列表查询参数
// 功能:
//   1. 将字符串参数转换为整数
//   2. 验证参数合法性
//   3. 应用默认值和范围限制
//
// 参数:
//   - page: 页码字符串（需要转换为整数）
//   - pageSize: 每页条数字符串（需要转换为整数）
//   - order: 排序方式字符串（"asc"或"desc"）
//
// 返回:
//   - *models.QueryDTParams: 验证后的查询参数对象
//   - error: 如果参数格式错误（无法转换为整数），返回错误
func (h HandlerSvc) QueryDTListParams(page string, pageSize string, order string) (*models.QueryDTParams, error) {
	// 将页码字符串转换为整数
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	// 验证页码（如果<=0则默认为1）
	pageVal := h.v.ValidatePage(pageInt)

	// 将每页条数字符串转换为整数
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, err
	}
	// 验证每页条数（如果<=0或>1000则默认为20）
	pageSizeVal := h.v.ValidatePageSize(pageSizeInt)

	// 验证排序方式（只接受"asc"或"desc"，否则默认为"desc"）
	orderBy := h.v.ValidateOrder(order)

	// 构建并返回验证后的查询参数对象
	return &models.QueryDTParams{
		Page:     pageVal,
		PageSize: pageSizeVal,
		Order:    orderBy,
	}, nil
}

// New 创建一个新的业务服务实例
// 参数:
//   - v: 参数验证器实例
//   - dtv: 充值代币数据访问层接口
// 返回:
//   - Service: 业务服务接口的实现
func New(v *Validator, dtv worker.DepositTokensView) Service {
	return &HandlerSvc{
		v:                 v,
		depositTokensView: dtv,
	}
}

// GetDepositList 获取充值列表（与GetDepositTokensList功能相同）
// 注意: 这个方法似乎是GetDepositTokensList的重复实现，可能需要重构
func (h HandlerSvc) GetDepositList(params *models.QueryDTParams) (*models.DepositTokensResponse, error) {
	depositList, total := h.depositTokensView.QueryDepositTokensList(params.Page, params.PageSize)
	return &models.DepositTokensResponse{
		Current: params.Page,
		Size:    params.PageSize,
		Total:   int64(total),
		Result:  depositList,
	}, nil
}
