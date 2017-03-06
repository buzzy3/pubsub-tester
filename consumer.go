package main

import (
	"log"
)

var (
	queues    = make(map[string]string)
	WorkQueue = make(chan WorkRequest, 100)
	c         Consumer
)

type BindOptions struct {
	q        string
	keys     []string
	preFetch int
}

type Consumer struct {
	tag  string
	done chan error
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
