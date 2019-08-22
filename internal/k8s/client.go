package k8s

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	// check if in cluster first
	_, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err == nil {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return clientset
	}

	// if not in cluster then load from config
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
			object:       cj.DeepCopy(),
			Jobs:         orderedOwnedJobs(js.Items, cj.Name, cj.Namespace),
		}
		if *cj.Spec.Suspend {
			cronJob.Schedule = "Disabled"
		}

		if val, ok := cj.Annotations["kubernetes-job-runner.io/config"]; ok {
			err := json.Unmarshal([]byte(val), &cronJob.Config)
			if jsonError, ok := err.(*json.SyntaxError); ok {
				line, character, _ := lineAndCharacter(val, int(jsonError.Offset))
				cronJob.Config.Error = fmt.Sprintf("Cannot parse JSON schema due to a syntax error at line %d, character %d: %v\n", line, character, jsonError.Error())
			} else if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
				line, character, _ := lineAndCharacter(val, int(jsonError.Offset))
				cronJob.Config.Error = fmt.Sprintf("The JSON type '%v' cannot be converted into the Go '%v' type on struct '%s', field '%v'. See input file line %d, character %d\n", jsonError.Value, jsonError.Type.Name(), jsonError.Struct, jsonError.Field, line, character)
			} else if err != nil {
				cronJob.Config.Error = err.Error()
			}
			cronJob.Config.Raw = val
			for i, container := range cj.Spec.JobTemplate.Spec.Template.Spec.Containers {
				for j, option := range cronJob.Config.Options {
					if option.Container == container.Name {
						cronJob.Config.Options[j].ContainerIndex = i
					}
				}
			}
		} else if !viper.GetBool("configured-only") {
			for i, container := range cj.Spec.JobTemplate.Spec.Template.Spec.Containers {
				for _, v := range container.Env {
					cronJob.Config.Options = append(cronJob.Config.Options, Option{
						EnvVar:         v.Name,
						Default:        v.Value,
						Type:           "string",
						Container:      container.Name,
						ContainerIndex: i,
					})
				}
			}
		}
		// Sort to make it easier to display options grouped by container on the frontend
		sort.Sort(ByContainerIndex(cronJob.Config.Options))
		cronJobs = insertCronJobIntoSliceByCreationTime(cronJobs, cronJob)
	}
	c.Lock()
	c.CronJobs = cronJobs
	c.Unlock()
}

func lineAndCharacter(input string, offset int) (line int, character int, err error) {
	lf := rune(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0, fmt.Errorf("couldn't find offset %d within the input.", offset)
	}

	// Humans tend to count from 1.
	line = 1

	for i, b := range input {
		if b == lf {
			line++
			character = 0
		}
		character++
		if i == offset {
			break
		}
	}

	return line, character, nil
}

func orderedOwnedJobs(js []batchv1.Job, cronJobName string, cronJobNamespace string) []Job {
	var jobs []Job
	for _, j := range js {
		for _, owner := range j.GetOwnerReferences() {
			if owner.Name == cronJobName && j.Namespace == cronJobNamespace {
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
				} else if j.Status.Failed >= *j.Spec.BackoffLimit+int32(1) {
					job.Status = "failed"
				}
				jobs = insertJobIntoSliceByCreationTime(jobs, job)
			}
		}
	}
	return jobs
}

func getImageTag(image string) string {
	s := strings.Split(image, ":")
	if len(s) == 2 {
		return s[1]
	}
	return "latest"
}

func (c *Collection) GetCronJob(namespace, cronJobName string) CronJob {
	c.Lock()
	defer c.Unlock()
	for _, cj := range c.CronJobs {
		if cj.Name == cronJobName && cj.Namespace == namespace {
			return cj
		}
	}
	return CronJob{}
}

func (c *Collection) GetJob(namespace, cronJobName, jobName string) Job {
	cronJob := c.GetCronJob(namespace, cronJobName)
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
		}
		jobs = append(jobs, j)
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
		}
		cronJobs = append(cronJobs, cj)
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

