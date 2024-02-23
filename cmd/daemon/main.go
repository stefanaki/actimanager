package main

import (
	"cslab.ece.ntua.gr/actimanager/internal/daemon/cpupinning"
	"flag"
	"github.com/go-logr/logr"
	"k8s.io/klog/v2/textlogger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var nodeName = flag.String("node-name", "minikube", "Name of the node")
	var runtime = flag.String("container-runtime", "docker", "Container Runtime (Default: containerd, Values: containerd, docker, kind)")
	var cgroupsPath = flag.String("cgroups-path", "/sys/fs/cgroup/", "Specify Path to cgroups")
	var driver = flag.String("cgroups-driver", "systemd", "Set cgroup cgroupsDriver used by kubelet. Values: systemd, cgroupfs")
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

	containerRuntime := cpupinning.ParseRuntime(*runtime)
	cgroupsDriver := cpupinning.ParseCgroupsDriver(*driver)

	cpuPinningController, err := cpupinning.NewCpuPinningController(containerRuntime, cgroupsDriver, *cgroupsPath, logger)
	if err != nil {
		logger.Error(err, "cannot create cpu pinnning controller")
		os.Exit(1)
	}
	cpuPinningServer := cpupinning.NewCpuPinningServer(cpuPinningController)
	daemonServer := NewDaemonServer(":8089", cpuPinningServer, logger)

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
}

func createLogger() logr.Logger {
	flags := flag.NewFlagSet("klog", flag.ContinueOnError)
	config := textlogger.NewConfig(textlogger.Verbosity(3))
	config.AddFlags(flags)
	return textlogger.NewLogger(config)
}
