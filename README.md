# asyncwork
Execute tasks concurrently in golang.

## Install
``` shell
go get -u github.com/nbys/asyncwork/worker
```

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

Use context.Context to stop running goroutines
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

## Execute the tasks
Start task execution:
```go
resultChannel := worker.PerformTasks(ctx, tasks)
```

## Get result from completed tasks
Loop over the channel with results.
```go
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
```
