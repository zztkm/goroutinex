//	Kevin Chen (2017)
//	Patterns from Pike's Google I/O talk, "Go Concurrency Patterns"

//	Golang restoring sequencing after multiplexing

package main

import (
	"fmt"
	"time"
)

// Message ...
type Message struct {
	str   string
	block chan int
}

// fanIn is itself a generator
func fanIn(ch1, ch2 <-chan Message) <-chan Message { // receives two read-only channels
	newCh := make(chan Message)
	go func() {
		for {
			newCh <- <-ch1
		}
	}() // launch two goroutine while loops to continuously pipe to new channel
	go func() {
		for {
			newCh <- <-ch2
		}
	}()
	return newCh
}

func generator(msg string) <-chan Message { // returns receive-only channel
	ch := make(chan Message)
	blockingStep := make(chan int) // channel within channel to control exec, set false default
	go func() {                    // anonymous goroutine
		for i := 0; ; i++ {
			ch <- Message{fmt.Sprintf("%s %d", msg, i), blockingStep}
			time.Sleep(time.Second)
			blockingStep <- 1 // block by waiting for input
		}
	}()
	return ch
}

func main() {
	ch := fanIn(generator("Hello"), generator("Bye"))
	for i := 0; i < 10; i++ {
		msg1 := <-ch
		fmt.Println(msg1.str)

		msg2 := <-ch
		fmt.Println(msg2.str)

		<-msg1.block // reset channel, stop blocking
		<-msg2.block
	}
}
