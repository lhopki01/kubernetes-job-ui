package site

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lhopki01/kubernetes-job-ui/internal/k8s"
)

func Serve() {
	r := gin.Default()

	r.Static("/static", "./static")

	cronjobs := []k8s.Cronjob{}
	cronjobs = append(cronjobs, k8s.Cronjob{
		Name:     "foobar",
		Schedule: "*/5 * * * *",
		Jobs: []k8s.Job{
			k8s.Job{
				Name:   "job1",
				Passed: true,
			},
			k8s.Job{
				Name:   "job2",
				Passed: true,
			},
			k8s.Job{
				Name:   "job3",
				Passed: false,
			},
			k8s.Job{
				Name:   "job4",
				Passed: true,
			},
			k8s.Job{
				Name:   "job5",
				Passed: true,
			},
		},
	})
	cronjobs = append(cronjobs, k8s.Cronjob{
		Name:     "barfoo",
		Schedule: "",
		Jobs: []k8s.Job{
			k8s.Job{
				Name:   "job1",
				Passed: false,
			},
			k8s.Job{
				Name:   "job2",
				Passed: true,
			},
		},
	})

	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/hello", func(c *gin.Context) {
		c.HTML(http.StatusOK, "helloworld.html.tmpl", gin.H{
			"cronjobs": cronjobs,
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
