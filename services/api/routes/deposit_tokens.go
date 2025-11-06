// Package routes 定义HTTP路由处理器
package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"
)

// DepositTokensHandler 处理充值代币列表查询请求
//
// HTTP端点: GET /api/v1/deposit/tokens
// 查询参数:
//   - page: 页码（默认为1）
//   - pageSize: 每页条数（默认为20，最大1000）
//
// 响应:
//   - 200 OK: 返回充值代币列表
//     示例: {"Current":1,"Size":20,"Total":100,"result":[...]}
//   - 400 Bad Request: 查询参数无效
//   - 500 Internal Server Error: 数据库查询失败
//
// 处理流程:
//   1. 从URL查询参数中提取page和pageSize
//   2. 验证和标准化参数（通过服务层）
//   3. 查询数据库获取充值代币列表
//   4. 返回分页后的JSON响应
func (h Routes) DepositTokensHandler(w http.ResponseWriter, r *http.Request) {
	// 步骤1: 提取查询参数
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")

	// 步骤2: 验证并标准化查询参数
	// 将字符串参数转换为整数，并应用验证规则
	params, err := h.svc.QueryDTListParams(pageQuery, pageSizeQuery, "asc")
	if err != nil {
		// 参数验证失败，返回400错误
		http.Error(w, "invalid query params", http.StatusBadRequest)
		log.Error("error reading request params", "err", err.Error())
		return
	}

	// 步骤3: 从数据库查询充值代币列表
	depositTokenRet, err := h.svc.GetDepositTokensList(params)
	if err != nil {
		// 数据库查询失败，返回500错误
		http.Error(w, "Internal server error reading state root list", http.StatusInternalServerError)
		log.Error("Unable to read state root list from DB", "err", err.Error())
		return
	}

	// 步骤4: 将结果序列化为JSON并返回
	err = jsonResponse(w, depositTokenRet, http.StatusOK)
	if err != nil {
		log.Error("Error writing response", "err", err.Error())
	}
}
