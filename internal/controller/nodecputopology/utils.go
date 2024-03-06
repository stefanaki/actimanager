package nodecputopology

import (
	"context"
	pbtopo "cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/topology"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	corev1 "k8s.io/api/core/v1"
)

func (r *NodeCpuTopologyReconciler) getNodeAddress(ctx context.Context, node *corev1.Node) (string, error) {
	nodeAddress := ""
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			nodeAddress = address.Address
			break
		}
	}
	if nodeAddress == "" {
		return "", fmt.Errorf("failed to get IP address of node " + node.Name)
	}
	return nodeAddress, nil
}

func (r *NodeCpuTopologyReconciler) getTopology(ctx context.Context, node *corev1.Node) (*pbtopo.TopologyResponse, error) {
	nodeAddress, err := r.getNodeAddress(ctx, node)
	if err != nil {
		return nil, fmt.Errorf("failed to get IP address of node: %v", err)
	}
	conn, err := grpc.Dial(fmt.Sprintf("%v:8089", nodeAddress), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	topologyClient := pbtopo.NewTopologyClient(conn)
	topologyResponse, err := topologyClient.GetTopology(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to get topology: %v", err)
	}
	return topologyResponse, nil
}
