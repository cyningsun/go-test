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
	trace.Log(ctx, "myID", myID)
	fmt.Println("Hello Trace")
}
