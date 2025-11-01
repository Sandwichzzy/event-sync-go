package worker

import (
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
)

type WithdrawManagerUpdate struct {
	GUID            uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockNumber     *big.Int       `gorm:"serializer:u256" json:"block_number"`
	WithdrawManager common.Address `gorm:"serializer:bytes" json:"withdraw_manager"`
	Timestamp       uint64         `json:"timestamp"`
}

func (WithdrawManagerUpdate) TableName() string {
	return "withdraw_manager_update"
}

type WithdrawManagerUpdateView interface {
	QueryWithdrawManagerUpdateList(page int, pageSize int, order string) ([]WithdrawManagerUpdate, uint64)
}

type WithdrawManagerUpdateDB interface {
	WithdrawManagerUpdateView
	StoreWithdrawManagerUpdates([]WithdrawManagerUpdate) error
}

type withdrawManagerUpdateDB struct {
	gorm *gorm.DB
}

func NewWithdrawManagerUpdateDB(db *gorm.DB) WithdrawManagerUpdateDB {
	return &withdrawManagerUpdateDB{gorm: db}
}

func (db *withdrawManagerUpdateDB) QueryWithdrawManagerUpdateList(page int, pageSize int, order string) ([]WithdrawManagerUpdate, uint64) {
	var (
		updates []WithdrawManagerUpdate
		total   int64
	)

	// 默认排序字段
	if order == "" {
		order = "timestamp desc"
	}

	// 统计总数
	if err := db.gorm.Model(&WithdrawManagerUpdate{}).Count(&total).Error; err != nil {
		fmt.Printf("count withdraw_manager_update error: %v\n", err)
		return nil, 0
	}

	// 分页查询
	offset := (page - 1) * pageSize
	result := db.gorm.
		Order(order).
		Limit(pageSize).
		Offset(offset).
		Find(&updates)

	if result.Error != nil {
		fmt.Printf("query withdraw_manager_update error: %v\n", result.Error)
		return nil, 0
	}

	return updates, uint64(total)
}

func (db *withdrawManagerUpdateDB) StoreWithdrawManagerUpdates(updateList []WithdrawManagerUpdate) error {
	if len(updateList) == 0 {
		return nil
	}
	result := db.gorm.CreateInBatches(&updateList, len(updateList))
	return result.Error
}
