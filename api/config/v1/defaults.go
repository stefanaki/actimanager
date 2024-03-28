package v1

var (
	// DefaultPolicy is the default capacity of the WorkloadAware plugin.
	DefaultPolicy WorkloadAwarePolicy = PolicyMaximumUtilization
)

// SetDefaults_WorkloadAwareArgs sets the default parameters for WorkloadAware plugin.
func SetDefaults_WorkloadAwareArgs(obj *WorkloadAwareArgs) {
	if obj.Policy == nil || *obj.Policy == "" {
		obj.Policy = &DefaultPolicy
	}
	if obj.Features == nil {
		obj.Features = &[]Feature{}
	}
}
