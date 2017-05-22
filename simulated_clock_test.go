package redsync

import (
	"fmt"
	"testing"
	"time"
		"github.com/garyburd/redigo/redis"

)

func TestSkewedClock(t *testing.T) {
  orderCh := make(chan int)
  //flag := -1
  pools := newVMPools(6)
  rs := New(pools[:5])

  mutex := rs.NewMutex("test-redsync")

  sleepTime := [2]int{10, 0}
  for idx, _ := range sleepTime {
		go func(idx int) {
			err := mutex.Lock()
			if err != nil {
				t.Fatalf("Expected err == nil, got %q", err)
			}

			if idx == 0 {
				// network partition, to test against an 8 seconds expiry
				time.Sleep(10 * time.Second)
			}

			// check if the clock expireds or not
			if SkewedNow().Before(mutex.until) {
		    	// write the index of the routine to the redis, simulate two servers
		    	fmt.Println(idx, "is writing..")
		    	// connect to a redis instance
		    	conn := pools[5].Get()
				value, err := redis.String(conn.Do("SET", "WHO", idx))
				conn.Close()
				if err != nil && err != redis.ErrNil {
					panic(err)
				}

				if value == "OK"  {
   					fmt.Println(idx, "is writing successfully.")
  				} else {
  					fmt.Println(idx, "is writing unsuccessfully.")
  				}
		    	//flag = idx
			}

			defer mutex.Unlock()
			orderCh <- idx
		}(idx)
	}

	for range sleepTime {
		<-orderCh
	}

	conn := pools[5].Get()
	value, err := redis.String(conn.Do("GET", "WHO"))
	conn.Close()
	if err != nil && err != redis.ErrNil {
		panic(err)
	}

	fmt.Println(value, "wrote to the database last.")

  	if value != "1"  {
   		t.Fatalf("1 is expected to write last.")
  	}
}

func SkewedNow() time.Time {
	// to simulate a 5 seconds offset from different machines
	return time.Now().Add( -5 * time.Second)
}
