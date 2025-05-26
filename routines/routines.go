package main

import (
	"fmt"
	"time"
)

func Bar(pipe chan<- int, done chan<- bool) {
	for i := range 30 {
		pipe <- i
		// time.Sleep(1 * time.Second)
	}
	close(pipe)
	done <- true
}

func main() {
	done := make(chan bool)
	pipe := make(chan int, 1)

	go Bar(pipe, done)

	time.Sleep(10 * time.Second)

	go func() {
		for data := range pipe {
			fmt.Println(data)
		}
	}()

	<-done
	fmt.Println("bar call done")
}
