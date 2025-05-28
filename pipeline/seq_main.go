package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func produce(num int) chan int {
	out := make(chan int)
	go func() {
		for range num {
			out <- rand.IntN(1000)
			time.Sleep(1 * time.Millisecond)
		}
		close(out)
	}()
	return out
}

func double(input <-chan int) chan int {
	out := make(chan int)
	go func() {
		for value := range input {
			out <- value * 2
			time.Sleep(1 * time.Millisecond)
		}
		close(out)
	}()
	return out
}

func filterBelow(thresh int, input <-chan int) chan int {
	out := make(chan int)
	go func() {
		for value := range input {
			if value < thresh {
				out <- value
				time.Sleep(1 * time.Millisecond)
			}
		}
		close(out)
	}()
	return out
}

func print(input <-chan int) {
	for range input {
		// fmt.Printf("value is %d\n", value)
	}
}

func main() {
	start := time.Now().UnixMilli()
	input := produce(1_000)
	doubled := double(input)
	filtered := filterBelow(1000, doubled)
	print(filtered)
	fmt.Printf("total time: %d mili\n", time.Now().UnixMilli()-start)
}
