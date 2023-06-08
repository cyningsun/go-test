package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	cliNum  = 100
	tps     = 1000
	backoff = time.Second / (tps / cliNum)
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	fmt.Println("backoff", backoff)

	for w := 1; w <= cliNum; w++ {
		go call(w)
	}

	<-c
}

func call(id int) {
	for {
		req, err := http.NewRequest("POST", "http://localhost:8090/hello", bytes.NewBufferString("John Doe"))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/octet-stream")
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1200)
		defer cancel()
		req = req.WithContext(ctx)

		begin := time.Now()
		resp, err := http.DefaultClient.Do(req)
		duration := time.Now().Sub(begin)
		if err != nil {
			fmt.Printf("cli-%d: %s\n", id, err)
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("cli-%d(%dms): %s\n", id, duration/time.Millisecond, string(body))
	}
}
