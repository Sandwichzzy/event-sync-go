package event_sync_go

import (
	"context"
	"sync/atomic"

	"github.com/Sandwichzzy/event-sync-go/event"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Sandwichzzy/event-sync-go/config"
	"github.com/Sandwichzzy/event-sync-go/database"
	"github.com/Sandwichzzy/event-sync-go/synchronizer"
	"github.com/Sandwichzzy/event-sync-go/synchronizer/node"
)

type EventSync struct {
	synchronizer   *synchronizer.Synchronizer
	eventProcessor *event.EventProcessor

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewEventSync(ctx context.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*EventSync, error) {
	ethClient, err := node.DialEthClient(ctx, cfg.Chain.ChainRpcUrl)
	if err != nil {
		return nil, err
	}

	db, err := database.NewDB(ctx, cfg.MasterDB)
	if err != nil {
		log.Error("init database fail", err)
		return nil, err
	}

	syncer, err := synchronizer.NewSynchronizer(cfg, db, ethClient, shutdown)
	if err != nil {
		log.Error("new synchronizer fail", "err", err)
		return nil, err
	}

	eventConfig := &event.EventProcessorConfig{
		LoopInterval:    cfg.Chain.LoopInterval,
		EventStartBlock: cfg.Chain.StartingHeight,
		EventBlockStep:  cfg.Chain.BlockStep,
	}

	eventProcessor, err := event.NewEventProcessor(db, eventConfig, shutdown)
	if err != nil {
		log.Error("new event processor fail", "err", err)
		return nil, err
	}

	out := &EventSync{
		synchronizer:   syncer,
		eventProcessor: eventProcessor,
		shutdown:       shutdown,
	}
	return out, nil
}

func (es *EventSync) Start(ctx context.Context) error {
	err := es.synchronizer.Start()
	if err != nil {
		return err
	}
	err = es.eventProcessor.Start()
	if err != nil {
		return err
	}
	return nil
}

func (es *EventSync) Stop(ctx context.Context) error {
	err := es.synchronizer.Close()
	if err != nil {
		return err
	}
	err = es.eventProcessor.Close()
	if err != nil {
		return err
	}
	return nil
}

func (es *EventSync) Stopped() bool {
	return es.stopped.Load()
}
