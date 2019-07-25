package site

import (
	"fmt"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/lhopki01/kubernetes-job-ui/internal/k8s"
)

func Serve(collection *k8s.Collection) {
	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/cronjobs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cronjobs.html.tmpl", gin.H{
			"cronJobs": collection.GetCronJobs(),
		})
	})
	r.GET("/cronjob", func(c *gin.Context) {
		cronJobName := c.Query("cronjob")
		c.HTML(http.StatusOK, "cronjob.html.tmpl", gin.H{
			"cronJob":  collection.GetCronJob(cronJobName),
			"cronJobs": collection.GetCronJobs(),
		})
	})
	r.GET("/job", func(c *gin.Context) {
		cronJobName := c.Query("cronjob")
		jobName := c.Query("job")
		job := collection.GetJob(cronJobName, jobName)
		c.HTML(http.StatusOK, "job.html.tmpl", gin.H{
			"cronJob": collection.GetCronJob(cronJobName),
			"job":     job,
			"pods":    k8s.GetPodLogs(collection.Client, jobName),
		})
	})
	r.GET("/createjob", func(c *gin.Context) {
		cronJobName := c.Query("cronjob")
		cronJob := collection.GetCronJob(cronJobName)

		c.HTML(http.StatusOK, "createjob.html.tmpl", gin.H{
			"jobOptions": cronJob.Config,
			"cronJob":    cronJob.Name,
		})
	})
	r.POST("/createjob", func(c *gin.Context) {
		cronJobName := c.Query("cronjob")
		cronJob := collection.GetCronJob(cronJobName)

		envVars := map[string]string{}
		spew.Dump(c.PostForm())
		for _, option := range cronJob.Config.Options {
			envVars[option.EnvVar] = c.PostForm(option.EnvVar)
		}
		spew.Dump(envVars)
		jobName, err := collection.RunJob(cronJob.Object, envVars)
		if err != nil {
			panic(err)
		}
		collection.UpdateCollection()
		c.Redirect(
			http.StatusSeeOther,
			fmt.Sprintf("/job?cronjob=%s&job=%s", cronJobName, jobName),
		)
	})

	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println("failed to start server")
		os.Exit(1)
	}
}
