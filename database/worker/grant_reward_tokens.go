package worker

import (
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
)

type GrantRewardTokens struct {
	GUID         uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockNumber  *big.Int       `gorm:"serializer:u256" json:"block_number"`
	TokenAddress common.Address `gorm:"serializer:bytes" json:"token_address"`
	Granter      common.Address `gorm:"serializer:bytes" json:"granter"`
	Amount       *big.Int       `gorm:"serializer:u256" json:"amount"`
	Timestamp    uint64         `json:"timestamp"`
}

func (GrantRewardTokens) TableName() string {
	return "grant_reward_tokens"
}

type GrantRewardTokensView interface {
	QueryGrantRewardTokensList(page int, pageSize int, order string) ([]GrantRewardTokens, uint64)
}

type GrantRewardTokensDB interface {
	GrantRewardTokensView
	StoreGrantRewardTokens([]GrantRewardTokens) error
}

type grantRewardTokensDB struct {
	gorm *gorm.DB
}

func NewGrantRewardTokensDB(db *gorm.DB) GrantRewardTokensDB {
	return &grantRewardTokensDB{gorm: db}
}

func (db *grantRewardTokensDB) QueryGrantRewardTokensList(page int, pageSize int, order string) ([]GrantRewardTokens, uint64) {
	var (
		rewards []GrantRewardTokens
		total   int64
	)

	// 默认排序字段
	if order == "" {
		order = "timestamp desc"
	}

	// 统计总数
	if err := db.gorm.Model(&GrantRewardTokens{}).Count(&total).Error; err != nil {
		fmt.Printf("count grant_reward_tokens error: %v\n", err)
		return nil, 0
	}

	// 分页查询
	offset := (page - 1) * pageSize
	result := db.gorm.
		Order(order).
		Limit(pageSize).
		Offset(offset).
		Find(&rewards)

	if result.Error != nil {
		fmt.Printf("query grant_reward_tokens error: %v\n", result.Error)
		return nil, 0
	}

	return rewards, uint64(total)
}

func (db *grantRewardTokensDB) StoreGrantRewardTokens(grantRewardTokensList []GrantRewardTokens) error {
	if len(grantRewardTokensList) == 0 {
		return nil
	}
	result := db.gorm.CreateInBatches(&grantRewardTokensList, len(grantRewardTokensList))
	return result.Error
}
