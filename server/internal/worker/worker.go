package worker

import (
	"encoding/json"
	"log"

	"github.com/krzysu/web-crawler/internal/crawler"
	"github.com/krzysu/web-crawler/internal/database"
	"github.com/krzysu/web-crawler/internal/models"
)

// Job represents a crawling job.
type Job struct {
	URL string
	ID  string // ID of the crawl result in the database, if it exists
}

// JobQueue is a buffered channel that holds incoming jobs.
var JobQueue chan Job

// Worker represents the worker that executes the jobs.
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

// NewWorker creates a new Worker.
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

// Start starts the worker by listening for jobs on its JobChannel.
func (w Worker) Start() {
	go func() {
		for {
			// Add my JobChannel to the worker pool.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// We have received a work request.
				log.Printf("Processing job for URL: %s\n", job.URL)

				var result *models.CrawlResult
				var err error

				if job.ID != "" {
					// If ID is provided, it's a re-crawl, so fetch existing result
					result, err = database.GetCrawlResult(job.ID)
					if err != nil {
						log.Printf("Error getting existing crawl result for ID %s: %v\n", job.ID, err)
						continue
					}
					result.Status = "running"
					database.UpdateCrawlResult(result)
				} else {
					// New crawl, create a placeholder result
					result = &models.CrawlResult{URL: job.URL, Status: "queued"}
					database.CreateCrawlResult(result)
				}

				crawledResult, crawlErr := crawler.Crawl(job.URL)
				if crawlErr != nil {
					log.Printf("Error crawling URL %s: %v\n", job.URL, crawlErr)
					crawledResult.Status = "error"
					crawledResult.ErrorMessage = crawlErr.Error()
				}

				// Update the database with the crawled result
				if job.ID != "" {
					crawledResult.ID = job.ID // Ensure we update the correct record
					database.UpdateCrawlResult(crawledResult)
				} else {
					database.UpdateCrawlResult(crawledResult)
				}

			case <-w.quit:
				// We have received a signal to stop
				log.Println("Worker stopping")
				return
			}
		}
	}()
}

// Stop tells the worker to stop.
func (w Worker) Stop() {
	w.quit <- true
}

// Dispatcher manages the worker pool.
type Dispatcher struct {
	maxWorkers int
	WorkerPool chan chan Job
}

// NewDispatcher creates a new Dispatcher.
func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		maxWorkers: maxWorkers,
		WorkerPool: make(chan chan Job, maxWorkers),
	}
}

// Run starts the workers and listens for jobs.
func (d *Dispatcher) Run() {
	// Start the workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

// dispatch listens for jobs on the JobQueue and dispatches them to available workers.
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// A job request has been received
			go func(job Job) {
				// Try to obtain a worker job channel that is available.
				// This will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// Dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

func init() {
	JobQueue = make(chan Job, 100) // Buffer up to 100 jobs
}
