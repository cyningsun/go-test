package main

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"
)

func main() {
	f, _ := os.Create("trace.out")
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	RunMyProgram()
}

func RunMyProgram() {
	ctx := context.Background()
	myID := "123"

	ctx, task := trace.NewTask(ctx, "Say Hello")
	defer task.End()
	//Log Something on the Task
	trace.Log(ctx, "myID", myID)
	saidHello := make(chan bool)
	//Say Hello in a goroutine using WithRegion
	go func() {
		trace.WithRegion(ctx, "sayHello", sayHello)
		saidHello <- true
	}()
	//another way to create a region
	<-saidHello
	trace.StartRegion(ctx, "sayGoodbye").End()
	sayGoodBye()
}
func sayHello()   { fmt.Println("Hello Trace") }
func sayGoodBye() { fmt.Println("goodbye") }
