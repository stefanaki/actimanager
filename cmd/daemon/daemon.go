package main

import (
	"context"
	cpupinningserver "cslab.ece.ntua.gr/actimanager/internal/daemon/cpupinning"
	topologyserver "cslab.ece.ntua.gr/actimanager/internal/daemon/topology"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/cpupinning"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/topology"
	"fmt"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"time"
)

type DaemonServer struct {
	grpcServer       *grpc.Server
	endpoint         string
	cpuPinningServer *cpupinningserver.Server
	topologyServer   *topologyserver.Server
	logger           logr.Logger
}

func NewDaemonServer(endpoint string,
	cpuPinningServer *cpupinningserver.Server,
	topologyServer *topologyserver.Server,
	logger logr.Logger) *DaemonServer {
	return &DaemonServer{
		grpcServer:       nil,
		endpoint:         endpoint,
		cpuPinningServer: cpuPinningServer,
		topologyServer:   topologyServer,
		logger:           logger.WithName("daemon-server"),
	}
}

func (s *DaemonServer) Start() error {
	s.logger.Info("Starting daemon server")
	lis, err := net.Listen("tcp", s.endpoint)
	if err != nil {
		return fmt.Errorf("cannot create tcp listener: %v", err.Error())
	}
	s.grpcServer = grpc.NewServer()

	// Register servers below
	cpupinning.RegisterCPUPinningServer(s.grpcServer, *s.cpuPinningServer)
	topology.RegisterTopologyServer(s.grpcServer, *s.topologyServer)
	healthv1.RegisterHealthServer(s.grpcServer, health.NewServer())

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Error(err, "failed to serve")
		}
	}()
	conn, err := grpc.DialContext(context.Background(), s.endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithIdleTimeout(5*time.Second),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := &net.Dialer{}
			return d.DialContext(ctx, "tcp", addr)
		}),
	)
	if err != nil {
		return fmt.Errorf("cannot connect to server: %v", err.Error())
	}
	s.logger.Info("Daemon server started serving", "endpoint", s.endpoint)
	defer conn.Close()
	return nil
}

func (s *DaemonServer) Stop() {
	s.logger.Info("Stopping daemon server")
	s.grpcServer.GracefulStop()
	s.grpcServer = nil
}
