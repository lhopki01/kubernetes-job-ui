package site

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lhopki01/kubernetes-job-ui/internal/k8s"
)

func Serve(collection *k8s.Collection) {
	r := gin.Default()
	r.Use(cors.Default())

	//r.Static("/static", "./static")

	//r.LoadHTMLGlob("templates/*.tmpl")

	r.GET("api/v1/cronjobs", func(c *gin.Context) {
		c.JSON(http.StatusOK, collection.GetCronJobs())
	})
	r.GET("api/v1/namespaces/:namespace/jobs/:jobname", func(c *gin.Context) {
		jobName := c.Param("jobname")
		c.JSON(http.StatusOK, collection.GetPodLogs(jobName))
	})
	r.POST("api/v1/namespaces/:namespace/cronjobs/:cronJobName", func(c *gin.Context) {
		namespace := c.Param("namespace")
		cronJobName := c.Param("cronJobName")
		var options []k8s.ResponseOption
		c.BindJSON(&options)
		spew.Dump(options)
		validationErrors := collection.ValidateEnvVars(namespace, cronJobName, options)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusUnprocessableEntity, validationErrors)
		} else {
			job, err := collection.RunJob(namespace, cronJobName, options)
			if err != nil {
				spew.Dump(err)
			} else {
				c.JSON(http.StatusOK, job)
			}
		}
	})

	//r.GET("api/v1/cronjobs/:cronJobName", func(c *gin.Context) {
	//	cronJobName := c.Param("cronJobName")
	//	c.JSON(200, collection.GetCronJob(cronJobName))
	//})
	//r.GET("api/v1/cronjobs/:cronJobName/jobs", func(c *gin.Context) {
	//	cronJobName := c.Param("cronJobName")
	//	c.JSON(200, collection.GetCronJob(cronJobName).Jobs)
	//})
	//r.GET("api/v1/cronjobs/:cronJobName/jobs/:jobName", func(c *gin.Context) {
	//	cronJobName := c.Param("cronJobName")
	//	jobName := c.Param("jobName")
	//	c.JSON(200, collection.GetJob(cronJobName, jobName))
	//})

	//r.GET("/cronjobs", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "cronjobs.html.tmpl", gin.H{
	//		"cronJobs": collection.GetCronJobs(),
	//	})
	//})
	//r.GET("/cronjob", func(c *gin.Context) {
	//	cronJobName := c.Query("cronjob")
	//	c.HTML(http.StatusOK, "cronjob.html.tmpl", gin.H{
	//		"cronJob":  collection.GetCronJob(cronJobName),
	//		"cronJobs": collection.GetCronJobs(),
	//	})

	//})
	//r.GET("/job", func(c *gin.Context) {
	//	cronJobName := c.Query("cronjob")
	//	jobName := c.Query("job")
	//	job := collection.GetJob(cronJobName, jobName)
	//	c.HTML(http.StatusOK, "job.html.tmpl", gin.H{
	//		"cronJob": collection.GetCronJob(cronJobName),
	//		"job":     job,
	//		"pods":    collection.GetPodLogs(jobName),
	//	})
	//})
	//r.GET("/createjob", func(c *gin.Context) {
	//	cronJobName := c.Query("cronjob")
	//	cronJob := collection.GetCronJob(cronJobName)

	//	c.HTML(http.StatusOK, "createjob.html.tmpl", gin.H{
	//		"jobOptions": cronJob.Config,
	//		"cronJob":    cronJob,
	//	})
	//})
	//r.POST("/createjob", func(c *gin.Context) {
	//	cronJobName := c.Query("cronjob")
	//	cronJob := collection.GetCronJob(cronJobName)

	//	envVars := make([]string, len(cronJob.Config.Options))
	//	i := 0
	//	for i < len(cronJob.Config.Options) {
	//		envVars[i] = c.PostForm(strconv.Itoa(i))
	//		i++
	//	}
	//	spew.Dump(envVars)
	//	jobName, err := collection.RunJob(cronJob.Name, envVars)
	//	if err != nil {
	//		panic(err)
	//	}
	//	collection.UpdateCollection()
	//	c.Redirect(
	//		http.StatusSeeOther,
	//		fmt.Sprintf("/job?cronjob=%s&job=%s", cronJobName, jobName),
	//	)
	//})
	target := "localhost:3000"
	r.NoRoute(func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println("failed to start server")
		os.Exit(1)
	}
}
