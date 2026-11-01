package useSchedule

import (
	"log"
	"sync"
)

// NewQueue create a new job Queue
// and assign all required parameters
func NewQueue(workers int, capacity int, jobCallBack JobCallBack) Queue {
	var wg sync.WaitGroup
	jobQueue := make(chan any, capacity)
	quit := make(chan struct{})
	return Queue{
		Workers:     workers,
		JobQueue:    jobQueue,
		JobCallBack: jobCallBack,
		Wg:          &wg,
		QuitChan:    quit,
	}
}

// Stop close all the running goroutines
// and stops processing any more jobs
func (q *Queue) Stop() {
	q.QuitChan <- struct{}{}
}

// EnqueueJobNonBlocking  use this to queue the jobs you need to execute
// in an unblocking way...(i.e if the [Queue.JobQueue] is full it will not block)
// Returns false if the buffer is full
// else if it is accepted the job it returns true
// use case imagine you are receiving events and you want to prevent anymore
// events from being submitted if the buffered channel
// is full,you can return an error to the user if this function returns false...
// although a better approach would be to store it in redis if is it is rejected for
// later processing once the [JobQueue] has available space to prevent loss of events
// Note  if you are using a for loop to consume the jobs, it's better  to use [Queue.EnqueueJobBlocking ]
// to prevent you from having a Busy wait(continuous pooling to check if a space is available to queue the job)
// which might utilize 90% of your cpu or more.
func (q *Queue) EnqueueJobNonBlocking(job interface{}) bool {
	select {
	case q.JobQueue <- job:
		return true
	default:
		return false
	}
}

// EnqueueJobBlocking queues jobs and blocks if [JobQueue] is full
// once the [JobQueue] is no longer full the job will be accepted
// this is better  for your cpu utilization  unlike [Queue.EnqueueJobNonBlocking]
// when consuming a job via a for loop.
// since it stop looping and blocks if you have no space in the Queue
// it will only continue looping when you have more space in the  Queue
func (q *Queue) EnqueueJobBlocking(job interface{}) {
	q.JobQueue <- job
}

// StartWorkers start  goroutines  and add them to wait group
// the goroutines added(started) will be determined by the number of [Queue.Workers] you specified for this queue
// for example if you specified 10 Workers 10 goroutines will be used to process the job
// they can now start picking jobs and processing them
func (q *Queue) StartWorkers() {
	for i := 0; i < q.Workers; i++ {
		// add the goroutine  to a wait group to prevent the program from exiting
		// be a goroutine return
		q.Wg.Add(1)
		go q.worker()
	}
	q.Wg.Wait()
}

func (q *Queue) worker() {
	defer q.Wg.Done()
	for {
		select {
		// terminate the goroutine
		case <-q.QuitChan:
			log.Println("closing the  workers")
			return
		case job := <-q.JobQueue:
			// a job has been received  call this function
			q.JobCallBack(job)
		}
	}
}
