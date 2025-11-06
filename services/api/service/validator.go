// Package service 提供业务逻辑验证功能
package service

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

// Validator 验证器结构体，提供各种参数验证方法
// 用于在业务逻辑层验证和标准化用户输入
type Validator struct{}

// ParseValidateAddress 解析并验证以太坊地址
// 功能:
//   1. 检查地址是否为有效的十六进制格式
//   2. 确保地址不是零地址（0x0）
//   3. 将字符串地址转换为common.Address类型
//
// 参数:
//   - addr: 地址字符串（应该是0x开头的十六进制格式）
//
// 返回:
//   - common.Address: 解析后的地址对象
//   - error: 如果地址格式无效或为零地址，返回错误
func (v *Validator) ParseValidateAddress(addr string) (common.Address, error) {
	var parsedAddr common.Address

	// 特殊情况: 允许"0x00"作为有效输入
	if addr != "0x00" {
		// 验证是否为有效的十六进制地址格式
		if !common.IsHexAddress(addr) {
			return common.Address{}, errors.New("address must be represented as a valid hexadecimal string")
		}

		// 将十六进制字符串转换为Address对象
		parsedAddr = common.HexToAddress(addr)

		// 拒绝零地址
		if parsedAddr == common.HexToAddress("0x0") {
			return common.Address{}, errors.New("address cannot be the zero address")
		}
	}

	return parsedAddr, nil
}

// ValidatePage 验证并标准化页码参数
// 规则:
//   - 如果page <= 0，则返回1（第一页）
//   - 否则返回原值
//
// 参数:
//   - page: 用户提供的页码
//
// 返回:
//   - int: 验证后的页码（最小为1）
func (v *Validator) ValidatePage(page int) int {
	var validPage int
	if page <= 0 {
		validPage = 1  // 默认第一页
	} else {
		validPage = page
	}
	return validPage
}

// ValidatePageSize 验证并标准化每页条数参数
// 规则:
//   - 如果pageSize <= 0 或 pageSize > 1000，则返回20（默认值）
//   - 否则返回原值
//
// 这个限制是为了：
//   1. 防止单次查询返回过多数据，影响性能
//   2. 提供合理的默认值
//
// 参数:
//   - pageSize: 用户提供的每页条数
//
// 返回:
//   - int: 验证后的每页条数（范围1-1000，默认20）
func (v *Validator) ValidatePageSize(pageSize int) int {
	var validPageSize int
	if pageSize <= 0 || pageSize > 1000 {
		validPageSize = 20  // 默认每页20条
	} else {
		validPageSize = pageSize
	}
	return validPageSize
}

// ValidateOrder 验证并标准化排序方式参数
// 规则:
//   - 只接受 "asc", "ASC", "desc", "DESC" 四种值
//   - 其他任何值都返回"desc"（默认降序）
//
// 参数:
//   - order: 用户提供的排序方式字符串
//
// 返回:
//   - string: 验证后的排序方式（"asc"/"ASC"/"desc"/"DESC"，默认"desc"）
func (v *Validator) ValidateOrder(order string) string {
	if order == "asc" || order == "ASC" || order == "DESC" || order == "desc" {
		return order
	} else {
		return "desc"  // 默认降序
	}
}

// ValidateIdOrIndex 验证ID或索引参数
// 规则:
//   - 必须大于0
//
// 参数:
//   - idOrIndex: ID或索引值
//
// 返回:
//   - error: 如果值 <= 0，返回错误
func (v *Validator) ValidateIdOrIndex(idOrIndex uint64) error {
	if idOrIndex <= 0 {
		return errors.New("page size must be more than 0")
	}
	return nil
}
