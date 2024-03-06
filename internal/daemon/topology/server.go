package topology

import (
	"context"
	topo "cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/topology"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"os/exec"
)

var (
	LscpuCommand = "/usr/bin/lscpu"
	LscpuArgs    = []string{"-p=socket,node,core,cpu"}
)

type Server struct {
	topo.UnimplementedTopologyServer
}

func NewTopologyServer() *Server {
	return &Server{}
}

// GetTopology returns the topology of the node
func (s Server) GetTopology(ctx context.Context, in *emptypb.Empty) (*topo.TopologyResponse, error) {
	output, err := exec.Command(LscpuCommand, LscpuArgs...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing lscpu: %v", err)
	}
	return ParseTopology(string(output))
}
