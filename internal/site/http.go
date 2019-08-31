package site

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/lhopki01/kubernetes-job-ui/internal/k8s"
)

func Serve(collection *k8s.Collection) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(static.Serve("/", BinaryFileSystem("")))
	r.NoRoute(func(c *gin.Context) {
		index, err := Asset("index.html")
		if err != nil {
			fmt.Println("index.html not found")
		}
		c.Writer.WriteHeader(http.StatusOK)
		_, err = c.Writer.Write(index)
		if err != nil {
			fmt.Println("could not write index.html to response body")
		}
	})

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

	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println("failed to start server")
		os.Exit(1)
	}
}

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: root}
	return &binaryFileSystem{
		fs,
	}
}
