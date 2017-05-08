package work

import (
	"eworker/bl"
	"log"
)

// Job represents the job to be run
type Job struct {
	payload bl.Payload
}

// A buffered channel that we can send work requests on.
// var JobQueue chan Job

// A pool of workers that are instantianted to perform the work
// var WorkerPool chan chan Job

// Worker represents the worker that executes the job
type Worker struct {
	ID         int
	JobChannel chan Job
	WorkerPool chan chan Job
	QuitChan   chan bool
}

func NewWorker(id int, workerPool chan chan Job) Worker {
	worker := Worker{
		ID:         id,
		JobChannel: make(chan Job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}

	return worker
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:

				// we have received a work request.
				log.Println(job.payload.HealthCheck())

			case <-w.QuitChan:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// func Dispatch(maxWorkers, maxQueue int) {

// 	WorkerPool = make(chan chan Job, maxWorkers)
// 	JobQueue = make(chan Job, maxQueue)

// 	// starting n number of workers
// 	for i := 0; i < maxWorkers; i++ {

// 		worker := NewWorker(i+1, WorkerPool)
// 		worker.Start()
// 	}

// 	go func() {

// 		for {
// 			select {
// 			case job := <-JobQueue:
// 				// a job request has been received
// 				go func(jobChannel Job) {
// 					// try to obtain a worker that is available.
// 					// this will block until a worker is idle
// 					worker := <-WorkerPool

// 					// dispatch the job to the worker, dequeuing from
// 					// the jobChannel
// 					worker <- jobChannel
// 				}(job)
// 			default:
// 				//log.Println("default job")
// 			}
// 		}
// 	}()
// }

// func Add(p Payload) {
// 	worker := Job{payload: p}
// 	JobQueue <- worker
// }
