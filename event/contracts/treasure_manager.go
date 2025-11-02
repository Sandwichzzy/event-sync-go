package contracts

import (
	"context"

	"github.com/Sandwichzzy/event-sync-go/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type TreasureManager struct {
	TmAbi     *abi.ABI
	TmFilter  *bindings.TreasureManagerFilterer
	TmContext context.Context
}

func NewTreasureManager() (*TreasureManager, error) {
	treasureManagerAbi, err = bindings.TreasureManagerMetaData.GetAbi()
	if err != nil {
		log.Error("binding treasure manager data abi fail", "err", err)
		return nil, err
	}

	treasureManagerFilter, err := bindings.NewTreasureManagerFilterer(common.Address{}, nil)
}
