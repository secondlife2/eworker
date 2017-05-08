package work

import "eworker/bl"

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	JobQueue   chan Job
	MaxWorkers int
	MaxQueue   int
}

func NewDispatcher(maxWorkers, maxQueue int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	queue := make(chan Job, maxQueue)
	return &Dispatcher{
		WorkerPool: pool,
		JobQueue:   queue,
		MaxWorkers: maxWorkers,
		MaxQueue:   maxQueue,
	}
}
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i+1, d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			// a job request has been received
			go func(jobChannel Job) {
				// try to obtain a worker that is available.
				// this will block until a worker is idle
				worker := <-d.WorkerPool

				// dispatch the job to the worker, dequeuing from
				// the jobChannel
				worker <- jobChannel
			}(job)
		default:
			//log.Println("default job")
		}
	}
}

func (d *Dispatcher) Add(p bl.Payload) {
	worker := Job{payload: p}
	d.JobQueue <- worker
}
