package main

import (
	"cslab.ece.ntua.gr/actimanager/internal/daemon/cpupinning"
	"flag"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	"net"
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
		"ContainerInfo Runtime (Default: containerd, Possible values: containerd, docker, kind)",
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

	cR := parseRuntime(runtime)
	driver := parseCGroupDriver(cgroupDriver)

	cpuPinningController, err := cpupinning.NewCpuPinningController(cR, driver, cgroupPath, logger)
	if err != nil {
		klog.Fatalf("cannot create cpu pinnning controller: %v", err.Error())
	}

	cpuPinningServer := cpupinning.NewCpuPinningServer(cpuPinningController)
	healthServer := health.NewServer()

	srv := grpc.NewServer()

	cpupinning.RegisterCpuPinningServer(srv, cpuPinningServer)
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
