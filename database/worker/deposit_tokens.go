package worker

import (
	"fmt"
	"math/big"

	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type DepositTokens struct {
	GUID         uuid.UUID      `gorm:"primaryKey" json:"guid"`
	BlockNumber  *big.Int       `gorm:"serializer:u256" json:"block_number"`
	TokenAddress common.Address `json:"token_address" gorm:"serializer:bytes"`
	Sender       common.Address `json:"sender" gorm:"serializer:bytes"`
	Amount       *big.Int       `gorm:"serializer:u256"`
	Timestamp    uint64
}

func (DepositTokens) TableName() string {
	return "deposit_tokens"
}

type DepositTokensView interface {
	QueryDepositTokensList(page int, pageSize int) ([]DepositTokens, uint64)
	QueryDepositTokensById(string) (*DepositTokens, error)
}

type DepositTokensDB interface {
	DepositTokensView

	StoreDepositTokens([]DepositTokens) error
}

type depositTokensDB struct {
	gorm *gorm.DB
}

func (db depositTokensDB) QueryDepositTokensById(guid string) (*DepositTokens, error) {
	var dts DepositTokens
	if err := db.gorm.Table("deposit_tokens").Where("guid", guid).Take(&dts).Error; err != nil {
		return nil, err
	}
	return &dts, nil
}

func (db depositTokensDB) QueryDepositTokensList(page int, pageSize int) ([]DepositTokens, uint64) {
	var (
		depositTokens []DepositTokens
		total         int64
	)

	if err := db.gorm.Table("deposit_tokens").Count(&total).Error; err != nil {
		fmt.Printf("count deposit_tokens error: %v\n", err)
		return nil, 0
	}

	offset := (page - 1) * pageSize
	result := db.gorm.
		Limit(pageSize).
		Offset(offset).
		Find(&depositTokens)

	if result.Error != nil {
		fmt.Printf("query deposit_tokens error: %v\n", result.Error)
		return nil, 0
	}

	return depositTokens, uint64(total)
}

func (db depositTokensDB) StoreDepositTokens(depositTokensList []DepositTokens) error {
	result := db.gorm.CreateInBatches(&depositTokensList, len(depositTokensList))
	return result.Error
}

func NewDepositTokensDB(db *gorm.DB) DepositTokensDB {
	return &depositTokensDB{gorm: db}
}
