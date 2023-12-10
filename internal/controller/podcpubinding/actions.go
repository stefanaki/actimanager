package podcpubinding

import (
	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	"reflect"
)

func needsUpdate(cpuBinding *v1alpha1.PodCpuBinding) bool {
	return cpuBinding.ObjectMeta.Annotations[v1alpha1.ActionUpdateAnnotationKey] == "true" ||
		!reflect.DeepEqual(cpuBinding.Spec, cpuBinding.Status.LastSpec)
}

func needsDelete(cpuBinding *v1alpha1.PodCpuBinding) bool {
	return cpuBinding.ObjectMeta.Annotations[v1alpha1.ActionDeleteAnnotationKey] == "true"
}
