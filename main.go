package main

import (
	"github.com/lhopki01/kubernetes-job-ui/internal/k8s"
	"github.com/lhopki01/kubernetes-job-ui/internal/site"
)

func main() {
	collection := k8s.InitializeStruct()
	site.Serve(collection)
}
