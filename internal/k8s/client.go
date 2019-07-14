package k8s

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
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

func NewCollection() (collection *Collection) {
	c := NewClient()
	collection = new(Collection)
	(*collection).Client = c
	UpdateCollection(collection)
	return collection
}

func UpdateCollection(collection *Collection) {
	_, cjs := GetCronJobs(collection.Client)
	_, js := GetJobs(collection.Client)

	collection.CronJobs = make(map[string]CronJob)

	for _, cj := range cjs.Items {
		cronJob := CronJob{
			Name:      cj.Name,
			Namespace: cj.Namespace,
			Schedule:  cj.Spec.Schedule,
			Object:    cj.DeepCopy(),
			Jobs:      make(map[string]Job),
		}
		if *cj.Spec.Suspend {
			cronJob.Schedule = "Disabled"
		}
		json.Unmarshal([]byte(cj.Annotations["kubernetes-job-runner.io/config"]), &cronJob.Config)
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
					cronJob.Jobs[j.Name] = job
				}
			}
		}
		collection.CronJobs[cronJob.Name] = cronJob
	}
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

func GetJobs(clientset *kubernetes.Clientset) (names []string, jobs *batchv1.JobList) {
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
		pod := Pod{
			Name:         p.Name,
			Namespace:    p.Namespace,
			CreationTime: p.CreationTimestamp,
			Phase:        p.Status.Phase,
		}
		if p.Status.Phase == "Pending" {
			pod.Logs = "Pod not running yet"

		} else {
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
			pod.Logs = str
		}
		pods = append(pods, pod)
	}
	return pods
}

func RunJob(
	c *kubernetes.Clientset,
	cronJob *v1beta1.CronJob,
	envVars map[string]string,
) (name string, err error) {
	job, err := c.BatchV1().Jobs(cronJob.Namespace).Create(createJobFromCronJob(cronJob, envVars))
	return job.Name, err
}

func createJobFromCronJob(
	cronJob *v1beta1.CronJob,
	envVars map[string]string,
) *batchv1.Job {
	annotations := make(map[string]string)
	annotations["cronjob.kubernetes.io/instantiate"] = "manual"
	for k, v := range cronJob.Spec.JobTemplate.Annotations {
		annotations[k] = v
	}

	spec := cronJob.Spec.JobTemplate.Spec
	envVarSlice := make([][]corev1.EnvVar, len(spec.Template.Spec.Containers))
	for k, v := range envVars {
		envVarSlice[0] = append(envVarSlice[0], corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	for _, v := range spec.Template.Spec.Containers[0].Env {
		if _, ok := envVars[v.Name]; ok {
			continue
		}
		envVarSlice[0] = append(envVarSlice[0], v)
	}
	spec.Template.Spec.Containers[0].Env = envVarSlice[0]

	name := fmt.Sprintf("%s-m-%v", cronJob.Name, time.Now().Unix())

	return &batchv1.Job{
		// this is ok because we know exactly how we want to be serialized
		TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "Job"},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annotations,
			Labels:      cronJob.Spec.JobTemplate.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cronJob, appsv1.SchemeGroupVersion.WithKind("CronJob")),
			},
		},
		Spec: cronJob.Spec.JobTemplate.Spec,
	}
}
