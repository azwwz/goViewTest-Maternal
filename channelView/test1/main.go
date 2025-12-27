package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Demo: channels are concurrency-safe for send/receive")
	basicSafeSendReceive()

	fmt.Println("\nDemo: close must be coordinated (unsafe if multiple closers)")
	unsafeCloseExample()

	fmt.Println("\nDemo: safe close with single owner")
	safeCloseExample()
}

func basicSafeSendReceive() {
	ch := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		defer wg.Done()
		for v := range ch {
			fmt.Printf("recv %d\n", v)
		}
	}()

	wg.Wait()
}

// Demonstrates a race on close: multiple goroutines attempt to close.
// This will likely panic: "close of closed channel".
func unsafeCloseExample() {
	ch := make(chan int)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("panic recovered:", r)
			}
		}()
		close(ch)
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("panic recovered:", r)
			}
		}()
		close(ch)
	}()

	// Give goroutines time to run.
	time.Sleep(50 * time.Millisecond)
}

// Safe close pattern: single owner closes, others only send/receive.
func safeCloseExample() {
	ch := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Printf("recv %d\n", v)
	}

	wg.Wait()
}
