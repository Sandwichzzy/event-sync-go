package grpc

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Sandwichzzy/event-sync-go/config"
	"github.com/Sandwichzzy/event-sync-go/database"
	"github.com/Sandwichzzy/event-sync-go/services/grpc/eventpb"
)

const MaxReceivedMessageSize = 1024 * 1024 * 30000

type RpcService struct {
	conf *config.Config
	db   *database.DB
	eventpb.UnimplementedEventServiceServer
	stopped atomic.Bool
}

func (rs *RpcService) Stop(ctx context.Context) error {
	rs.stopped.Store(true)
	return nil
}

func (rs *RpcService) Stopped() bool {
	return rs.stopped.Load()
}

func NewRpcService(conf *config.Config, db *database.DB) (*RpcService, error) {
	rpcService := &RpcService{
		conf: conf,
		db:   db,
	}
	return rpcService, nil
}

func (rs *RpcService) start(ctx context.Context) error {
	go func(s *RpcService) {
		addr := fmt.Sprintf("%s:%d", rs.conf.GrpcServer.Host, rs.conf.GrpcServer.Port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("Could not start grpc server on ", addr)
			return
		}

		gs := grpc.NewServer(
			grpc.MaxRecvMsgSize(MaxReceivedMessageSize),
			grpc.ChainUnaryInterceptor(nil))

		reflection.Register(gs) // grpcui -plaintext 127.0.0.1:port

		eventpb.RegisterEventServiceServer(gs, rs)

		log.Info("Grpc info", "Port", rs.conf.GrpcServer.Port, "addr", listener.Addr())

		if err := gs.Serve(listener); err != nil {
			log.Error("Could not GRPC services")
		}

	}(rs)
	return nil
}
