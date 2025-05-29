package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func baz() error {
	panic("foo")
}

func bar(ctx context.Context) {
	fmt.Println("bar: started")
	timer := time.Tick(1 * time.Second)

	forever := true
	for forever {
		select {
		case <-timer:
			fmt.Println("bar: still alive")
		case <-ctx.Done():
			fmt.Println("bar: received context stop")
			forever = false
		}
	}
	fmt.Println("bar: done")
}

func foo(ctx context.Context) {
	fmt.Println("foo: started")
	go bar(ctx)
	<-ctx.Done()
	fmt.Println("foo: done")
}

func main() {
	fmt.Println("main: started")
	defer func() {
		fmt.Println("after panic?")
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("received interrupt signal")
			cancel()
		}
	}()

	baz()

	go foo(ctx)

	<-ctx.Done()

	time.Sleep(1 * time.Second)
	fmt.Println("main: done")
}
