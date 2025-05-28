package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func SimpleProduce(num int) []int {
	var out []int
	for range num {
		value := rand.IntN(1000)
		out = append(out, value)
		time.Sleep(1 * time.Millisecond)
	}
	return out
}

func SimpleDouble(input []int) []int {
	out := []int{}
	for _, value := range input {
		out = append(out, value*2)
		// time.Sleep(1 * time.Millisecond)
	}
	return out
}

func SimpleFilterBelow(thresh int, input []int) []int {
	out := []int{}
	for _, value := range input {
		if value < thresh {
			out = append(out, value)
			// time.Sleep(1 * time.Millisecond)
		}
	}
	return out
}

func SimplePrint(input []int) {
	for range input {
		// fmt.Printf("value is %d\n", value)
	}
}

func main() {
	start := time.Now().UnixMilli()

	input := SimpleProduce(1_000)
	doubled := SimpleDouble(input)
	filtered := SimpleFilterBelow(1000, doubled)
	SimplePrint(filtered)

	fmt.Printf("total time: %d mili\n", time.Now().UnixMilli()-start)
}
