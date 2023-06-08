package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Message struct {
	enqueue  time.Time
	playload string
}

var (
	workerNum = 1
	queueSize = 1000
	jobs      = make(chan string, queueSize)
	results   = make(chan string, queueSize)
)

func worker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		results <- fmt.Sprintf("worker-%d: hello %s", id, string(j))
		time.Sleep(time.Millisecond * 50)
	}
}

func dispatch(w http.ResponseWriter, req *http.Request) {
	rBody, _ := ioutil.ReadAll(req.Body)

	select {
	case jobs <- string(rBody):
		fmt.Println("queue len:", len(jobs))
	default:
		fmt.Fprintf(w, "enqueue error")
		return
	}

	// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	// defer cancel()

	select {
	case resp := <-results:
		fmt.Fprintf(w, resp)
		// case <-ctx.Done():
		// 	fmt.Fprintf(w, "timeout error")
	}
}

func main() {
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	http.HandleFunc("/hello", dispatch)

	http.ListenAndServe(":8090", nil)
}
