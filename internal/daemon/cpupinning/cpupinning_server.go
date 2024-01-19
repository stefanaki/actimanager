package cpupinning

import (
	"context"
	"fmt"
	"strings"
)

// Server represents the CPU pinning server.
type Server struct {
	Controller *CpuPinningController
	UnimplementedCpuPinningServer
}

// NewCpuPinningServer creates a new instance of the CPU pinning server.
func NewCpuPinningServer(controller *CpuPinningController) *Server {
	return &Server{Controller: controller}
}

// ApplyPinning applies CPU pinning based on the provided request.
func (s Server) ApplyPinning(ctx context.Context, request *ApplyPinningRequest) (*Response, error) {
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
			return &Response{
				Status: ResponseStatus_ERROR,
			}, fmt.Errorf("failed to apply CPU pinning: %v", err.Error())
		}
	}

	return &Response{
		Status: ResponseStatus_SUCCESSFUL,
	}, nil
}

// RemovePinning removes the CPU pinning configuration.
func (s Server) RemovePinning(ctx context.Context, request *RemovePinningRequest) (*Response, error) {
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
			return &Response{
				Status: ResponseStatus_ERROR,
			}, fmt.Errorf("failed to remove CPU pinning: %v", err.Error())
		}
	}

	return &Response{
		Status: ResponseStatus_SUCCESSFUL,
	}, nil
}
