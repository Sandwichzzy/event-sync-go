package config

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	"github.com/Sandwichzzy/event-sync-go/flags"
)

const (
	defaultConfirmations = 64
	defaultLoopInterval  = 5000
	TreasureManagerAddr  = "0x388fF618Ca5c1b8F28D4E845B431Ca3D4200140e"
)

type Config struct {
	Migrations    string
	Chain         ChainConfig
	MasterDB      DBConfig
	SlaveDB       DBConfig
	SlaveDbEnable bool
	//ApiCacheEnable bool
	HTTPServer ServerConfig
}

type ChainConfig struct {
	ChainRpcUrl    string
	ChainId        uint
	StartingHeight uint64
	Confirmations  uint64
	BlockStep      uint64
	Contracts      []common.Address
	LoopInterval   time.Duration
}

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type ServerConfig struct {
	Host string
	Port int
}

func LoadConfig(cliCtx *cli.Context) (Config, error) {
	var cfg Config
	cfg = NewConfig(cliCtx)

	if cfg.Chain.Confirmations == 0 {
		cfg.Chain.Confirmations = defaultConfirmations
	}
	if cfg.Chain.LoopInterval == 0 {
		cfg.Chain.LoopInterval = defaultLoopInterval
	}
	log.Info("loaded chain config", "config", cfg.Chain)
	return cfg, nil
}

func LoadContracts() []common.Address {
	var Contracts []common.Address
	Contracts = append(Contracts, common.HexToAddress(TreasureManagerAddr))
	return Contracts
}

func NewConfig(cliCtx *cli.Context) Config {
	return Config{
		Migrations: cliCtx.String(flags.MigrationsFlag.Name),
		Chain: ChainConfig{
			ChainId:        cliCtx.Uint(flags.ChainIdFlag.Name),
			ChainRpcUrl:    cliCtx.String(flags.ChainRpcFlag.Name),
			StartingHeight: cliCtx.Uint64(flags.StartingHeightFlag.Name),
			Confirmations:  cliCtx.Uint64(flags.ConfirmationsFlag.Name),
			BlockStep:      cliCtx.Uint64(flags.BlocksStepFlag.Name),
			Contracts:      LoadContracts(),
			LoopInterval:   cliCtx.Duration(flags.LoopIntervalFlag.Name),
		},
		MasterDB: DBConfig{
			Host:     cliCtx.String(flags.MasterDbHostFlag.Name),
			Port:     cliCtx.Int(flags.MasterDbPortFlag.Name),
			Name:     cliCtx.String(flags.MasterDbNameFlag.Name),
			User:     cliCtx.String(flags.MasterDbUserFlag.Name),
			Password: cliCtx.String(flags.MasterDbPasswordFlag.Name),
		},
		SlaveDB: DBConfig{
			Host:     cliCtx.String(flags.SlaveDbHostFlag.Name),
			Port:     cliCtx.Int(flags.SlaveDbPortFlag.Name),
			Name:     cliCtx.String(flags.SlaveDbNameFlag.Name),
			User:     cliCtx.String(flags.SlaveDbUserFlag.Name),
			Password: cliCtx.String(flags.SlaveDbPasswordFlag.Name),
		},
		SlaveDbEnable: cliCtx.Bool(flags.SlaveDbEnableFlag.Name),
		HTTPServer: ServerConfig{
			Host: cliCtx.String(flags.HttpHostFlag.Name),
			Port: cliCtx.Int(flags.HttpPortFlag.Name),
		},
	}
}