func (c *Collection) GetPodLogs(namespace, cronJobName, jobName string) (pods []Pod) {
	c.UpdateCollection()
	job := c.GetJob(namespace, cronJobName, jobName)
	ps, err := c.Client.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
	})
	if err != nil {
		println(err)
	}
	for _, p := range ps.Items {
		pod := Pod{
			Name:         p.Name,
			Namespace:    p.Namespace,
			CreationTime: p.CreationTimestamp,
			Phase:        p.Status.Phase,
		}
		for _, container := range p.Spec.Containers {
			if p.Status.Phase == "Pending" {
				pod.Containers = append(pod.Containers, Container{
					Name: container.Name,
					Logs: "Pod not running yet",
				})
			} else {
				podLogOpts := corev1.PodLogOptions{
					Container: container.Name,
				}
				req := c.Client.CoreV1().Pods(p.Namespace).GetLogs(p.Name, &podLogOpts)
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
				pod.Containers = append(pod.Containers, Container{
					Name: container.Name,
					Logs: str,
				})
			}
		}
		pods = append(pods, pod)
	}
	return pods
}

func (c *Collection) ValidateEnvVars(namespace, cronJobName string, envVars []ResponseOption) []ValidationError {
	cronJob := c.GetCronJob(namespace, cronJobName)
	var validationErrors []ValidationError
	for i, option := range cronJob.Config.Options {
		for _, envVar := range envVars {
			if option.Container == envVar.Container && option.EnvVar == envVar.EnvVar {
				err := ValidateEnvVar(option, envVar)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{
						EnvVar:      envVar.EnvVar,
						Container:   envVar.Container,
						OptionIndex: i,
						Error:       err.Error(),
					})
				}
			}
		}
	}
	spew.Dump(validationErrors)
	return validationErrors
}

func ValidateEnvVar(option Option, envVar ResponseOption) error {
	switch option.Type {
	case "list":
		for _, v := range option.Values {
			if v == envVar.Value {
				return nil
			}
		}
		return fmt.Errorf(
			"'%s' not one of ['%s']",
			envVar.Value,
			strings.Join(option.Values, "', '"),
		)
	case "string":
		return nil
	}
	return nil
}

func (c *Collection) RunJob(namespace, cronJobName string, envVars []ResponseOption) (CreateResponse, error) {
	newJobObject := c.createJobFromCronJob(namespace, cronJobName, envVars)
	job, err := c.Client.BatchV1().Jobs(namespace).Create(
		newJobObject,
	)
	return CreateResponse{Job: job.Name}, err
}

func (c *Collection) createJobFromCronJob(namespace, cronJobName string, envVars []ResponseOption) *batchv1.Job {
	cronJob := c.GetCronJob(namespace, cronJobName)
	annotations := make(map[string]string)
	annotations["cronjob.kubernetes.io/instantiate"] = "manual"
	for k, v := range cronJob.object.Spec.JobTemplate.Annotations {
		annotations[k] = v
	}

	spec := cronJob.object.Spec.JobTemplate.Spec
	envVarSlice := make([][]corev1.EnvVar, len(spec.Template.Spec.Containers))
	for _, v := range envVars {
		for _, w := range cronJob.Config.Options {
			if v.EnvVar == w.EnvVar && v.Container == w.Container {
				envVarSlice[w.ContainerIndex] = append(
					envVarSlice[w.ContainerIndex],
					corev1.EnvVar{
						Name:  w.EnvVar,
						Value: v.Value,
					},
				)
			}
		}
	}

	for i, c := range spec.Template.Spec.Containers {
		for _, v := range c.Env {
			if !envVarInSlice(v, envVarSlice[i]) {
				envVarSlice[i] = append(envVarSlice[i], v)
			}
		}
		spec.Template.Spec.Containers[i].Env = envVarSlice[i]
	}

	name := fmt.Sprintf("%s-%v-m", cronJob.Name, time.Now().Unix())

	return &batchv1.Job{
		// this is ok because we know exactly how we want to be serialized
		TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "Job"},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annotations,
			Labels:      cronJob.object.Spec.JobTemplate.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cronJob.object, appsv1.SchemeGroupVersion.WithKind("CronJob")),
			},
		},
		Spec: cronJob.object.Spec.JobTemplate.Spec,
	}
}

func envVarInSlice(envVar corev1.EnvVar, slice []corev1.EnvVar) bool {
	for _, v := range slice {
		if v.Name == envVar.Name {
			return true
		}
	}
	return false
}
