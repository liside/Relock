package redlock

import (
	"fmt"
	"testing"
	"time"
)

func TestMonotonicTime(t *testing.T) {
  flag := 0
  pools := vmPools(5)
  rs := New(pools)

  mutex := rs.NewMutex("test-redsync")

  sleepTime := [10, 0]
  for idx, st := range sleepTime {
		go func() {
			err := mutex.Lock()
			if err != nil {
				t.Fatalf("Expected err == nil, got %q", err)
			}

			if i == 0 {
				time.Sleep(st * time.Second)
			}

      // write something to the disk
      flag = idx


			defer mutex.Unlock()
			orderCh <- i
		}()
	}

	for range mutexes {
		<-orderCh
	}

  if flag != 2 {
    t.Fatalf("2 is expected")
  }
}

func vmPools() (n int) []Pool {
	pools := []Pool{}
	for index, server := range servers {
		func() {
			pools = append(pools, &redis.Pool{
				MaxIdle:     3,
				IdleTimeout: 240 * time.Second,
				Dial: func() (redis.Conn, error) {
					return redis.Dial("tcp", fmt.Sprintf("localhost:%d", 8000 + idx))
				},
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			})
		}()
		if len(pools) == n {
			break
		}
	}
	return pools
}
