package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Pod struct {
	Name         string
	Namespace    string
	CreationTime metav1.Time
	Passed       string
	Logs         string
}

type Job struct {
	Name         string
	Namespace    string
	CreationTime metav1.Time
	Running      bool
	Passed       bool
	Pods         []Pod
}

type Cronjob struct {
	Name      string
	Namespace string
	Schedule  string
	Jobs      []Job
}

type Collection struct {
	Cronjobs []Cronjob
	Jobs     map[string]Job
	Client   *kubernetes.Clientset
}

type JobOptions struct {
	Options []Option
}

type Option struct {
	EnvVar      string
	Options     []string
	Default     string
	Description string
}
