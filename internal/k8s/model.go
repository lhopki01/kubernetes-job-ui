package k8s

type Pod struct {
	Name         string
	CreationTime int
	Passed       string
}

type Job struct {
	Name         string
	CreationTime int
	Running      bool
	Passed       bool
	Pods         []Pod
}

type Cronjob struct {
	Name     string
	Schedule string
	Jobs     []Job
}
