package flags

import (
	"time"

	"github.com/urfave/cli/v2"
)

const envVarPrefix = "EVENT_SYNC"

func prefixEnvVars(name string) []string {
	return []string{envVarPrefix + "_" + name}
}

var (
	MigrationsFlag = &cli.StringFlag{
		Name:    "migrations-dir",
		Value:   "./migrations",
		Usage:   "path to migrations folder",
		EnvVars: prefixEnvVars("MIGRATIONS_DIR"),
	}

	ChainIdFlag = &cli.UintFlag{
		Name:     "chain-id",
		Usage:    "The port of the api",
		EnvVars:  prefixEnvVars("CHAIN_ID"),
		Value:    1,
		Required: true,
	}

	ChainRpcFlag = &cli.StringFlag{
		Name:     "chain-rpc",
		Usage:    "HTTP provider URL for L1",
		EnvVars:  prefixEnvVars("CHAIN_RPC"),
		Required: true,
	}
	//同步起始区块高度
	StartingHeightFlag = &cli.Uint64Flag{
		Name:    "starting-height",
		Usage:   "The starting height of chain",
		EnvVars: prefixEnvVars("STARTING_HEIGHT"),
		Value:   0,
	}
	// 交易确认数
	ConfirmationsFlag = &cli.Uint64Flag{
		Name:    "confirmations",
		Usage:   "The confirmation depth of l1",
		EnvVars: prefixEnvVars("CONFIRMATIONS"),
		Value:   64,
	}
	// 同步循环间隔
	LoopIntervalFlag = &cli.DurationFlag{
		Name:    "loop-interval",
		Usage:   "The interval of synchronization",
		EnvVars: prefixEnvVars("LOOP_INTERVAL"),
		Value:   time.Second * 5,
	}
	BlocksStepFlag = &cli.UintFlag{
		Name:    "blocks-step",
		Usage:   "Scanner blocks step",
		EnvVars: prefixEnvVars("BLOCKS_STEP"),
		Value:   5,
	}
	//主数据库连接配置
	// MasterDbHostFlag MasterDb Flags
	MasterDbHostFlag = &cli.StringFlag{
		Name:     "master-db-host",
		Usage:    "The host of the master database",
		EnvVars:  prefixEnvVars("MASTER_DB_HOST"),
		Required: true,
	}
	MasterDbPortFlag = &cli.IntFlag{
		Name:     "master-db-port",
		Usage:    "The port of the master database",
		EnvVars:  prefixEnvVars("MASTER_DB_PORT"),
		Required: true,
	}
	MasterDbUserFlag = &cli.StringFlag{
		Name:     "master-db-user",
		Usage:    "The user of the master database",
		EnvVars:  prefixEnvVars("MASTER_DB_USER"),
		Required: true,
	}
	MasterDbPasswordFlag = &cli.StringFlag{
		Name:     "master-db-password",
		Usage:    "The host of the master database",
		EnvVars:  prefixEnvVars("MASTER_DB_PASSWORD"),
		Required: true,
	}
	MasterDbNameFlag = &cli.StringFlag{
		Name:     "master-db-name",
		Usage:    "The db name of the master database",
		EnvVars:  prefixEnvVars("MASTER_DB_NAME"),
		Required: true,
	}

	HttpHostFlag = &cli.StringFlag{
		Name:     "http-host",
		Usage:    "The host of the api",
		EnvVars:  prefixEnvVars("HTTP_HOST"),
		Required: true,
	}
	HttpPortFlag = &cli.IntFlag{
		Name:     "http-port",
		Usage:    "The port of the api",
		EnvVars:  prefixEnvVars("HTTP_PORT"),
		Value:    8987,
		Required: true,
	}

	GrpcHostFlag = &cli.StringFlag{
		Name:    "grpc-host",
		Usage:   "The host of the api",
		EnvVars: prefixEnvVars("GRPC_HOST"),
		Value:   "127.0.0.1",
	}
	GrpcPortFlag = &cli.IntFlag{
		Name:    "grpc-port",
		Usage:   "The port of the api",
		EnvVars: prefixEnvVars("GRPC_PORT"),
		Value:   8987,
	}

	SlaveDbEnableFlag = &cli.BoolFlag{
		Name:     "slave-db-enable",
		Usage:    "Whether to use slave db",
		EnvVars:  prefixEnvVars("SLAVE_DB_ENABLE"),
		Required: true,
	}

	// SlaveDbHostFlag Slave DB  flags
	SlaveDbHostFlag = &cli.StringFlag{
		Name:    "slave-db-host",
		Usage:   "The host of the slave database",
		EnvVars: prefixEnvVars("SLAVE_DB_HOST"),
	}
	SlaveDbPortFlag = &cli.IntFlag{
		Name:    "slave-db-port",
		Usage:   "The port of the slave database",
		EnvVars: prefixEnvVars("SLAVE_DB_PORT"),
	}
	SlaveDbUserFlag = &cli.StringFlag{
		Name:    "slave-db-user",
		Usage:   "The user of the slave database",
		EnvVars: prefixEnvVars("SLAVE_DB_USER"),
	}
	SlaveDbPasswordFlag = &cli.StringFlag{
		Name:    "slave-db-password",
		Usage:   "The host of the slave database",
		EnvVars: prefixEnvVars("SLAVE_DB_PASSWORD"),
	}
	SlaveDbNameFlag = &cli.StringFlag{
		Name:    "slave-db-name",
		Usage:   "The db name of the slave database",
		EnvVars: prefixEnvVars("SLAVE_DB_NAME"),
	}
)

var requiredFlags = []cli.Flag{
	MigrationsFlag,
	ChainIdFlag,
	ChainRpcFlag,
	MasterDbHostFlag,
	MasterDbPortFlag,
	MasterDbUserFlag,
	MasterDbPasswordFlag,
	MasterDbNameFlag,
	LoopIntervalFlag,
	BlocksStepFlag,
	SlaveDbEnableFlag,
	HttpHostFlag,
	HttpPortFlag,
}

var optionalFlags = []cli.Flag{
	StartingHeightFlag,
	ConfirmationsFlag,
	SlaveDbHostFlag,
	SlaveDbPortFlag,
	SlaveDbUserFlag,
	SlaveDbPasswordFlag,
	SlaveDbNameFlag,
	GrpcHostFlag,
	GrpcPortFlag,
}

var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
}
