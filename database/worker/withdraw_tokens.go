package worker

import (
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
)

type WithdrawTokens struct {
	GUID         uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockNumber  *big.Int       `gorm:"serializer:u256" json:"block_number"`
	TokenAddress common.Address `gorm:"serializer:bytes" json:"token_address"`
	Sender       common.Address `gorm:"serializer:bytes" json:"sender"`
	Receiver     common.Address `gorm:"serializer:bytes" json:"receiver"`
	Amount       *big.Int       `gorm:"serializer:u256" json:"amount"`
	Timestamp    uint64         `json:"timestamp"`
}

func (WithdrawTokens) TableName() string {
	return "withdraw_tokens"
}

type WithdrawTokensView interface {
	QueryWithdrawTokensList(page int, pageSize int, order string) ([]WithdrawTokens, uint64)
}

type WithdrawTokensDB interface {
	WithdrawTokensView
	StoreWithdrawTokens([]WithdrawTokens) error
}

type withdrawTokensDB struct {
	gorm *gorm.DB
}

func NewWithdrawTokensDB(db *gorm.DB) WithdrawTokensDB {
	return &withdrawTokensDB{gorm: db}
}

func (db *withdrawTokensDB) QueryWithdrawTokensList(page int, pageSize int, order string) ([]WithdrawTokens, uint64) {
	var (
		withdraws []WithdrawTokens
		total     int64
	)

	if order == "" {
		order = "timestamp desc"
	}

	// 统计总数
	if err := db.gorm.Model(&WithdrawTokens{}).Count(&total).Error; err != nil {
		fmt.Printf("count withdraw_tokens error: %v\n", err)
		return nil, 0
	}

	// 分页查询
	offset := (page - 1) * pageSize
	result := db.gorm.
		Order(order).
		Limit(pageSize).
		Offset(offset).
		Find(&withdraws)

	if result.Error != nil {
		fmt.Printf("query withdraw_tokens error: %v\n", result.Error)
		return nil, 0
	}

	return withdraws, uint64(total)
}

func (db *withdrawTokensDB) StoreWithdrawTokens(withdrawTokensList []WithdrawTokens) error {
	if len(withdrawTokensList) == 0 {
		return nil
	}
	result := db.gorm.CreateInBatches(&withdrawTokensList, len(withdrawTokensList))
	return result.Error
}
