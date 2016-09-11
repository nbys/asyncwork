# asyncwork
Execute tasks concurrently in golang.

## Prepare the tasks
The task is a function with following signature:
```go
type TaskFunction func() interface{}
```

It can be declared like this:
```go
task1 := func() interface{} {
  time.Sleep(time.Second * 4)
  fmt.Println("very slow function")
  return "I'm ready"
}
```

All task have to be collected into the slice:
```go
tasks := []worker.TaskFunction{task1, task2, task3}
```

Initialize the channel which will orchestrate canceling of task execution.
```go
done := make(chan struct{})
defer close(done)
```

## Execute the tasks
Start task execution:
```go
resultChannel := worker.PerformTasks(tasks, done)
```

## Get result from completed tasks
Loop over the channel with results. 
```go
for result := range resultChannel {
		switch result.(type) {
		case string:
			fmt.Println("Here is a string:", result.(string))
		case int:
			fmt.Println("Here is an integer:", result.(int))
```
