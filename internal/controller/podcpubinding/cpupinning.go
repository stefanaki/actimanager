package podcpubinding

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/cpupinning"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// applyCpuPinning applies CPU pinning for a given pod on a specified node
func (r *PodCpuBindingReconciler) applyCpuPinning(
	ctx context.Context,
	cpuSet []int,
	memSet []int,
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
		Pod:    cpupinning.ParsePodInfo(pod),
		CpuSet: convertIntSliceToInt32(cpuSet),
		MemSet: convertIntSliceToInt32(memSet),
	}
	logger.Info("Requesting CPU pinning", "request", applyCpuPinningRequest)

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
	removeCpuPinningRequest := &cpupinning.RemovePinningRequest{Pod: cpupinning.ParsePodInfo(pod)}
	logger.Info("Removing CPU pinning", "request", removeCpuPinningRequest)

	res, err := cpuPinningClient.RemovePinning(ctx, removeCpuPinningRequest)
	if err != nil {
		return fmt.Errorf("failed to remove CPU pinning: %v", err.Error())
	}

	if res.Status == cpupinning.ResponseStatus_ERROR {
		return fmt.Errorf("failed to remove CPU pinning: unknown error")
	}

	return nil
}

// convertIntSliceToInt32 maps an int slice to an int32 slice
func convertIntSliceToInt32(intSlice []int) []int32 {
	var int32Slice []int32
	for _, i := range intSlice {
		int32Slice = append(int32Slice, int32(i))
	}
	return int32Slice
}
