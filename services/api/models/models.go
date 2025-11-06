// Package models 定义API层的数据模型和请求/响应结构体
package models

import "github.com/Sandwichzzy/event-sync-go/database/worker"

// QueryDTParams 充值代币列表查询参数
// 用于分页查询时传递查询条件
type QueryDTParams struct {
	Page     int    // 页码
	PageSize int    // 每页条数
	Order    string // 排序方式（"asc"升序或"desc"降序）
}

// DepositTokensResponse 充值代币列表的API响应结构
// 采用标准的分页响应格式，包含元数据和结果列表
type DepositTokensResponse struct {
	Current int                    `json:"Current"` // 当前页码
	Size    int                    `json:"Size"`    // 当前页条数
	Total   int64                  `json:"Total"`   // 总记录数
	Result  []worker.DepositTokens `json:"result"`  // 当前页的充值代币数据列表
}
