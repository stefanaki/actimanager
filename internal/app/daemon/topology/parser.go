package topology

import (
	pbtopo "cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/topology"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// ParseTopology parses the output of lscpu command
func ParseTopology(output string) (*pbtopo.TopologyResponse, error) {
	res := &pbtopo.TopologyResponse{}
	res.Sockets = make(map[string]*pbtopo.Socket)
	res.NumaNodes = make(map[string]*pbtopo.NUMANode)
	res.Cpus = make([]int64, 0)

	for _, lsLine := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		if strings.HasPrefix(lsLine, "#") {
			continue
		}
		values := strings.Split(lsLine, ",")
		if len(values) != 4 {
			fmt.Printf("Invalid format for socket,node,core,cpu: %s\n", lsLine)
			return nil, fmt.Errorf("invalid format for socket,node,core,cpu: %s", lsLine)
		}

		var err error
		socketIdStr, nodeIdStr, coreIdStr, cpuIdStr := values[0], values[1], values[2], values[3]
		socketId, err := strconv.Atoi(socketIdStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse socket ID %q: %v", socketIdStr, err)
		}
		nodeId, err := strconv.Atoi(nodeIdStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse node ID %q: %v", nodeIdStr, err)
		}
		coreId, err := strconv.Atoi(coreIdStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse core ID %q: %v", coreIdStr, err)
		}
		cpuId, err := strconv.Atoi(cpuIdStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse cpu ID %q: %v", cpuIdStr, err)
		}
		socketIdStr, nodeIdStr, coreIdStr = strconv.Itoa(socketId), strconv.Itoa(nodeId), strconv.Itoa(coreId)
		if _, ok := res.Sockets[socketIdStr]; !ok {
			res.Sockets[socketIdStr] = &pbtopo.Socket{
				Id:    int64(socketId),
				Cores: make(map[string]*pbtopo.Core),
			}
		}
		if _, ok := res.NumaNodes[nodeIdStr]; !ok {
			res.NumaNodes[nodeIdStr] = &pbtopo.NUMANode{
				Id:   int64(nodeId),
				Cpus: make([]int64, 0),
			}
		}
		if _, ok := res.Sockets[socketIdStr].Cores[coreIdStr]; !ok {
			res.Sockets[socketIdStr].Cores[coreIdStr] = &pbtopo.Core{
				Id:   int64(coreId),
				Cpus: make([]int64, 0),
			}
		}
		if !slices.Contains(res.Sockets[socketIdStr].Cores[coreIdStr].Cpus, int64(cpuId)) {
			res.Cpus = append(res.Cpus, int64(cpuId))
			res.Sockets[socketIdStr].Cores[coreIdStr].Cpus = append(res.Sockets[socketIdStr].Cores[coreIdStr].Cpus, int64(cpuId))
		}
		if !slices.Contains(res.NumaNodes[nodeIdStr].Cpus, int64(cpuId)) {
			res.NumaNodes[nodeIdStr].Cpus = append(res.NumaNodes[nodeIdStr].Cpus, int64(cpuId))
		}
	}
	return res, nil
}
