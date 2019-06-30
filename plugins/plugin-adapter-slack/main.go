package main

import (
	"log"

	"github.com/andersnormal/autobot/pkg/plugins"
)

type slackPlugin struct {
}

func main() {
	// create root context
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// plugin ....
	plugin := plugins.New("slack-adapter")

	go func() {
		// subscribe ...
		for {
			select {
			case e, ok := <-plugin.SubscribeMessages():
				if !ok {
					return
				}

				log.Printf("received event: %v", e)
			}
		}
	}()

	if err := plugin.Wait(); err != nil {
		log.Fatalf("stopped plugin: %v", err)
	}
}
