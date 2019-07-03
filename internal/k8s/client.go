package k8s

import (
	"flag"
	"os"
	"path/filepath"

	v1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func NewClient() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func GetCronJobs(clientset *kubernetes.Clientset) (names []string, cronJobs *v1beta1.CronJobList) {
	cronJobs, err := clientset.BatchV1beta1().CronJobs("luke").List(metav1.ListOptions{})
	if err != nil {
		println(err)
	}

	for _, cronJob := range cronJobs.Items {
		names = append(names, cronJob.Name)
	}

	return names, cronJobs
}

func GetJobs(clientset *kubernetes.Clientset) (names []string, jobs *v1.JobList) {
	jobs, err := clientset.BatchV1().Jobs("luke").List(metav1.ListOptions{})
	if err != nil {
		println(err)
	}

	for _, job := range jobs.Items {
		names = append(names, job.Name)
	}

	return names, jobs
}

func GetPods(clientset *kubernetes.Clientset) (names []string, pods *v1.PodList) {
	jobs, err := clientset.CoreV1().Pods("luke").List(metav1.ListOptions{})
	if err != nil {
		println(err)
	}

	for _, pod := range pods.Items {
		names = append(names, pod.Name)
	}

	return names, pods
}
