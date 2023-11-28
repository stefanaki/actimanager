package daemon

import (
	"context"
)

type Server struct {
	UnimplementedCpuPinningDaemonServer
}

func (s Server) ApplyPinning(ctx context.Context, request *ApplyPinningRequest) (*ApplyPinningResponse, error) {

	return &ApplyPinningResponse{}, nil
}

func (s Server) UpdatePinning(ctx context.Context, request *UpdatePinningRequest) (*UpdatePinningResponse, error) {
	return &UpdatePinningResponse{}, nil
}

func (s Server) RemovePinning(ctx context.Context, request *RemovePinningRequest) (*RemovePinningResponse, error) {
	println("ok")
	return &RemovePinningResponse{}, nil
}
