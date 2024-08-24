package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"cslab.ece.ntua.gr/actimanager/internal/app/daemon/client"
	"cslab.ece.ntua.gr/actimanager/internal/app/daemon/cpupinning"
	"cslab.ece.ntua.gr/actimanager/internal/app/daemon/topology"
	clients "cslab.ece.ntua.gr/actimanager/internal/pkg/client"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2/textlogger"
)

func main() {
	var verbosity = flag.Int("verbosity", 3, "Log verbosity level (0 = least verbose, 5 = most verbose)")
	var nodeName = flag.String("node-name", "minikube-m02", "Name of the node")
	var runtime = flag.String("container-runtime", "docker", "Container Runtime (Default: containerd, Values: containerd, docker, kind)")
	var cgroupsPath = flag.String("cgroups-path", "/sys/fs/cgroup/", "Specify Path to cgroups")
	var driver = flag.String("cgroups-driver", "systemd", "Set cgroups driver used by kubelet. Values: systemd, cgroupfs")
	var reconcilePeriod = flag.String("reconcile-period", "10s", "Reconcile period")
	flag.Parse()

	var logger = createLogger(verbosity)
	defer handlePanic(logger)

	logger.Info("args", "runtime", runtime, "nodeName", nodeName, "driver", driver, "cgroupsPath", cgroupsPath, "reconcilePeriod", reconcilePeriod)

	coreClient, err := clients.NewClient()
	cslabClient, err := clients.NewCSLabClient()
	podCPUBindingClient, err := client.NewPodCPUBindingClient(*cslabClient, logger)
	podClient, err := client.NewPodClient(*coreClient, logger)
	if err != nil {
		logger.Error(err, "cannot create clients")
		os.Exit(1)
	}

	err = podCPUBindingClient.Start()
	err = podClient.Start()
	if err != nil {
		logger.Error(err, "cannot start the clients")
		os.Exit(1)
	}

	cpuTopology, err := topology.GetCPUTopology()
	if err != nil {
		logger.Error(err, "cannot get cpu topology")
		os.Exit(1)
	}

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
		*reconcilePeriod,
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

	// Graceful shutdown on SIGINT SIGTERM, SIGQUIT
	<-signalCh
	logger.Info("Received signal, shutting down")
	daemonServer.Stop()
	podCPUBindingClient.Stop()
	podClient.Stop()
}

func createLogger(verbosity *int) logr.Logger {
	flags := flag.NewFlagSet("klog", flag.ContinueOnError)
	config := textlogger.NewConfig(textlogger.Verbosity(*verbosity))
	config.AddFlags(flags)
	return textlogger.NewLogger(config)
}

func handlePanic(logger logr.Logger) {
	if err := recover(); err != nil {
		logger.Info("Fatal error", "value", err)
	}
}
