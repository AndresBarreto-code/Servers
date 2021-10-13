package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Job struct {
	Name   string
	Server string
}

type Worker struct {
	Id         int
	JobQueue   chan Job
	WorkerPool chan chan Job
	QuitChan   chan bool
}

type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	JobQueue   chan Job
}

func NewWorker(id int, workerPool chan chan Job) *Worker {
	return &Worker{
		Id:         id,
		JobQueue:   make(chan Job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				fmt.Printf("Worker whit id %d started.\n", w.Id)
				message := checkServer(job.Server)
				fmt.Printf("Worker whit id %d finished. Result: %s. Name: %s.\n", w.Id, message, job.Name)
			case <-w.QuitChan:
				fmt.Printf("Worker whit id %d stopped.\n.", w.Id)

			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func NewDispacher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	return &Dispatcher{
		JobQueue:   jobQueue,
		MaxWorkers: maxWorkers,
		WorkerPool: make(chan chan Job, maxWorkers),
	}
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			go func() {
				workerJobQueue := <-d.WorkerPool
				workerJobQueue <- job
			}()
		}
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i, d.WorkerPool)
		worker.Start()
	}
	go d.Dispatch()
}

func main() {
	const (
		maxWorkers   = 4
		maxQueueSize = 20
		port         = ":8081"
	)
	jobQueue := make(chan Job, maxQueueSize)
	dispatcher := NewDispacher(jobQueue, maxWorkers)
	dispatcher.Run()

	start := time.Now()
	servers := []string{
		"http://platzi.com",
		"http://google.com",
		"http://instagram.com",
		"http://facebook.com",
		"http://twitter.com",
	}

	for _, server := range servers {
		job := Job{
			Name:   server,
			Server: server,
		}
		jobQueue <- job
	}

	timeTaken := time.Since(start)
	fmt.Printf("exec time %s\n", timeTaken)
	log.Fatal(http.ListenAndServe(port, nil))

}

type Writer struct{}

func (Writer) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}

func checkServer(server string) string {
	response, err := http.Get(server)
	if err != nil {
		w := Writer{}
		io.Copy(w, response.Body)
		return server + " is not working"
	} else {
		return server + " is working correctly"
	}
}
