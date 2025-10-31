package main

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelInfo, true)))
	app := NewCli()

}
