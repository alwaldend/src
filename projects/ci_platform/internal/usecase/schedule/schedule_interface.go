package useSchedule

import "sync"

type JobCallBack func(job any)

type Options struct {
	// number of goroutines(workers,consumers) to be used to process the jobs
	Workers int
	// number of items that will be held in the JobQueue channel (its capacity) at a time before blocking
	// i.e capacity of our JobQueue channel
	Capacity int
	// JobQueue a buffered channel of capacity [Queue.Capacity] that will temporary hold
	// jobs before they are assigned to a worker
	JobQueue chan any
	// Wg will be used to make sure the program does not terminate before all of our goroutines
	// complete the job assigned to them
	Wg *sync.WaitGroup
	// QuitChan will be used to stop all goroutines [Queue.Workers]
	QuitChan chan struct{}
	// JobCallBack is the function to be called when a job (event) is received
	// it should implement JobCallBack i.e a function or method with only one parameter and no return value
	JobCallBack JobCallBack
}

type Queue struct {
	// Workers Number of goroutines(workers,consumers) to be used to process the jobs
	Workers int
	// Capacity is the number of  items that will be held in the JobQueue channel (its capacity) at a time before blocking
	// i.e capacity of our JobQueue channel
	Capacity int
	// JobQueue a buffered channel of capacity [Queue.Capacity] that will temporary hold
	// jobs before they are assigned to a worker
	JobQueue chan any
	// Wg will be used to make sure the program does not terminate before all of our goroutines
	// complete the job assigned to them
	Wg *sync.WaitGroup
	// QuitChan will be used to stop all goroutines [Queue.Workers]
	QuitChan chan struct{}
	// JobCallBack is the function to be called when a job (event) is received
	// it should implement JobCallBack i.e a function or method with only one parameter and no return value
	JobCallBack JobCallBack
}

type Job struct {
	Run      func(data any) error
	Callback func() error
}

type Interactor interface{}
