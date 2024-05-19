package cpupinning

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/cpupinning"
	"fmt"
)

// Server represents the CPU pinning server.
type Server struct {
	Controller *CPUPinningController
	cpupinning.UnimplementedCPUPinningServer
}

// NewCPUPinningServer creates a new instance of the CPU pinning server.
func NewCPUPinningServer(controller *CPUPinningController) *Server {
	return &Server{Controller: controller}
}

// ApplyPinning applies the requested CPU pinning on the containers of a Pod.
func (s Server) ApplyPinning(ctx context.Context, request *cpupinning.ApplyPinningRequest) (*cpupinning.Response, error) {
	pod := request.GetPod()
	cpuSet := ConvertIntSliceToString(request.CpuSet)
	memSet := ConvertIntSliceToString(request.MemSet)

	if err := s.Controller.Apply(pod, cpuSet, memSet); err != nil {
		return &cpupinning.Response{
			Status: cpupinning.ResponseStatus_ERROR,
		}, fmt.Errorf("failed to apply CPU pinning: %v", err.Error())
	}

	return &cpupinning.Response{
		Status: cpupinning.ResponseStatus_SUCCESSFUL,
	}, nil
}

// RemovePinning removes the CPU pinning of the containers of a Pod.
func (s Server) RemovePinning(ctx context.Context, request *cpupinning.RemovePinningRequest) (*cpupinning.Response, error) {
	if err := s.Controller.Remove(request.GetPod()); err != nil {
		return &cpupinning.Response{
			Status: cpupinning.ResponseStatus_ERROR,
		}, fmt.Errorf("failed to remove CPU pinning: %v", err.Error())
	}

	return &cpupinning.Response{
		Status: cpupinning.ResponseStatus_SUCCESSFUL,
	}, nil
}
