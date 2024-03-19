package main

import (
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/internal/daemon/client"
	"cslab.ece.ntua.gr/actimanager/internal/daemon/cpupinning"
	"cslab.ece.ntua.gr/actimanager/internal/daemon/topology"
	clients "cslab.ece.ntua.gr/actimanager/internal/pkg/client"
	nctutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/nodecputopology"
	"flag"
	"github.com/go-logr/logr"
	"k8s.io/klog/v2/textlogger"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	var nodeName = flag.String("node-name", "minikube-m02", "Name of the node")
	var runtime = flag.String("container-runtime", "docker", "Container Runtime (Default: containerd, Values: containerd, docker, kind)")
	var cgroupsPath = flag.String("cgroups-path", "/sys/fs/cgroup/", "Specify Path to cgroups")
	var driver = flag.String("cgroups-driver", "systemd", "Set cgroups driver used by kubelet. Values: systemd, cgroupfs")
	flag.Parse()

	var logger = createLogger()
	defer func() {
		err := recover()
		if err != nil {
			logger.Info("Fatal error", "value", err)
		}
	}()

	logger.Info(
		"args",
		"runtime", runtime,
		"nodeName", nodeName,
		"driver", driver,
		"cgroupsPath", cgroupsPath,
	)

	var err error
	stopChannel := make(chan struct{})
	coreClient, err := clients.NewClient()
	cslabClient, err := clients.NewCSLabClient()
	podCPUBindingClient, err := client.NewPodCPUBindingClient(*cslabClient, logger)
	podClient, err := client.NewPodClient(*coreClient, logger)
	if err != nil {
		logger.Error(err, "cannot create clients")
		os.Exit(1)
	}
	err = podCPUBindingClient.Start(&stopChannel)
	err = podClient.Start(&stopChannel)
	if err != nil {
		logger.Error(err, "cannot start the clients")
		os.Exit(1)
	}
	cpuTopology, err := getCPUTopology()
	if err != nil {
		logger.Error(err, "cannot get cpu topology")
		os.Exit(1)
	}
	logger.Info("cpu topology", "topology", cpuTopology)
	containerRuntime := cpupinning.ParseRuntime(*runtime)
	cgroupsDriver := cpupinning.ParseCgroupsDriver(*driver)
	cpuPinningController, err := cpupinning.NewCPUPinningController(
		containerRuntime,
		cgroupsDriver,
		*cgroupsPath,
		podCPUBindingClient,
		podClient,
		*cpuTopology,
		*nodeName,
		logger,
	)
	if err != nil {
		logger.Error(err, "cannot create cpu pinnning controller")
		os.Exit(1)
	}
	cpuPinningServer := cpupinning.NewCPUPinningServer(cpuPinningController)
	topologyServer := topology.NewTopologyServer()
	daemonServer := NewDaemonServer(":8089", cpuPinningServer, topologyServer, logger)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	err = daemonServer.Start()
	if err != nil {
		logger.Error(err, "error starting the server")
		os.Exit(1)
	}

	// Graceful shutdown on SIGINT and SIGTERM
	<-signalCh
	logger.Info("Received signal, shutting down")
	daemonServer.Stop()
	podCPUBindingClient.Stop()
	podClient.Stop()
}

func createLogger() logr.Logger {
	flags := flag.NewFlagSet("klog", flag.ContinueOnError)
	config := textlogger.NewConfig(textlogger.Verbosity(3))
	config.AddFlags(flags)
	return textlogger.NewLogger(config)
}

func getCPUTopology() (*v1alpha1.CPUTopology, error) {
	output, err := exec.Command(topology.LscpuCommand, topology.LscpuArgs...).CombinedOutput()
	t, err := topology.ParseTopology(string(output))
	if err != nil {
		return nil, err
	}
	return nctutils.TopologyToV1Alpha1(t), nil
}
