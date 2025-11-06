// Package routes 提供路由处理器的通用工具函数
package routes

import (
	"encoding/json"
	"net/http"
)

const (
	// InternalServerError 服务器内部错误提示信息
	InternalServerError = "Internal server error"
	// defaultPageLimit 默认分页查询的最大条数限制
	defaultPageLimit    = 100
)

// jsonResponse 将数据序列化为JSON格式并写入HTTP响应
// 这是一个通用的响应辅助函数，用于标准化所有API的JSON响应格式
//
// 参数:
//   - w: HTTP响应写入器
//   - data: 要序列化的数据（任意类型）
//   - statusCode: HTTP状态码（如200, 400, 500等）
//
// 返回:
//   - error: 如果序列化或写入失败，返回错误
//
// 功能:
//   1. 设置响应头Content-Type为application/json
//   2. 将data序列化为JSON字节流
//   3. 设置HTTP状态码
//   4. 将JSON数据写入响应体
func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) error {
	// 步骤1: 设置响应内容类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 步骤2: 将数据序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		// JSON序列化失败，返回500错误
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return err
	}

	// 步骤3: 设置HTTP状态码
	w.WriteHeader(statusCode)

	// 步骤4: 写入JSON数据到响应体
	_, err = w.Write(jsonData)
	if err != nil {
		// 写入失败（通常是网络问题）
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return err
	}

	return nil
}
