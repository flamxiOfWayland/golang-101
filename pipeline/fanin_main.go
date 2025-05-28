package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func gen(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}
func fanInSimple(cs ...<-chan string) <-chan string {
	c := make(chan string)
	for _, ci := range cs {
		go func(cv <-chan string) {
			for {
				input := <-cv
				c <- input
			}
		}(ci)

	}
	return c
}

func main() {
	c := fanInSimple(
		gen("Joe1"), gen("Ahn1"),
		gen("Joe2"), gen("Ahn2"),
		gen("Joe3"), gen("Ahn3"),
		gen("Joe4"), gen("Ahn4"),
		gen("Joe5"), gen("Ahn5"),
		gen("Joe6"), gen("Ahn6"),
		gen("Joe7"), gen("Ahn7"),
	)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 2; i++ {
			fmt.Printf("remote printer: %s\n", <-c)
		}
		wg.Done()
	}()

	for i := 0; i < 2; i++ {
		fmt.Println(<-c)
	}

	wg.Wait()
	fmt.Println("Done")
}
