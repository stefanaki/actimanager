package main

import (
	"cslab.ece.ntua.gr/actimanager/internal/daemon"
	"google.golang.org/grpc"
	"log"
	"net"
)

func runDaemon() {
	lis, err := net.Listen("tcp", ":8089")

	if err != nil {
		log.Fatalf("cannot create tcp listener: %v", err.Error())
	}

	serverRegistrar := grpc.NewServer()
	service := &daemon.Server{}

	daemon.RegisterCpuPinningDaemonServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)

	if err != nil {
		log.Fatalf("cannot serve: %v", err.Error())
	}
}
