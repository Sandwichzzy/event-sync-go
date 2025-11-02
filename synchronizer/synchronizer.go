package synchronizer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Sandwichzzy/event-sync-go/common/retry"
	"github.com/Sandwichzzy/event-sync-go/common/tasks"
	"github.com/Sandwichzzy/event-sync-go/config"
	"github.com/Sandwichzzy/event-sync-go/database"
	common2 "github.com/Sandwichzzy/event-sync-go/database/common"
	"github.com/Sandwichzzy/event-sync-go/database/event"
	"github.com/Sandwichzzy/event-sync-go/database/utils"
	"github.com/Sandwichzzy/event-sync-go/synchronizer/node"
)

type Synchronizer struct {
	ethClient node.EthClient
	db        *database.DB

	loopInterval     time.Duration         // 同步循环间隔
	headerBufferSize uint64                // 每次处理的区块数量
	headerTraversal  *node.HeaderTraversal // 区块遍历器

	headers      []types.Header // 待处理的区块头缓存
	latestHeader *types.Header  // 最新区块头

	startHeight       *big.Int
	confirmationDepth *big.Int
	chainCfg          *config.ChainConfig

	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group // 任务管理组
}

func NewSynchronizer(cfg *config.Config, db *database.DB, client node.EthClient, shutdown context.CancelCauseFunc) (*Synchronizer, error) {
	//确定同步起始点
	latestHeader, err := db.Blocks.LatestBlockHeader()
	if err != nil {
		log.Error("query block header database error", "err", err)
		return nil, err
	}
	var fromHeader *types.Header
	if latestHeader != nil {
		// 从数据库最后同步的区块继续
		log.Info("sync detected last indexed block", "number", latestHeader.Number, "hash", latestHeader.Hash)
		fromHeader = latestHeader.RLPHeader.Header()
	} else if cfg.Chain.StartingHeight > 0 {
		// 从配置的起始高度开始
		log.Info("no sync indexed state starting from supplied ethereum height", "height", cfg.Chain.StartingHeight)
		header, err := client.BlockHeaderByNumber(big.NewInt(int64(cfg.Chain.StartingHeight)))
		if err != nil {
			return nil, fmt.Errorf("could not fetch starting block header: %w", err)
		}
		fromHeader = header
	} else {
		// 从头开始同步
		log.Info("no eth wallet indexed state")
	}

	//创建区块遍历器（注意：确认深度为0，实时同步）
	headerTraversal := node.NewHeaderTraversal(client, fromHeader, big.NewInt(0), cfg.Chain.ChainId)

	// 创建同步器实例
	resCtx, resCancel := context.WithCancel(context.Background())
	return &Synchronizer{
		loopInterval:     time.Duration(cfg.Chain.LoopInterval) * time.Second,
		headerBufferSize: uint64(cfg.Chain.BlockStep),
		headerTraversal:  headerTraversal,
		ethClient:        client,
		latestHeader:     fromHeader,
		db:               db,
		chainCfg:         &cfg.Chain,
		resourceCtx:      resCtx,
		resourceCancel:   resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in Synchronizer: %w", err))
		}},
	}, nil

}

func (syncer *Synchronizer) Start() error {
	// tickerSyncer := time.NewTicker(syncer.loopInterval)
	tickerSyncer := time.NewTicker(time.Second * 3)
	syncer.tasks.Go(func() error {
		for range tickerSyncer.C {
			if len(syncer.headers) > 0 {
				// 重试机制：先处理之前失败的批次
				log.Info("retrying previous batch")
			} else {
				// 获取新的区块头批次
				newHeaders, err := syncer.headerTraversal.NextHeaders(uint64(syncer.chainCfg.BlockStep))
				if err != nil {
					log.Error("error querying for headers", "err", err)
					continue
				} else if len(newHeaders) == 0 {
					log.Warn("no new headers. syncer at head?")
				} else {
					syncer.headers = newHeaders // 缓存待处理区块头
				}
				latestHeader := syncer.headerTraversal.LatestHeader()
				if latestHeader != nil {
					log.Info("Latest header", "latestHeader Number", latestHeader.Number)
				}
			}

			// 处理当前批次
			err := syncer.processBatch(syncer.headers, syncer.chainCfg)
			if err == nil {
				syncer.headers = nil // 处理成功，清空缓存
			}
			// 如果处理失败，headers不会被清空，下次循环会重试
		}
		return nil
	})
	return nil
}

