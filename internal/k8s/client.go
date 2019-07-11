package k8s

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	v1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
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

func InitializeStruct() (collection Collection) {
	c := NewClient()
	collection.Client = c
	_, cjs := GetCronJobs(c)
	_, js := GetJobs(c)

	collection.Jobs = make(map[string]Job)

	for _, cj := range cjs.Items {
		cronjob := Cronjob{
			Name:      cj.Name,
			Namespace: cj.Namespace,
			Schedule:  cj.Spec.Schedule,
		}
		for _, j := range js.Items {
			for _, owner := range j.GetOwnerReferences() {
				if owner.Name == cj.Name {
					job := Job{
						Name:         j.Name,
						Namespace:    j.Namespace,
						CreationTime: j.CreationTimestamp,
					}
					if j.Status.Succeeded > 0 {
						job.Passed = true
					}
					//job.Pods = GetPods(c, j.Name)
					cronjob.Jobs = append(cronjob.Jobs, job)
					collection.Jobs[j.Name] = job
				}
			}
		}
		collection.Cronjobs = append(collection.Cronjobs, cronjob)
	}
	return collection

}

func GetCronJobs(clientset *kubernetes.Clientset) (names []string, cronJobs *v1beta1.CronJobList) {
	cronJobs, err := clientset.BatchV1beta1().CronJobs("").List(metav1.ListOptions{})
	if err != nil {
		println(err)
	}

	for _, cronJob := range cronJobs.Items {
		names = append(names, cronJob.Name)
	}

	return names, cronJobs
}

func GetJobs(clientset *kubernetes.Clientset) (names []string, jobs *v1.JobList) {
	jobs, err := clientset.BatchV1().Jobs("").List(metav1.ListOptions{})
	if err != nil {
		println(err)
	}

	for _, job := range jobs.Items {
		names = append(names, job.Name)
	}

	return names, jobs
}

func GetPod(clientset *kubernetes.Clientset, job string) (pods []Pod) {
	ps, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", job),
	})
	if err != nil {
		println(err)
	}

	for _, pod := range ps.Items {
		pod := Pod{
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			CreationTime: pod.CreationTimestamp,
		}
		pods = append(pods, pod)
	}

	return pods
}

func GetPodLogs(clientset *kubernetes.Clientset, job string) (pods []Pod) {
	ps, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", job),
	})
	if err != nil {
		println(err)
	}
	podLogOpts := corev1.PodLogOptions{}
	for _, p := range ps.Items {
		req := clientset.CoreV1().Pods(p.Namespace).GetLogs(p.Name, &podLogOpts)
		podLogs, err := req.Stream()
		if err != nil {
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
		}
		//str := strings.ReplaceAll(buf.String(), "\n", "<br>")
		str := buf.String()

		pods = append(pods, Pod{
			Name:         p.Name,
			Namespace:    p.Namespace,
			CreationTime: p.CreationTimestamp,
			Logs:         str,
		})

	}
	return pods
}
