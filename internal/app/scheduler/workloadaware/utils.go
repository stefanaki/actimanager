package workloadaware

import (
	"cslab.ece.ntua.gr/actimanager/api/config"
	"cslab.ece.ntua.gr/actimanager/api/config/validation"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func parseArgs(obj runtime.Object) (*config.WorkloadAwareArgs, error) {
	args, ok := obj.(*config.WorkloadAwareArgs)
	if !ok {
		return nil, fmt.Errorf("want args to be of type PodGroupsArgs, got %T", obj)
	}
	err := validation.ValidateWorkloadAwareArgs(nil, args)
	if err != nil {
		return nil, err
	}
	return args, nil
}

func (w *WorkloadAware) getState(state *framework.CycleState) (*State, error) {
	data, err := state.Read(framework.StateKey(Name))
	if err != nil {
		return nil, fmt.Errorf("could not read state data: %v", err)
	}
	stateData, ok := data.(*State)
	if !ok {
		return nil, fmt.Errorf("could not cast state data")
	}
	return stateData, nil
}
