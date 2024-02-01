package main

import (
	cpupinningserver "cslab.ece.ntua.gr/actimanager/internal/pkg/cpupinning"
	"flag"
	"net"

	"cslab.ece.ntua.gr/actimanager/internal/daemon/cpupinning"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/textlogger"
)

func main() {
	var cgroupPath, cgroupDriver, nodeName, runtime string
	var logger = createLogger()

	defer func() {
		err := recover()
		if err != nil {
			logger.Info("Fatal error", "value", err)
		}
	}()

	flag.StringVar(
		&runtime,
		"runtime",
		"docker",
		"Container Runtime (Default: containerd, Possible values: containerd, docker, kind)",
	)
	flag.StringVar(&cgroupPath, "cpath", "/sys/fs/cgroup/", "Specify Path to cgroups")
	flag.StringVar(&nodeName, "node-name", "", "Node name")
	flag.StringVar(&cgroupDriver, "cgroup-driver", "systemd", "Set cgroup driver used by kubelet. Values: systemd, cgroupfs")
	flag.Parse()

	logger.Info(
		"args",
		"runtime", runtime,
		"nodeName", nodeName,
		"cgroupDriver", cgroupDriver,
		"cgroupPath", cgroupPath,
	)

	cR := cpupinning.ParseRuntime(runtime)
	driver := cpupinning.ParseCgroupsDriver(cgroupDriver)

	cpuPinningController, err := cpupinning.NewCpuPinningController(cR, driver, cgroupPath, logger)
	if err != nil {
		klog.Fatalf("cannot create cpu pinnning controller: %v", err.Error())
	}

	cpuPinningServer := cpupinning.NewCpuPinningServer(cpuPinningController)
	healthServer := health.NewServer()
	srv := grpc.NewServer()

	cpupinningserver.RegisterCpuPinningServer(srv, cpuPinningServer)
	healthv1.RegisterHealthServer(srv, healthServer)

	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		klog.Fatalf("cannot create tcp listener: %v", err.Error())
	}

	klog.Infof("server listening at %v\n", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		klog.Fatalf("cannot serve: %v", err.Error())
	}
}

func createLogger() logr.Logger {
	flags := flag.NewFlagSet("klog", flag.ContinueOnError)

	config := textlogger.NewConfig(textlogger.Verbosity(3))
	config.AddFlags(flags)

	return textlogger.NewLogger(config)
}
