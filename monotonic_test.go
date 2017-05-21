package redsync

import (
	"fmt"
	"testing"
	"time"
)

func TestMonotonicTime(t *testing.T) {
	orderCh := make(chan int)
  flag := 0
  // pools := vmPools(5)
	pools := newMockPools(8)
  rs := New(pools)

  mutex := rs.NewMutex("test-redsync")

  sleepTime := [2]int{10, 0}
  for idx, _ := range sleepTime {
		fmt.Println("index is %d", idx)

		go func(idx int) {
			err := mutex.Lock()
			if err != nil {
				t.Fatalf("Expected err == nil, got %q", err)
			}

			fmt.Println("========index is %d", idx)

			if idx == 0 {
				// network partition
				fmt.Println("it actually is gonna sleep")
				time.Sleep(10 * time.Second)
			}

      // write something to the disk
      flag = idx

			defer mutex.Unlock()
			orderCh <- idx
		}(idx)
	}

	for range sleepTime {
		<-orderCh
	}

  if flag != 1 {
		fmt.Println("receive %d", flag)
    t.Fatalf("1 is expected")
  }
}
