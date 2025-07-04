package worker

import (
	"log"
	"sync"

	"github.com/krzysu/web-crawler/internal/crawler"
	"github.com/krzysu/web-crawler/internal/database"
	"github.com/krzysu/web-crawler/internal/models"
	
	"time"
)

// Job represents a crawling job.
type Job struct {
	URL string
	ID  uint // ID of the crawl result in the database, if it exists
}

// Worker represents the worker that executes the jobs.
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
	db         *database.DB
	wg         *sync.WaitGroup // Add WaitGroup to Worker
}

// NewWorker creates a new Worker.
func NewWorker(workerPool chan chan Job, db *database.DB, wg *sync.WaitGroup) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		db:         db,
		wg:         wg,
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

				if job.ID != 0 {
					// If ID is provided, it's a re-crawl, so fetch existing result
					result, err = w.db.GetCrawlResult(job.ID)
					if err != nil {
						log.Printf("Error getting existing crawl result for ID %d: %v\n", job.ID, err)
						w.wg.Done()
						continue
					}
					result.Headings = make(map[string]int)
					result.BrokenLinks = make([]map[string]interface{}, 0)
					result.Status = "running"
					if err := w.db.UpdateCrawlResult(result); err != nil {
						log.Printf("Error updating crawl result for ID %d: %v\n", result.ID, err)
					}
				} else {
					// New crawl, create a placeholder result
					result = &models.CrawlResult{
						URL:       job.URL,
						Status:    "queued",
						Headings:  make(map[string]int),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
					err = w.db.CreateCrawlResult(result)
					if err != nil {
						log.Printf("Error creating new crawl result for URL %s: %v\n", job.URL, err)
						w.wg.Done()
						continue
					}
				}

				crawlErr := crawler.Crawl(result)
				if crawlErr != nil {
					log.Printf("Error crawling URL %s: %v\n", job.URL, crawlErr)
					result.Status = "error"
					result.ErrorMessage = crawlErr.Error()
				}

				// Update the database with the crawled result
				result.UpdatedAt = time.Now()
				if err := w.db.UpdateCrawlResult(result); err != nil {
					log.Printf("Error updating crawl result for URL %s: %v\n", result.URL, err)
				}

				w.wg.Done()

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
	JobQueue   chan Job // Add JobQueue to Dispatcher
	db         *database.DB
	wg         *sync.WaitGroup // Add WaitGroup to Dispatcher
}

// NewDispatcher creates a new Dispatcher.
func NewDispatcher(maxWorkers int, db *database.DB, wg *sync.WaitGroup) *Dispatcher {
	return &Dispatcher{
		maxWorkers: maxWorkers,
		WorkerPool: make(chan chan Job, maxWorkers),
		JobQueue:   make(chan Job, 100), // Initialize JobQueue here
		db:         db,
		wg:         wg, // Use the provided WaitGroup
	}
}

// Run starts the workers and listens for jobs.
func (d *Dispatcher) Run() {
	// Start the workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.db, d.wg)
		worker.Start()
	}

	go d.dispatch()
}

// dispatch listens for jobs on the JobQueue and dispatches them to available workers.
func (d *Dispatcher) dispatch() {
	for job := range d.JobQueue { // Use d.JobQueue
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
