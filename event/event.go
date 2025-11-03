package event

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Sandwichzzy/event-sync-go/common/bigint"
	"github.com/Sandwichzzy/event-sync-go/common/tasks"
	"github.com/Sandwichzzy/event-sync-go/database"
	"github.com/Sandwichzzy/event-sync-go/database/common"
	"github.com/Sandwichzzy/event-sync-go/database/event"
	"github.com/Sandwichzzy/event-sync-go/event/contracts"
)

type EventProcessorConfig struct {
	LoopInterval    time.Duration
	EventStartBlock uint64 // 事件起始区块
	EventBlockStep  uint64 // 每次处理的区块步长
}

type EventProcessor struct {
	db                *database.DB
	eventBlocksConfig *EventProcessorConfig
	resourceCtx       context.Context
	resourceCancel    context.CancelFunc
	tasks             tasks.Group
	LatestBlockHeader *common.BlockHeader
	TreasureManager   *contracts.TreasureManager
}

func NewEventProcessor(db *database.DB, eventBlocksConfig *EventProcessorConfig, shutdown context.CancelCauseFunc) (*EventProcessor, error) {
	//创建合约处理器（核心业务逻辑）
	treasureManager, err := contracts.NewTreasureManager(db)
	if err != nil {
		log.Error("new treasure manager fail", "err", err)
	}
	//获取最新处理的事件区块
	latestBlockHeader, err := db.EventBlocks.LatestEventBlockHeader()
	if err != nil {
		log.Error("get latest event block header fail", "err", err)
		return nil, err
	}

	resCtx, resCancel := context.WithCancel(context.Background())

	return &EventProcessor{
		db:                db,
		eventBlocksConfig: eventBlocksConfig,
		resourceCtx:       resCtx,
		resourceCancel:    resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in bridge processor: %w", err))
		}},
		LatestBlockHeader: latestBlockHeader,
		TreasureManager:   treasureManager,
	}, nil
}

func (ep *EventProcessor) Start() error {
	log.Info("starting bridge processor...")
	// 创建定时器，按配置间隔执行
	tickerWorker := time.NewTicker(ep.eventBlocksConfig.LoopInterval)
	ep.tasks.Go(func() error {
		for range tickerWorker.C {
			err := ep.processEvent()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func (ep *EventProcessor) Close() error {
	ep.resourceCancel()
	return ep.tasks.Wait()
}

func (ep *EventProcessor) processEvent() error {
	// 1. 确定起始区块号
	lastBlockNumber := big.NewInt(int64(ep.eventBlocksConfig.EventStartBlock)) //配置里读
	if ep.LatestBlockHeader != nil {
		lastBlockNumber = ep.LatestBlockHeader.Number //如果数据库有就用数据库的
	}

	log.Info("process event latest block number", "lastBlockNumber", lastBlockNumber)
	// 2. 构建数据库查询作用域
	latestHeaderScope := func(db *gorm.DB) *gorm.DB {
		newQuery := db.Session(&gorm.Session{NewDB: true})
		headers := newQuery.Model(common.BlockHeader{}).Where("number > ?", lastBlockNumber)
		return db.Where("number = (?)", newQuery.Table("(?) as block_numbers", headers.Order("number ASC").Limit(int(ep.eventBlocksConfig.EventBlockStep))).Select("MAX(number)"))
	}

	if latestHeaderScope == nil {
		return nil
	}
	// 3. 查询需要处理的最新区块头
	latestHeader, err := ep.db.Blocks.BlockHeaderWithScope(latestHeaderScope)
	if err != nil {
		return fmt.Errorf("failed to query for latest unfinalized L1 state: %w", err)
	} else if latestHeader == nil {
		log.Debug("no new  state to process event")
		return nil
	}
	// 4. 计算处理范围
	fromHeight, toHeight := new(big.Int).Add(lastBlockNumber, bigint.One), latestHeader.Number

	// 5. 构建事件区块记录
	eventBlocks := make([]event.EventBlocks, 0, toHeight.Uint64()-fromHeight.Uint64())

	for index := fromHeight.Uint64(); index < toHeight.Uint64(); index++ {
		blockHeader, err := ep.db.Blocks.BlockHeaderByNumber(big.NewInt(int64(index)))
		if err != nil {
			return err
		}
		evBlock := event.EventBlocks{
			GUID:       uuid.New(),
			Hash:       blockHeader.Hash,
			ParentHash: blockHeader.ParentHash,
			Number:     blockHeader.Number,
			Timestamp:  blockHeader.Timestamp,
		}
		eventBlocks = append(eventBlocks, evBlock)
	}
	// 6. 处理合约事件（核心业务逻辑）
	log.Info("parse contract event start", "fromHeight", fromHeight.String(), "toHeight", toHeight.String())
	depositTokens, grantsRewardTokens, withdrawManagerUpdates, withdrawTokens, err := ep.TreasureManager.ProcessTreasureManagerEvents(fromHeight, toHeight)
	if err != nil {
		log.Error("parse treasure manager contracts events fail", "err", err)
		return err
	}
	// 7. 数据库事务：保存所有处理结果
	if err := ep.db.Transaction(func(tx *database.DB) error {
		if err != nil {
			log.Error("process delegation event fail", "err", err)
			return err
		}

		if len(depositTokens) > 0 {
			err := ep.db.DepositTokens.StoreDepositTokens(depositTokens)
			if err != nil {
				log.Error("store deposit tokens fail", "err", err)
				return err
			}
		}

		if len(withdrawTokens) > 0 {
			err := ep.db.WithdrawTokens.StoreWithdrawTokens(withdrawTokens)
			if err != nil {
				log.Error("store withdraw tokens fail", "err", err)
				return err
			}
		}

		if len(grantsRewardTokens) > 0 {
			err := ep.db.GrantRewardTokens.StoreGrantRewardTokens(grantsRewardTokens)
			if err != nil {
				log.Error("store grants reward tokens fail", "err", err)
				return err
			}
		}

		if len(withdrawManagerUpdates) > 0 {
			err := ep.db.WithdrawManagerUpdate.StoreWithdrawManagerUpdates(withdrawManagerUpdates)
			if err != nil {
				log.Error("store withdraw manager update fail", "err", err)
				return err
			}
		}

		if len(eventBlocks) > 0 {
			err = tx.EventBlocks.StoreEventBlocks(eventBlocks)
			if err != nil {
				log.Error("store event block fail", "err", err)
				return err
			}
		}
		return nil
	}); err != nil {
		log.Error("exec database fail", "err", err)
		return err
	}
	// 8. 更新最新处理区块头
	ep.LatestBlockHeader = latestHeader
	return nil

}
