package main

import (
	"asyncwork/worker"
	"errors"
	"fmt"
	"time"
)

func doSomeWork() {

	// Define three slow functions which will be runnning concurrently
	slowFunction := func() interface{} {
		time.Sleep(time.Second * 2)
		fmt.Println("slow function")
		return 2
	}

	verySlowFunction := func() interface{} {
		time.Sleep(time.Second * 4)
		fmt.Println("very slow function")
		return "I'm ready"
	}

	// One function returns an error
	errorFunction := func() interface{} {
		time.Sleep(time.Second * 3)
		fmt.Println("function with an error")
		return errors.New("Error in function")
	}

	tasks := []worker.TaskFunction{slowFunction, verySlowFunction, errorFunction}

	// The channel to signal other goroutines don't send a value
	done := make(chan struct{})
	defer close(done)

	resultChannel := worker.PerformTasks(tasks, done)

	for result := range resultChannel {
		switch result.(type) {
		case error:
			return
		case string:
			fmt.Println("Here is a string:", result.(string))
		case int:
			fmt.Println("Here is an integer:", result.(int))
		default:
			fmt.Println("Some unknown type ")
		}
	}
}

func main() {
	doSomeWork()
}
