package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/lhopki01/kubernetes-job-ui/internal/k8s"
	"k8s.io/client-go/kubernetes"
)

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	clientset := k8s.NewClient()

	ticker := time.NewTicker(250 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Done!")
				return
			case <-ticker.C:
				g.Update(func(g *gocui.Gui) error {
					cronJobs(g, clientset)
					return nil
				})
				g.Update(func(g *gocui.Gui) error {
					jobs(g, clientset)
					return nil
				})
			}
		}
	}()
	if err = g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("cronJobs", -1, -1, maxX/4, maxY); err != nil {
		v.Highlight = true
	}
	if v, err := g.SetView("jobs", maxX/4, -1, maxX/2, maxY); err != nil {
		v.Highlight = true
	}
	return nil
}

func cronJobs(g *gocui.Gui, clientset *kubernetes.Clientset) {
	v, _ := g.View("cronJobs")
	v.Clear()
	names, _ := k8s.GetCronJobs(clientset)
	for _, name := range names {
		fmt.Fprintf(v, "%s\n", name)
	}
}

func jobs(g *gocui.Gui, clientset *kubernetes.Clientset) {
	v, _ := g.View("jobs")
	v.Clear()
	names, _ := k8s.GetJobs(clientset)
	for _, name := range names {
		fmt.Fprintf(v, "%s\n", name)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
