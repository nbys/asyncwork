package worker

import "sync"

// There is my implementation of the "pipeling". The original idea is described in
// the "Go Blog": https://blog.golang.org/pipelines

// Result is a structure which will contain either value or error from task execution.
type Result struct {
	Val interface{}
	Err error
}

// TaskFunction is a function type for tasks to be performed.
// All incoming tasks have to conform to this function type.
type TaskFunction func() Result

// PerformTasks is a function which will be called by the client to perform
// multiple task concurrently.
// Input:
// tasks: the slice with functions (type TaskFunction)
// done:  the channel to trigger the end of task processing and return
// Output: the channel with results
func PerformTasks(tasks []TaskFunction, done chan struct{}) chan Result {

	// Create a worker for each incoming task
	workers := make([]chan Result, 0, len(tasks))

	for _, task := range tasks {
		resultChannel := newWorker(task, done)
		workers = append(workers, resultChannel)
	}

	// Merge results from all workers
	out := merge(workers, done)
	return out
}

func newWorker(task TaskFunction, done chan struct{}) chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)

		select {
		case <-done:
			// Received a signal to abandon furher processing
			return
		case out <- task():
			// Got some result
		}
	}()

	return out
}

func merge(workers []chan Result, done chan struct{}) chan Result {
	// Merged channel with results
	out := make(chan Result)

	// Synchronization over channels: do not close "out" before all tasks are completed
	var wg sync.WaitGroup

	// Start output goroutine for each outbound channel from the workers
	// get all values from channel (c) before channel is closed
	// if interruption signal has received decrease the counter of running tasks via wg.Done()
	output := func(c <-chan Result) {
		defer wg.Done()
		for result := range c {
			select {
			case <-done:
				// Received a signal to abandon furher processing
				return
			case out <- result:
				// some message or nothing
			}
		}
	}

	wg.Add(len(workers))
	for _, workerChannel := range workers {
		go output(workerChannel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
