package podcpubinding

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/internal/daemon/cpupinning"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// applyCpuPinning applies CPU pinning for a given pod on a specified node
func (r *PodCpuBindingReconciler) applyCpuPinning(
	ctx context.Context,
	cpuSet []v1alpha1.Cpu,
	pod *corev1.Pod) error {
	logger := log.FromContext(ctx).WithName("apply-pinning")

	nodeAddress, err := r.getNodeAddress(ctx, pod.Spec.NodeName)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(fmt.Sprintf("%v:8089", nodeAddress), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %v", err.Error())
	}

	defer conn.Close()

	cpuPinningClient := cpupinning.NewCpuPinningClient(conn)
	applyCpuPinningRequest := &cpupinning.ApplyPinningRequest{
		Pod:    parsePodInfo(pod),
		CpuSet: &cpupinning.CpuSet{Cpu: convertCpuListToInt32(cpuSet)},
	}
	logger.Info("dispatching cpu pinning request", "applyCpuPinningRequest", applyCpuPinningRequest)

	res, err := cpuPinningClient.ApplyPinning(ctx, applyCpuPinningRequest)
	if err != nil {
		return fmt.Errorf("failed to apply CPU pinning: %v", err.Error())
	}

	if res.Status == cpupinning.ResponseStatus_ERROR {
		return fmt.Errorf("failed to apply CPU pinning: unknown error")
	}

	return nil
}

// removeCpuPinning removes CPU pinning for a given pod on a specified node
func (r *PodCpuBindingReconciler) removeCpuPinning(
	ctx context.Context,
	pod *corev1.Pod) error {
	logger := log.FromContext(ctx).WithName("remove-pinning")

	nodeAddress, err := r.getNodeAddress(ctx, pod.Spec.NodeName)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(fmt.Sprintf("%v:8089", nodeAddress), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	cpuPinningClient := cpupinning.NewCpuPinningClient(conn)
	removeCpuPinningRequest := &cpupinning.RemovePinningRequest{Pod: parsePodInfo(pod)}
	logger.Info("trying to remove cpu pinning", "removeCpuPinningRequest", removeCpuPinningRequest)

	res, err := cpuPinningClient.RemovePinning(ctx, removeCpuPinningRequest)
	if err != nil {
		return fmt.Errorf("failed to remove CPU pinning: %v", err.Error())
	}

	if res.Status == cpupinning.ResponseStatus_ERROR {
		return fmt.Errorf("failed to remove CPU pinning: unknown error")
	}

	return nil
}

// parsePodInfo extracts relevant information from a Pod to create a cpupinning.Pod object
func parsePodInfo(pod *corev1.Pod) *cpupinning.Pod {
	p := &cpupinning.Pod{
		Id:         string(pod.ObjectMeta.UID),
		Name:       pod.Name,
		Namespace:  pod.Namespace,
		Containers: nil,
	}

	containers := make([]*cpupinning.Container, 0)
	for _, containerStatus := range pod.Status.ContainerStatuses {
		containers = append(containers, &cpupinning.Container{
			Id:        containerStatus.ContainerID,
			Name:      containerStatus.Name,
			Resources: parseContainerResources(containerStatus.Name, pod),
		})
	}

	p.Containers = containers

	return p
}

// parseContainerResources extracts resource information from a container
func parseContainerResources(containerName string, pod *corev1.Pod) *cpupinning.ResourceInfo {
	resources := &cpupinning.ResourceInfo{}

	for _, container := range pod.Spec.Containers {
		if container.Name == containerName {
			limitCpus := container.Resources.Limits.Cpu()
			limitMemory := container.Resources.Limits.Memory()
			requestCpus := container.Resources.Requests.Cpu()
			requestMemory := container.Resources.Requests.Memory()

			resources = &cpupinning.ResourceInfo{
				RequestedCpus:   int32(requestCpus.MilliValue()),
				LimitCpus:       int32(limitCpus.MilliValue()),
				RequestedMemory: []byte(requestMemory.String()),
				LimitMemory:     []byte(limitMemory.String()),
			}

			return resources
		}
	}

	return resources
}

// convertCpuListToInt32 maps a Cpu list to an int32 array
func convertCpuListToInt32(cpuSet []v1alpha1.Cpu) []int32 {
	var cpuList []int32

	for _, cpu := range cpuSet {
		cpuList = append(cpuList, int32(cpu.CpuId))
	}

	return cpuList
}
