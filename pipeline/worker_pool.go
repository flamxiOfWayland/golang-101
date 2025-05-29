package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Duration(rand.IntN(5)) * time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {
	const numJobs = 10
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 10; w++ {
		go worker(w, jobs, results)
	}

	fmt.Println("sending jobs")
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	fmt.Println("no more jobs")

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}
