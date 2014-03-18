package lock

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLock(t *testing.T) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if e := file.Close(); e != nil {
		t.Fatal(e)
	}

	n := 10
	loop := 1000

	counter := 0

	ch := make(chan int)
	lockPath := file.Name()
	for i := 0; i < n; i++ {
		go func(id int) {
			defer func() { ch <- 0 }()

			for j := 0; j < loop; j++ {
				var lock *Locker
				for lock == nil {
					var err error
					lock, err = Lock(lockPath)
					if err != nil {
						t.Error(id, j, err)
						return
					}
				}

				counter++
				if e := lock.Unlock(); e != nil {
					t.Error(id, j, e)
					return
				}
			}
		}(i)
	}

	// 終了待ち。
	for i := 0; i < n; i++ {
		<-ch
	}

	if counter != n*loop {
		t.Error(counter, n, loop, file.Name())
	}
}
