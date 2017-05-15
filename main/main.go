package main

import (
	"asyncwork/worker"
	"context"
	"errors"
	"fmt"
	"time"
)

func doSomeWork() {

	// Define three slow functions which will be runnning concurrently
	slowFunction := func() interface{} {
		time.Sleep(time.Second * 10)
		fmt.Println("slow function")
		return 2
	}

	verySlowFunction := func() interface{} {
		time.Sleep(time.Second * 5)
		fmt.Println("very slow function")
		return "I'm ready"
	}

	// One function returns an error
	errorFunction := func() interface{} {
		time.Sleep(time.Second * 1)
		fmt.Println("function with an error")
		return errors.New("Error in function")
	}

	tasks := []worker.TaskFunction{slowFunction, verySlowFunction, errorFunction}

	// Use context to cancel goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChannel := worker.PerformTasks(ctx, tasks)

	// Print value from first goroutine and cancel others
	for result := range resultChannel {
		switch result.(type) {
		case error:
			fmt.Println("Received error")
			cancel()
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
