package main

import (
	"context"

	"github.com/Sandwichzzy/event-sync-go/common/opio"
	"github.com/Sandwichzzy/event-sync-go/config"
	"github.com/Sandwichzzy/event-sync-go/database"
	"github.com/ethereum/go-ethereum/log"

	"github.com/urfave/cli/v2"

	"github.com/Sandwichzzy/event-sync-go/common/cliapp"

	flags2 "github.com/Sandwichzzy/event-sync-go/flags"
)

func runIndexer(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	//log.Info("run event watcher indexer")
	//cfg, err := config.LoadConfig(ctx)
	//if err != nil {
	//	log.Error("failed to load config", "err", err)
	//	return nil, err
	//}
	//return event_sync.NewEventSync(ctx.Context, &cfg, shutdown)
	return nil, nil
}

func runApi(ctx *cli.Context, _ context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	return nil, nil
}

func runMigrations(ctx *cli.Context) error {
	log.Info("Running migrations...")
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return err
	}
	ctx.Context = opio.CancelOnInterrupt(ctx.Context)
	db, err := database.NewDB(ctx.Context, cfg.MasterDB)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return err
	}
	defer func(db *database.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)
	return db.ExecuteSQLMigration(cfg.Migrations)
}

func NewCli() *cli.App {
	flags := flags2.Flags
	return &cli.App{
		Version:              "v0.0.1",
		Description:          "An Event Sync Service",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "api",
				Flags:       flags,
				Description: "Runs the api service",
				Action:      cliapp.LifecycleCmd(runApi),
			},
			{
				Name:        "index",
				Flags:       flags,
				Description: "Runs the indexing service",
				Action:      cliapp.LifecycleCmd(runIndexer),
			},
			{
				Name:        "migrate",
				Flags:       flags,
				Description: "Runs the database migrations",
				Action:      runMigrations,
			},
			{
				Name:        "version",
				Description: "print version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}

}
