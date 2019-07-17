package main

import "github.com/lhopki01/kubernetes-job-ui/cmd"

func main() {
	//collection := k8s.NewCollection()

	//for k, v := range collection.CronJobs {
	//	fmt.Printf("key: %s\nvalue: %s\n", k, v.Object.Name)
	//}

	//go func(collection *k8s.Collection) {
	//	for {
	//		time.Sleep(time.Duration(15) * time.Second)
	//		k8s.UpdateCollection(collection)
	//	}
	//}(collection)
	//site.Serve(collection)

	cmd.AddCommands()
}
