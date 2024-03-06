package podcpubinding

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
)

// PodCpuBindingFinalizer removes the CPU binding that was applied on the Pod
func (r *PodCpuBindingReconciler) PodCpuBindingFinalizer(ctx context.Context, cpuBinding *v1alpha1.PodCpuBinding, logger logr.Logger) error {
	pod, err := r.getPod(ctx, types.NamespacedName{
		Namespace: cpuBinding.Namespace,
		Name:      cpuBinding.Spec.PodName,
	})
	if err != nil {
		logger.Info("failed to get pod on finalize", "error", err.Error())
		return nil
	}
	// Remove CPU pinning and delete CR
	err = r.removeCpuPinning(ctx, pod)
	if err != nil {
		logger.Info("error removing CPU pinning", "error", err.Error())
	}
	return nil
}
