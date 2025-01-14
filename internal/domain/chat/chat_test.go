package chat

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {

	ch := make(chan int)

	// 채널에 값을 보내는 고루틴 실행
	go func() {
		ch <- 123
	}()

	go func() {
		ch <- 456
	}()

	for {
		i := <-ch
		if i != 123 && i != 456 {
			t.Errorf("Expected %d, but got %d", i, i)
		}

		fmt.Printf("i : %d\n", i)
	}
}

func TestUnbufferedChannel(t *testing.T) {
	fmt.Println("Testing with unbuffered channel:")

	// Unbuffered channel
	daUnbuffered := make(chan int) // Unbuffered

	// Start the sender goroutine
	go func(chans chan<- int) {
		defer close(chans)
		fmt.Println("Sending 1...")
		chans <- 1 // Will block if unbuffered, won't block if buffered

		for i := 0; i < 10; i++ {
			fmt.Printf("Sending %d...\n", i)
			chans <- i                         // Will block if unbuffered
			time.Sleep(500 * time.Millisecond) // Simulate some delay in sending
		}
	}(daUnbuffered)

	// Expected values
	expected := []int{1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	index := 0

	// Receive values and validate
	for a := range daUnbuffered {
		if index >= len(expected) {
			t.Errorf("Received unexpected value: %d", a)
			break
		}

		if a != expected[index] {
			t.Errorf("Expected %d, but got %d at index %d", expected[index], a, index)
		}
		index++

		// Simulate slow receiver
		time.Sleep(1 * time.Second)
	}

	if index != len(expected) {
		t.Errorf("Did not receive all expected values. Expected %d, but got %d", len(expected), index)
	}
}

func TestChannelConcurrency(t *testing.T) {

	tests := []struct {
		title        string
		bufferSize   int
		goroutineNum int
	}{
		{
			title:        "buffer size 1",
			bufferSize:   1,
			goroutineNum: 100,
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {

			ch := make(chan int, tc.bufferSize) // 채널 버퍼 크기를 1로 설정

			for i := 0; i < tc.goroutineNum; i++ {
				go func(j int) {
					fmt.Printf("Goroutine %d sent value\n", i+1)
					ch <- j + 1 // 값을 채널에 보냄
				}(i)
			}

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()

				i := 0
				for {
					select {
					case msg := <-ch:
						fmt.Printf("Received value %d\n", msg)
						i++

						if i == tc.goroutineNum {
							fmt.Printf("Complete %d\n", i)
							return
						}
					}
				}
			}()

			wg.Wait()
		})
	}

}
