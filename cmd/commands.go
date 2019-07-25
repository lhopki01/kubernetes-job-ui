package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lhopki01/kubernetes-job-ui/internal/k8s"
	"github.com/lhopki01/kubernetes-job-ui/internal/site"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddCommands() {
	rootCmd := &cobra.Command{
		Use:   "docker-chain-builder",
		Short: "A tool to build docker images and all their dependencies",
	}

	serveCmd := &cobra.Command{
		Use:   "serve [options]",
		Short: "Start webserver",
		Run: func(cmd *cobra.Command, args []string) {
			runServeCommand()

		},
	}

	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("namespace", "", "namespace to find CronJobs in")
	serveCmd.Flags().Bool("configured-only", false, "only show CronJobs with Configuration")
	err := viper.BindPFlags(serveCmd.Flags())
	if err != nil {
		log.Fatalf("Binding flags failed: %s", err)
	}

	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServeCommand() {
	collection := k8s.NewCollection()

	go func(c *k8s.Collection) {
		for {
			time.Sleep(time.Duration(5) * time.Second)
			c.UpdateCollection()
		}
	}(collection)

	site.Serve(collection)
}
