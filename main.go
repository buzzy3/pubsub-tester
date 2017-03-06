package main

import (
	"flag"
	"github.com/cucumber-tony/pubsub/subscriber"
	"log"
	"os"
)

var (
	Environment = flag.String("environment", "test", "Set the development mode")
	// Creds        = flag.String("creds", "./google/pubsub-key.json", "A path to your JSON key file for BigQuery && PubSub")
	ProjectID    = flag.String("project-id", "emulator-project-id", "Set the GCE projcet")
	Subscription = flag.String("sub", "my-topic-sub", "Set the GCE sub")
	// TriggerTopic = flag.String("trigger-topic", "cucumber-triggers-v2", "Set the GCE triggers topic")

	NWorkers = flag.Int("workers", 1, "The number of workers, working")

	debug = flag.Bool("debug", false, "Enable debuggin")
)

func init() {
	flag.Parse()
	os.Setenv("ENV", *Environment)
}

func main() {

	agent, err := subscriber.NewAgent(*ProjectID)
	if err != nil {
		panic(err)
	}
	agent.Subscription = *Subscription

	it := agent.Subscribe()
	defer it.Stop()

	StartDispatcher(*NWorkers)

	for {
		msg, err := it.Next()
		if err != nil {
			log.Fatalf("could not pull: %v", err)
		}

		work := WorkRequest{msg.Data}
		WorkQueue <- work

		msg.Done(true)
	}
}
