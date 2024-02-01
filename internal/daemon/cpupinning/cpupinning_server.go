package cpupinning

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/cpupinning"
	"fmt"
	"strings"
)

// Server represents the CPU pinning server.
type Server struct {
	Controller *CpuPinningController
	cpupinning.UnimplementedCpuPinningServer
}

// NewCpuPinningServer creates a new instance of the CPU pinning server.
func NewCpuPinningServer(controller *CpuPinningController) *Server {
	return &Server{Controller: controller}
}

// ApplyPinning applies CPU pinning based on the provided request.
func (s Server) ApplyPinning(ctx context.Context, request *cpupinning.ApplyPinningRequest) (*cpupinning.Response, error) {
	pod := request.Pod
	cpuSet := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(request.CpuSet)), ","), "[]")
	memSet := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(request.MemSet)), ","), "[]")

	for _, container := range request.Pod.Containers {
		c := ContainerInfo{
			CID:  container.Id,
			PID:  pod.Id,
			Name: container.Name,
			QS: QoSFromLimit(
				container.Resources.LimitCpus,
				container.Resources.RequestedCpus,
				container.Resources.LimitMemory,
				container.Resources.RequestedMemory,
			),
			Cpus: container.Resources.RequestedCpus,
		}

		if err := s.Controller.Apply(c, cpuSet, memSet); err != nil {
			return &cpupinning.Response{
				Status: cpupinning.ResponseStatus_ERROR,
			}, fmt.Errorf("failed to apply CPU pinning: %v", err.Error())
		}
	}

	return &cpupinning.Response{
		Status: cpupinning.ResponseStatus_SUCCESSFUL,
	}, nil
}

// RemovePinning removes the CPU pinning configuration.
func (s Server) RemovePinning(ctx context.Context, request *cpupinning.RemovePinningRequest) (*cpupinning.Response, error) {
	pod := request.Pod

	for _, container := range request.Pod.Containers {
		c := ContainerInfo{
			CID:  container.Id,
			PID:  pod.Id,
			Name: container.Name,
			QS: QoSFromLimit(
				container.Resources.LimitCpus,
				container.Resources.RequestedCpus,
				container.Resources.LimitMemory,
				container.Resources.RequestedMemory,
			),
			Cpus: container.Resources.RequestedCpus,
		}

		if err := s.Controller.Remove(c); err != nil {
			return &cpupinning.Response{
				Status: cpupinning.ResponseStatus_ERROR,
			}, fmt.Errorf("failed to remove CPU pinning: %v", err.Error())
		}
	}

	return &cpupinning.Response{
		Status: cpupinning.ResponseStatus_SUCCESSFUL,
	}, nil
}