func (syncer *Synchronizer) processBatch(headers []types.Header, chainCfg *config.ChainConfig) error {
	if len(headers) == 0 {
		return nil
	}
	// 1. 记录处理范围
	firstHeader, lastHeader := headers[0], headers[len(headers)-1]
	log.Info("extracting batch", "size", len(headers), "startBlock", firstHeader.Number.String(), "endBlock", lastHeader.Number.String())

	// 2. 构建区块头哈希映射（用于后续日志关联）
	headerMap := make(map[common.Hash]*types.Header, len(headers))

	for i := range headers {
		header := headers[i]
		headerMap[header.Hash()] = &header
	}
	log.Info("chainCfg Contracts", "contract address", chainCfg.Contracts[0])
	// 3. 查询合约事件日志
	filterQuery := ethereum.FilterQuery{FromBlock: firstHeader.Number, ToBlock: lastHeader.Number, Addresses: chainCfg.Contracts}
	logs, err := syncer.ethClient.FilterLogs(filterQuery)
	if err != nil {
		log.Info("failed to extract logs", "err", err)
		return err
	}
	// 4. 验证区块一致性（防止链重组）
	if logs.ToBlockHeader.Number.Cmp(lastHeader.Number) != 0 {
		return fmt.Errorf("mismatch in FilterLog#ToBlock number")
	} else if logs.ToBlockHeader.Hash() != lastHeader.Hash() {
		return fmt.Errorf("mismatch in FitlerLog#ToBlock block hash")
	}

	// 5. 记录检测到的事件数量
	if len(logs.Logs) > 0 {
		log.Info("detected logs", "size", len(logs.Logs))
	}
	// 6. 准备区块头数据存储
	blockHeaders := make([]common2.BlockHeader, 0, len(headers))
	for i := range headers {
		if headers[i].Number == nil {
			continue
		}
		bHeader := common2.BlockHeader{
			Hash:       headers[i].Hash(),
			ParentHash: headers[i].ParentHash,
			Number:     headers[i].Number,
			Timestamp:  headers[i].Time,
			RLPHeader:  (*utils.RLPHeader)(&headers[i]),
		}
		blockHeaders = append(blockHeaders, bHeader)
	}

	// 7. 转换事件日志数据
	chainContractEvent := make([]event.ContractEvent, len(logs.Logs))
	for i := range logs.Logs {
		logEvent := logs.Logs[i]
		if _, ok := headerMap[logEvent.BlockHash]; !ok {
			continue // 跳过不属于当前批次的日志
		}
		timestamp := headerMap[logEvent.BlockHash].Time
		chainContractEvent[i] = event.ContractEventFromLog(&logs.Logs[i], timestamp)
	}
	// 8. 数据库存储（带重试机制）
	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](syncer.resourceCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := syncer.db.Transaction(func(tx *database.DB) error {
			// 原子性存储：区块头和事件要么都成功，要么都失败
			if err := tx.Blocks.StoreBlockHeaders(blockHeaders); err != nil {
				return err
			}
			if err := tx.ContractEvent.StoreContractEvents(chainContractEvent); err != nil {
				return err
			}
			return nil
		}); err != nil {
			log.Info("unable to persist batch", err)
			return nil, fmt.Errorf("unable to persist batch: %w", err)
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

func (syncer *Synchronizer) Close() error {
	return nil
}
