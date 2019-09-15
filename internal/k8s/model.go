package k8s

import (
	"sync"

	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Container struct {
	Name string `json:"name"`
	Logs string `json:"logs"`
}

type Pod struct {
	Name         string      `json:"name"`
	Namespace    string      `json:"namespace"`
	CreationTime metav1.Time `json:"creationTime"`
	Status       string      `json:"status"`
	Containers   []Container `json:"containers"`
}

type Job struct {
	Name         string      `json:"name"`
	Namespace    string      `json:"namespace"`
	CreationTime metav1.Time `json:"creationTime"`
	Manual       bool        `json:"manual"`
	Status       string      `json:"status"`
	Pods         []Pod       `json:"pods"`
}

type CronJob struct {
	Name         string      `json:"name"`
	Namespace    string      `json:"namespace"`
	CreationTime metav1.Time `json:"creationTime"`
	Schedule     string      `json:"schedule"`
	Jobs         []Job       `json:"jobs"`
	Config       Config      `json:"config"`
	object       *v1beta1.CronJob
}

type Collection struct {
	sync.Mutex
	cronJobs      []CronJob
	monitoredJobs map[string]Job
	Client        *kubernetes.Clientset
}

type Config struct {
	Description string   `json:"description" yaml:"description"`
	Options     []Option `json:"options" yaml:"options"`
	Error       string   `json:"error" yaml:"error"`
	Errors      []string `json:"errors" yaml:"errors"`
	Raw         string   `json:"raw" yaml:"raw"`
}

const (
	List   = "list"
	String = "string"
)

type Option struct {
	EnvVar         string   `json:"envVar" yaml:"envVar"`
	Type           string   `json:"type" yaml:"type"`
	Values         []string `json:"values" yaml:"values"`
	Regex          string   `json:"regex" yaml:"regex"`
	Default        string   `json:"default" yaml:"default"`
	Description    string   `json:"description" yaml:"description"`
	Container      string   `json:"container" yaml:"container"`
	containerIndex int
}

type ResponseOption struct {
	EnvVar    string `json:"envVar"`
	Container string `json:"container"`
	Value     string `json:"value"`
}

type ValidationError struct {
	EnvVar      string `json:"envVar"`
	Container   string `json:"container"`
	OptionIndex int    `json:"optionIndex"`
	Error       string `json:"error"`
}

type CreateResponse struct {
	Job string `json:"job"`
}

type ByContainerIndex []Option

func (a ByContainerIndex) Len() int           { return len(a) }
func (a ByContainerIndex) Less(i, j int) bool { return a[i].containerIndex < a[j].containerIndex }
func (a ByContainerIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
