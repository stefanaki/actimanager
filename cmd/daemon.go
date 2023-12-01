package main

import (
	"cslab.ece.ntua.gr/actimanager/internal/daemon/cpupinning"
	"flag"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	"net"
)

func parseRuntime(runtime string) cpupinning.ContainerRuntime {
	val, ok := map[string]cpupinning.ContainerRuntime{
		"containerd": cpupinning.ContainerdRunc,
		"kind":       cpupinning.Kind,
		"docker":     cpupinning.Docker,
	}[runtime]
	if !ok {
		klog.Fatalf("unknown runtime %s", runtime)
	}
	return val
}

func parseCGroupDriver(driver string) cpupinning.CGroupDriver {
	val, ok := map[string]cpupinning.CGroupDriver{
		"systemd":  cpupinning.DriverSystemd,
		"cgroupfs": cpupinning.DriverCgroupfs,
	}[driver]
	if !ok {
		klog.Fatalf("unknown cgroup driver %s", driver)
	}
	return val
}

func createLogger() logr.Logger {
	flags := flag.NewFlagSet("klog", flag.ContinueOnError)
	klog.InitFlags(flags)
	_ = flags.Parse([]string{"-v", "3"})
	return klogr.NewWithOptions(klogr.WithFormat(klogr.FormatKlog))
}

func RunDaemon() {
	var cgroupPath, cgroupDriver, nodeName, runtime string
	var logger = createLogger()

	flag.StringVar(
		&runtime,
		"runtime",
		"containerd",
		"MyContainer Runtime (Default: containerd, Possible values: containerd, docker, kind)",
	)
	flag.StringVar(&cgroupPath, "cpath", "/sys/fs/cgroup/", "Specify Path to cgroupds")
	flag.StringVar(&nodeName, "node-name", "", "Node name")
	flag.StringVar(&cgroupDriver, "cgroup-driver", "systemd", "Set cgroup driver used by kubelet. Values: systemd, cgroupfs")

	flag.Parse()

	defer func() {
		err := recover()
		if err != nil {
			logger.Info("Fatal error", "value", err)
		}
	}()

	lis, err := net.Listen("tcp", ":8089")

	if err != nil {
		klog.Fatalf("cannot create tcp listener: %v", err.Error())
	}

	logger.Info(
		"args",
		"runtime", runtime,
		"nodeName", nodeName,
		"cgroupDriver", cgroupDriver,
		"cgroupPath", cgroupPath,
	)

	cR := parseRuntime(runtime)
	driver := parseCGroupDriver(cgroupDriver)

	cpuPinningController, err := cpupinning.NewCpuPinningController(cR, driver, cgroupPath, logger)
	if err != nil {
		klog.Fatalf("cannot create cpu pinnning controller: %v", err.Error())
	}

	cpuPinningServer := &cpupinning.Server{Controller: cpuPinningController}
	srv := grpc.NewServer()

	cpupinning.RegisterCpuPinningServer(srv, cpuPinningServer)

	klog.Infof("server listening at %v\n", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		klog.Fatalf("cannot serve: %v", err.Error())
	}
}
