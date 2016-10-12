package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func worker(i int) {
	defer func() {
		fmt.Fprintf(os.Stderr, "Work B %d done.\n", i)
	}()

	fmt.Fprintf(os.Stderr, "Got B %d\n", i)
	resp, err := http.Get("http://www.google.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Work B %d got error: %v\n", i, err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Work B %d got error: %v\n", i, err)
		return
	}

	fmt.Fprintf(os.Stderr, "Work B %d got data, len %d\n", i, len(data))
}

func main() {
	chanA := make(chan string, 10)
	tickerB := time.NewTicker(time.Second).C
	tickerC := time.NewTicker(3 * time.Second).C

	doneChan := make(chan bool)
	go func() {
		time.Sleep(time.Minute)
		doneChan <- true
	}()

	go func() {
		for i := 0; i < 25; i++ {
			chanA <- "A"
			fmt.Fprintf(os.Stderr, "Sent A\n")
			time.Sleep(300 * time.Millisecond)
		}
	}()

	i := 0
	for {
		select {
		case a := <-chanA:
			fmt.Fprintf(os.Stderr, "Got %s\n", a)
		case <-tickerB:
			i++
			worker(i)
		case <-tickerC:
			fmt.Fprintf(os.Stderr, "Got C\n")
		case <-doneChan:
			return
		}
	}

}
