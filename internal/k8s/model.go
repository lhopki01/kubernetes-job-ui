package k8s

import (
	"sync"

	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Pod struct {
	Name         string
	Namespace    string
	CreationTime metav1.Time
	Phase        corev1.PodPhase
	Logs         string
}

type Job struct {
	Name         string
	Namespace    string
	CreationTime metav1.Time
	Running      bool
	Passed       bool
	Manual       bool
	Status       string
	Pods         []Pod
}

type CronJob struct {
	Name      string
	Namespace string
	Schedule  string
	Jobs      map[string]Job
	Config    JobOptions
	Object    *v1beta1.CronJob
}

type Collection struct {
	CronJobs map[string]CronJob
	Jobs     map[string]Job
	Client   *kubernetes.Clientset
	Mux      sync.Mutex
}

type JobOptions struct {
	Options []Option `json:"options"`
	Error   error
}

type Option struct {
	EnvVar         string   `json:"envvar"`
	Values         []string `json:"values"`
	Default        string   `json:"default"`
	Description    string   `json:"Description"`
	ContainerIndex int      `json:"container_index"`
}
