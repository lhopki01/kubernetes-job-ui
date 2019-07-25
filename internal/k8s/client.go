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

	"github.com/spf13/viper"
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
		kubeconfig = flag.String(
			"kubeconfig",
			filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file",
		)
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
	collection.Client = c
	collection.UpdateCollection()
	return collection
}

func (c *Collection) UpdateCollection() {
	cjs := getCronJobs(c.Client)
	_, js := getJobs(c.Client)

	var cronJobs []CronJob
	for _, cj := range cjs.Items {
		cronJob := CronJob{
			Name:         cj.Name,
			Namespace:    cj.Namespace,
			CreationTime: cj.CreationTimestamp,
			Schedule:     cj.Spec.Schedule,
			Object:       cj.DeepCopy(),
			Jobs:         orderedOwnedJobs(js.Items, cj.Name),
		}
		if *cj.Spec.Suspend {
			cronJob.Schedule = "Disabled"
		}

		if val, ok := cj.Annotations["kubernetes-job-runner.io/config"]; ok {
			err := json.Unmarshal([]byte(val), &cronJob.Config)
			if err != nil {
				cronJob.Config.Error = err
			}
		} else if viper.GetBool("configured-only") {
			continue
		} else {
			for i, container := range cj.Spec.JobTemplate.Spec.Template.Spec.Containers {
				for _, v := range container.Env {
					cronJob.Config.Options = append(cronJob.Config.Options, Option{
						EnvVar:         v.Name,
						Default:        v.Value,
						ContainerIndex: i,
					})
				}
			}
		}
		cronJobs = insertCronJobIntoSliceByCreationTime(cronJobs, cronJob)
	}
	c.Lock()
	c.CronJobs = cronJobs
	c.Unlock()
}

func orderedOwnedJobs(js []batchv1.Job, cronJobName string) []Job {
	var jobs []Job
	for _, j := range js {
		for _, owner := range j.GetOwnerReferences() {
			if owner.Name == cronJobName {
				job := Job{
					Name:         j.Name,
					Namespace:    j.Namespace,
					CreationTime: j.CreationTimestamp,
				}
				if j.Annotations["cronjob.kubernetes.io/instantiate"] == "manual" {
					job.Manual = true
				}
				if j.Status.Succeeded > 0 {
					job.Passed = true
					job.Status = "succeeded"
				} else if j.Status.Active > 0 {
					job.Status = "active"
				} else if j.Status.Failed == *j.Spec.BackoffLimit+int32(1) {
					job.Status = "failed"
				}
				jobs = insertJobIntoSliceByCreationTime(jobs, job)
			}
		}
	}
	return jobs
}

func (c *Collection) GetCronJob(cronJobName string) CronJob {
	c.Lock()
	defer c.Unlock()
	for _, cj := range c.CronJobs {
		if cj.Name == cronJobName {
			return cj
		}
	}
	return CronJob{}
}

func (c *Collection) GetJob(cronJobName string, jobName string) Job {
	cronJob := c.GetCronJob(cronJobName)
	c.Lock()
	defer c.Unlock()
	for _, j := range cronJob.Jobs {
		if j.Name == jobName {
			return j
		}
	}
	return Job{}
}

func insertJobIntoSliceByCreationTime(js []Job, job Job) []Job {
	var jobs []Job
	for i, j := range js {
		if j.CreationTime.Before(&job.CreationTime) {
			return append(append(jobs, job), js[i:]...)
		} else {
			jobs = append(jobs, j)
		}
	}
	jobs = append(jobs, job)
	return jobs
}

func insertCronJobIntoSliceByCreationTime(cjs []CronJob, cronJob CronJob) []CronJob {
	var cronJobs []CronJob
	for i, cj := range cjs {
		if cj.CreationTime.Before(&cronJob.CreationTime) {
			cronJobs = append(append(cronJobs, cronJob), cjs[i:]...)
			return cronJobs
		} else {
			cronJobs = append(cronJobs, cj)
		}
	}
	cronJobs = append(cronJobs, cronJob)
	return cronJobs
}

func (c *Collection) GetCronJobs() []CronJob {
	c.Lock()
	defer c.Unlock()
	cronJobs := c.CronJobs
	return cronJobs
}

func getCronJobs(clientset *kubernetes.Clientset) (cronJobs *v1beta1.CronJobList) {
	namespace := viper.GetString("namespace")
	cronJobs, err := clientset.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		println(err)
	}

	return cronJobs
}

func getJobs(clientset *kubernetes.Clientset) (names []string, jobs *batchv1.JobList) {
	namespace := viper.GetString("namespace")
	jobs, err := clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		println(err)
	}

	for _, job := range jobs.Items {
		names = append(names, job.Name)
	}

	return names, jobs
}

func GetPod(clientset *kubernetes.Clientset, job string) (pods []Pod) {

	namespace := viper.GetString("namespace")
	ps, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{
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
	namespace := viper.GetString("namespace")
	ps, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{
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
				fmt.Printf("failed to stream logs with err: %v\n", err)
			}
			defer podLogs.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, podLogs)
			if err != nil {
				fmt.Printf("failed to copy logs with err: %v\n", err)
			}
			//str := strings.ReplaceAll(buf.String(), "\n", "<br>")
			str := buf.String()
			pod.Logs = str
		}
		pods = append(pods, pod)
	}
	return pods
}

func (c *Collection) RunJob(cronJob *v1beta1.CronJob, envVars map[string]string) (name string, err error) {
	newJobObject := createJobFromCronJob(cronJob, envVars)
	job, err := c.Client.BatchV1().Jobs(cronJob.Namespace).Create(
		newJobObject,
	)
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

	name := fmt.Sprintf("%s-%v-m", cronJob.Name, time.Now().Unix())

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
