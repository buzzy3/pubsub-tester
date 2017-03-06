package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Worker is what holds the inbound jobs
type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// NewWorker is the struct that holds the jobs
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// Start will start the dispathcher
func (w Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				log.Printf("worker %d: Received work request\n", w.ID)
				Process(work)
			case <-w.QuitChan:
				log.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop will cease the workers
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// Process handles the incoming request
func Process(work WorkRequest) bool {
	var dat map[string]interface{}
	if err := json.Unmarshal(work.Body, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	return true
}
