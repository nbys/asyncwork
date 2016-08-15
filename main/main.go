package main

import (
	"asyncwork/worker"
	"errors"
	"fmt"
	"time"
)

func doSomeWork() {

	// Define three slow functions which will be runnning concurrently
	slowFunction := func() worker.Result {
		fmt.Println("slow function")
		time.Sleep(time.Second * 2)
		return worker.Result{Val: 2}
	}

	verySlowFunction := func() worker.Result {
		fmt.Println("very slow function")
		time.Sleep(time.Second * 4)
		return worker.Result{Val: "I'm ready"}
	}

	// One function returns an error
	errorFunction := func() worker.Result {
		fmt.Println("function with an error")
		return worker.Result{Err: errors.New("Error in function")}
	}

	tasks := []worker.TaskFunction{errorFunction, slowFunction, verySlowFunction}

	// The channel to signal other goroutines don't send a value
	done := make(chan struct{})
	defer close(done)

	resultChannel := worker.PerformTasks(tasks, done)

	for result := range resultChannel {
		if result.Err != nil {
			// An error has been received, skip furher processing
			fmt.Println(result.Err, time.Now())
			return
		} else {
			switch result.Val.(type) {
			case string:
				fmt.Println("Here is a string:", result.Val.(string))
			case int:
				fmt.Println("Here is an integer:", result.Val.(int))
			default:
				fmt.Println("Some unknown type ")
			}
		}
	}

}

func main() {
	doSomeWork()
	time.Sleep(time.Second * 8)
}
