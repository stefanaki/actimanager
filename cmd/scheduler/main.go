package main

import (
	"cslab.ece.ntua.gr/actimanager/internal/scheduler/workloadaware"
	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // for rest client metric registration
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"os"

	_ "cslab.ece.ntua.gr/actimanager/api/config/scheme"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(workloadaware.Name, workloadaware.New),
		// app.WithPlugin(podisolation.Name, podisolation.New),
	)

	code := cli.Run(command)
	os.Exit(code)
}
