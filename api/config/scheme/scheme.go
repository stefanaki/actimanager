package scheme

import (
	"cslab.ece.ntua.gr/actimanager/api/config"
	v1 "cslab.ece.ntua.gr/actimanager/api/config/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeschedulerscheme "k8s.io/kubernetes/pkg/scheduler/apis/config/scheme"
)

var (
	// Re-use the in-tree Scheme.
	Scheme = kubeschedulerscheme.Scheme

	// Codecs provides access to encoding and decoding for the scheme.
	Codecs = serializer.NewCodecFactory(Scheme, serializer.EnableStrict)
)

func init() {
	AddToScheme(Scheme)
}

// AddToScheme builds the kubescheduler scheme using all known versions of the kubescheduler api.
func AddToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(config.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
}
