package site

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/lhopki01/kubernetes-job-ui/internal/k8s"
)

func Serve(collection *k8s.Collection) {
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()

	//r.Static("/static", "./static")

	//r.LoadHTMLGlob("templates/*.tmpl")

	r.GET("api/v1/cronjobs", func(c *gin.Context) {
		c.JSON(http.StatusOK, collection.GetCronJobs())
	})
	r.GET("api/v1/namespaces/:namespace/cronjobs/:cronjobname/jobs/:jobname", func(c *gin.Context) {
		namespace := c.Param("namespace")
		jobName := c.Param("jobname")
		cronJobName := c.Param("cronjobname")
		c.JSON(http.StatusOK, collection.GetJobLogs(namespace, cronJobName, jobName))
	})
	r.POST("api/v1/namespaces/:namespace/cronjobs/:cronJobName", func(c *gin.Context) {
		namespace := c.Param("namespace")
		cronJobName := c.Param("cronJobName")

		var options []k8s.ResponseOption
		err := c.BindJSON(&options)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "{}")
		}

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
