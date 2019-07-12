package site

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lhopki01/kubernetes-job-ui/internal/k8s"
)

func Serve(collection k8s.Collection) {
	r := gin.Default()

	r.Static("/static", "./static")

	//Cronjobs := []k8s.Cronjob{}
	//Cronjobs = append(cronjobs, k8s.Cronjob{
	//	Name:     "foobar",
	//	Schedule: "*/5 * * * *",
	//	Jobs: []k8s.Job{
	//		k8s.Job{
	//			Name:   "job1",
	//			Passed: true,
	//		},
	//		k8s.Job{
	//			Name:   "job2",
	//			Passed: true,
	//		},
	//		k8s.Job{
	//			Name:   "job3",
	//			Passed: false,
	//		},
	//		k8s.Job{
	//			Name:   "job4",
	//			Passed: true,
	//		},
	//		k8s.Job{
	//			Name:   "job5",
	//			Passed: true,
	//		},
	//	},
	//})
	//Cronjobs = append(cronjobs, k8s.Cronjob{
	//	Name:     "barfoo",
	//	Schedule: "",
	//	Jobs: []k8s.Job{
	//		k8s.Job{
	//			Name:   "job1",
	//			Passed: false,
	//		},
	//		k8s.Job{
	//			Name:   "job2",
	//			Passed: true,
	//		},
	//	},
	//})

	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/cronjobs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "helloworld.html.tmpl", gin.H{
			"cronjobs": collection.Cronjobs,
		})
	})
	r.GET("/job", func(c *gin.Context) {
		job := c.Query("name")
		c.HTML(http.StatusOK, "job.html.tmpl", gin.H{
			"job":   collection.Jobs[job],
			"query": job,
			"pods":  k8s.GetPodLogs(collection.Client, job),
		})
	})
	jobOptions := k8s.JobOptions{
		Options: []k8s.Option{
			k8s.Option{
				EnvVar: "FOOBAR",
				Options: []string{
					"FOO",
					"BAR",
					"foobar",
				},
				Default:     "foobar",
				Description: "An option to select what type of foo",
			},
			k8s.Option{
				EnvVar: "BOOLEAN",
				Options: []string{
					"true",
					"false",
				},
				Default:     "true",
				Description: "Should we do something?",
			},
		},
	}
	r.GET("/createjob", func(c *gin.Context) {
		c.HTML(http.StatusOK, "createjob.html.tmpl", gin.H{
			"jobOptions": jobOptions,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
