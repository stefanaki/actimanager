package validation

import (
	"cslab.ece.ntua.gr/actimanager/api/config"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var validWorkloadAwarePolicies = sets.NewString(
	string(config.PolicyMaximumUtilization),
	string(config.PolicyBalanced),
)

var validWorkloadAwareFeatures = sets.NewString(
	string(config.FeatureMemoryBoundExclusiveSockets),
	string(config.FeaturePhysicalCores),
)

func ValidateWorkloadAwareArgs(path *field.Path, args *config.WorkloadAwareArgs) error {
	var errors field.ErrorList
	if !validWorkloadAwarePolicies.Has(string(args.Policy)) {
		errors = append(errors, field.Invalid(path.Child("policy"), args.Policy, "invalid policy"))
	}
	for i, feature := range args.Features {
		if !validWorkloadAwareFeatures.Has(string(feature)) {
			errors = append(errors, field.Invalid(path.Child("features").Index(i), feature, "invalid feature"))
		}
	}
	return errors.ToAggregate()
}
