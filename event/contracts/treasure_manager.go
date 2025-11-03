package contracts

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"

	"github.com/Sandwichzzy/event-sync-go/bindings"
	"github.com/Sandwichzzy/event-sync-go/config"
	"github.com/Sandwichzzy/event-sync-go/database"
	"github.com/Sandwichzzy/event-sync-go/database/event"
	"github.com/Sandwichzzy/event-sync-go/database/worker"
)

type TreasureManager struct {
	db        *database.DB
	TmAbi     *abi.ABI
	TmFilter  *bindings.TreasureManagerFilterer
	TmContext context.Context
}

func NewTreasureManager(db *database.DB) (*TreasureManager, error) {
	treasureManagerAbi, err := bindings.TreasureManagerMetaData.GetAbi()
	if err != nil {
		log.Error("binding treasure manager data abi fail", "err", err)
		return nil, err
	}

	treasureManagerFilter, err := bindings.NewTreasureManagerFilterer(common.Address{}, nil)
	if err != nil {
		log.Error("get latest event block header fail", "err", err)
		return nil, err
	}

	return &TreasureManager{
		db:        db,
		TmAbi:     treasureManagerAbi,
		TmFilter:  treasureManagerFilter,
		TmContext: context.Background(),
	}, nil
}

func (tm *TreasureManager) ProcessTreasureManagerEvents(fromHeight *big.Int, toHeight *big.Int) ([]worker.DepositTokens, []worker.GrantRewardTokens, []worker.WithdrawManagerUpdate, []worker.WithdrawTokens, error) {
	contractEventFilter := event.ContractEvent{ContractAddress: common.HexToAddress(config.TreasureManagerAddr)}
	log.Info("query contracts filter",
		"TreasureManagerAddr", config.TreasureManagerAddr,
		"fromHeight", fromHeight,
		"toHeight", toHeight,
	)
	contractEventList, err := tm.db.ContractEvent.ContractEventsWithFilter(contractEventFilter, fromHeight, toHeight)
	if err != nil {
		log.Error("filter contract event by address and start/end block fail", "err", err)
		return nil, nil, nil, nil, err
	}

	var (
		depositTokens          []worker.DepositTokens
		grantsRewardTokens     []worker.GrantRewardTokens
		withdrawManagerUpdates []worker.WithdrawManagerUpdate
		withdrawTokens         []worker.WithdrawTokens
	)

	for _, eventItem := range contractEventList {
		rlpLog := eventItem.RLPLog
		// DepositToken
		if eventItem.EventSignature.String() == tm.TmAbi.Events["DepositToken"].ID.String() {
			depositTokenEvent, err := tm.TmFilter.ParseDepositToken(*rlpLog)
			if err != nil {
				log.Error("Parse deposit token event fail", "err", err)
				return nil, nil, nil, nil, err
			}
			log.Info("DepositToken Information",
				"TokenAddress", depositTokenEvent.TokenAddress,
				"Amount", depositTokenEvent.Amount,
				"Sender", depositTokenEvent.Sender)

			tempDepositToken := worker.DepositTokens{
				GUID:         uuid.New(),
				BlockNumber:  big.NewInt(int64(eventItem.RLPLog.BlockNumber)),
				TokenAddress: depositTokenEvent.TokenAddress,
				Sender:       depositTokenEvent.Sender,
				Amount:       depositTokenEvent.Amount,
				Timestamp:    uint64(time.Now().Unix()),
			}

			depositTokens = append(depositTokens, tempDepositToken)
		}

		// WithdrawToken
		if eventItem.EventSignature.String() == tm.TmAbi.Events["WithdrawToken"].ID.String() {
			withdrawTokenEvent, err := tm.TmFilter.ParseWithdrawToken(*rlpLog)
			if err != nil {
				log.Error("Parse withdraw token event fail", "err", err)
				return nil, nil, nil, nil, err
			}
			log.Info("WithdrawToken Event",
				"TokenAddress", withdrawTokenEvent.TokenAddress,
				"Amount", withdrawTokenEvent.Amount,
				"Sender", withdrawTokenEvent.Sender,
			)
			tempWithdrawToken := worker.WithdrawTokens{
				GUID:         uuid.New(),
				BlockNumber:  big.NewInt(int64(eventItem.RLPLog.BlockNumber)),
				TokenAddress: withdrawTokenEvent.TokenAddress,
				Sender:       withdrawTokenEvent.Sender,
				Receiver:     common.Address{},
				Amount:       withdrawTokenEvent.Amount,
				Timestamp:    uint64(time.Now().Unix()),
			}
			withdrawTokens = append(withdrawTokens, tempWithdrawToken)
		}
		// GrantRewardTokenAmount
		if eventItem.EventSignature.String() == tm.TmAbi.Events["GrantRewardTokenAmount"].ID.String() {
			grantRewardEvent, err := tm.TmFilter.ParseGrantRewardTokenAmount(*rlpLog)
			if err != nil {
				log.Error("Parse grant reward token event fail", "err", err)
				return nil, nil, nil, nil, err
			}

			log.Info("GrantRewardTokenAmount Event",
				"TokenAddress", grantRewardEvent.TokenAddress,
				"Amount", grantRewardEvent.Amount,
				"Granter", grantRewardEvent.Granter,
			)

			tempgrantsRewardToken := worker.GrantRewardTokens{
				GUID:         uuid.New(),
				BlockNumber:  big.NewInt(int64(eventItem.RLPLog.BlockNumber)),
				TokenAddress: grantRewardEvent.TokenAddress,
				Granter:      grantRewardEvent.Granter,
				Amount:       grantRewardEvent.Amount,
				Timestamp:    uint64(time.Now().Unix()),
			}
			grantsRewardTokens = append(grantsRewardTokens, tempgrantsRewardToken)
		}

		// WithdrawManagerUpdate
		if eventItem.EventSignature.String() == tm.TmAbi.Events["WithdrawManagerUpdate"].ID.String() {
			withdrawManagerEvent, err := tm.TmFilter.ParseWithdrawManagerUpdate(*rlpLog)
			if err != nil {
				log.Error("Parse withdraw manager update event fail", "err", err)
				return nil, nil, nil, nil, err
			}

			log.Info("WithdrawManagerUpdate Event",
				"WithdrawManager", withdrawManagerEvent.WithdrawManager,
			)

			tempWithdrawManagerUpdate := worker.WithdrawManagerUpdate{
				GUID:            uuid.New(),
				BlockNumber:     big.NewInt(int64(eventItem.RLPLog.BlockNumber)),
				WithdrawManager: withdrawManagerEvent.WithdrawManager,
				Timestamp:       uint64(time.Now().Unix()),
			}
			withdrawManagerUpdates = append(withdrawManagerUpdates, tempWithdrawManagerUpdate)
		}
	}
	return depositTokens, grantsRewardTokens, withdrawManagerUpdates, withdrawTokens, nil
}
