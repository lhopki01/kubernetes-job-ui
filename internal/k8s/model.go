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
	Containers   []Container
}

type Container struct {
	Name string
	Logs string
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
	Name         string
	Namespace    string
	CreationTime metav1.Time
	Schedule     string
	Jobs         []Job
	Config       JobOptions
	Object       *v1beta1.CronJob
}

type Collection struct {
	sync.Mutex
	CronJobs []CronJob
	Client   *kubernetes.Clientset
}

type JobOptions struct {
	Description string
	Options     []Option
	Error       string
	Raw         string
}

const (
	List   = "list"
	String = "string"
)

type Option struct {
	EnvVar         string
	Type           string
	Values         []string
	Default        string
	Description    string
	Container      string
	ContainerIndex int
}

type ResponseOption struct {
	EnvVar    string
	Container string
	Value     string
}

type ValidationError struct {
	EnvVar      string
	Container   string
	OptionIndex int
	Error       string
}

type CreateResponse struct {
	Job string
}

type ByContainerIndex []Option

func (a ByContainerIndex) Len() int           { return len(a) }
func (a ByContainerIndex) Less(i, j int) bool { return a[i].ContainerIndex < a[j].ContainerIndex }
func (a ByContainerIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
