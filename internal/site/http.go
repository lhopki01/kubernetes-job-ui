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
		collection.Mux.Lock()
		c.HTML(http.StatusOK, "cronjobs.html.tmpl", gin.H{
			"collection": collection,
		})
		collection.Mux.Unlock()
	})
	r.GET("/cronjob", func(c *gin.Context) {
		cronJob := c.Query("cronjob")
		collection.Mux.Lock()
		c.HTML(http.StatusOK, "cronjob.html.tmpl", gin.H{
			"collection": collection,
			"cronJob":    cronJob,
		})
		collection.Mux.Unlock()
	})
	r.GET("/job", func(c *gin.Context) {
		cronJob := c.Query("cronjob")
		job := c.Query("job")
		collection.Mux.Lock()
		c.HTML(http.StatusOK, "job.html.tmpl", gin.H{
			"collection": collection,
			"cronJob":    cronJob,
			"job":        job,
			"pods":       k8s.GetPodLogs(collection.Client, job),
		})
		collection.Mux.Unlock()
	})
	r.GET("/createjob", func(c *gin.Context) {
		cronJob := c.Query("cronjob")
		collection.Mux.Lock()
		c.HTML(http.StatusOK, "createjob.html.tmpl", gin.H{
			"jobOptions": collection.CronJobs[cronJob].Config,
			"cronJob":    cronJob,
		})
		collection.Mux.Unlock()
	})
	r.POST("/createjob", func(c *gin.Context) {
		cronJob := c.Query("cronjob")

		collection.Mux.Lock()
		envVars := map[string]string{}
		for _, option := range collection.CronJobs[cronJob].Config.Options {
			envVars[option.EnvVar] = c.PostForm(option.EnvVar)
		}
		spew.Dump(envVars)

		jobName, err := k8s.RunJob(collection, collection.CronJobs[cronJob].Object, envVars)
		if err != nil {
			panic(err)
		}
		collection.Mux.Unlock()
		k8s.UpdateCollection(collection)
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/job?cronjob=%s&job=%s", cronJob, jobName))
	})
	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println("failed to start server")
		os.Exit(1)
	}
}
