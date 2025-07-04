package client

import "github.com/skiprco/go-micro-kubernetes-registry/client/watch"

// Kubernetes ...
type Kubernetes interface {
	ListPods(labels map[string]string) (*PodList, error)
	UpdatePod(podName string, pod *Pod) (*Pod, error)
	WatchPods(labels map[string]string) (watch.Watch, error)
}

// PodList ...
type PodList struct {
	Items []Pod `json:"items"`
}

// Pod is the top level item for a pod.
type Pod struct {
	Metadata *Meta   `json:"metadata"`
	Status   *Status `json:"status"`
}

// Meta ...
type Meta struct {
	Name              string             `json:"name,omitempty"`
	Labels            map[string]*string `json:"labels,omitempty"`
	Annotations       map[string]*string `json:"annotations,omitempty"`
	DeletionTimestamp string             `json:"deletionTimestamp,omitempty"`
}

// Status ...
type Status struct {
	PodIP string `json:"podIP"`
	Phase string `json:"phase"`
}
