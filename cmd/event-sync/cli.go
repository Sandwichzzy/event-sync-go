package main

import (
	"context"

	"github.com/Sandwichzzy/event-sync-go.git/common/cliapp"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

func runIndexer(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("run event")
	return nil, nil
}

func runApi(ctx *cli.Context, _ context.CancelCauseFunc) (cliapp.Lifecycle, error) {

}

func runMigrations(ctx *cli.Context) error {

}

func NewCli() *cli.App {

}
