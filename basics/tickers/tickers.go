package main

import "time"
import "fmt"

func dosomething(t time.Time) {
	fmt.Println("doing ", t)
	time.Sleep(2 * time.Second)
	fmt.Println("done ", t)
}

func main() {
	tickChan := time.NewTicker(time.Second).C

	doneChan := make(chan bool)
	go func() {
		time.Sleep(time.Second * 10)
		doneChan <- true
	}()

	for {
		select {
		case <-tickChan:
			fmt.Println("Ticker ticked")
			dosomething(time.Now())
		case <-doneChan:
			fmt.Println("Done")
			return
		}
	}
}
