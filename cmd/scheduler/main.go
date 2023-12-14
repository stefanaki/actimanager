package main

import (
	"cslab.ece.ntua.gr/actimanager/internal/scheduler/cpuaware"
	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // for rest client metric registration
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"os"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(cpuaware.Name, cpuaware.New),
	)

	code := cli.Run(command)
	os.Exit(code)
}
