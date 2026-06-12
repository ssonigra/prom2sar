module github.com/prometheus/prometheus-dump-operator

go 1.21

require (
	github.com/prometheus/prometheus v0.48.0
	k8s.io/api v0.28.3
	k8s.io/apimachinery v0.28.3
	k8s.io/client-go v0.28.3
	sigs.k8s.io/controller-runtime v0.16.3
)

require (
	github.com/go-logr/logr v1.2.4
	github.com/prometheus/common v0.45.0
)
