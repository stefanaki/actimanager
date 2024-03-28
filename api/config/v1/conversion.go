package v1

import (
	"k8s.io/apimachinery/pkg/conversion"

	"cslab.ece.ntua.gr/actimanager/api/config"
)

func Convert_v1_WorkloadAwareArgs_To_config_WorkloadAwareArgs(in *WorkloadAwareArgs, out *config.WorkloadAwareArgs, s conversion.Scope) error {
	if err := autoConvert_v1_WorkloadAwareArgs_To_config_WorkloadAwareArgs(in, out, s); err != nil {
		return err
	}

	if in.Policy != nil {
		out.Policy = config.WorkloadAwarePolicy(*in.Policy)
	}
	if in.Features != nil {
		out.Features = make([]config.Feature, len(*in.Features))
		for i := range *in.Features {
			out.Features[i] = config.Feature((*in.Features)[i])
		}
	}
	return nil
}

func Convert_config_WorkloadAwareArgs_To_v1_WorkloadAwareArgs(in *config.WorkloadAwareArgs, out *WorkloadAwareArgs, scope conversion.Scope) error {
	if err := autoConvert_config_WorkloadAwareArgs_To_v1_WorkloadAwareArgs(in, out, scope); err != nil {
		return err
	}
	policy := WorkloadAwarePolicy(in.Policy)
	out.Policy = &policy
	out.Features = &[]Feature{}
	for i := range in.Features {
		*out.Features = append(*out.Features, Feature(in.Features[i]))
	}
	return nil
}
