package lock

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}

	n := 10
	loop := 1000

	counter := 0

	timeout := time.After(time.Duration(int64(n*loop*1000) * int64(time.Microsecond)))

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
				if err := lock.Unlock(); err != nil {
					t.Error(id, j, err)
					return
				}
			}
		}(i)
	}

	// 終了待ち。
	for i := 0; i < n; i++ {
		select {
		case <-ch:
		case <-timeout:
			t.Fatal("Dead lock?")
		}
	}

	if counter != n*loop {
		t.Error(counter, n, loop, file.Name())
	}
}

func BenchmarkLock(b *testing.B) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(file.Name())
	if err := file.Close(); err != nil {
		b.Fatal(err)
	}

	timeout := time.After(time.Duration(int64(b.N*1000) * int64(time.Microsecond)))
	ch := make(chan error)
	lockPath := file.Name()
	go func() {
		for j := 0; j < b.N; j++ {
			var lock *Locker

			for lock == nil {
				var err error
				lock, err = Lock(lockPath)
				if err != nil {
					ch <- err
					return
				}
			}

			if err := lock.Unlock(); err != nil {
				ch <- err
				return
			}
		}

		ch <- nil
	}()

	// 終了待ち。
	select {
	case err := <-ch:
		if err != nil {
			b.Fatal(err)
		}
	case <-timeout:
		b.Fatal("Dead lock?")
	}
}

func TestReentrant(t *testing.T) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}

	lock, err := Lock(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Unlock()

	lock, err = TryLock(file.Name())
	if err != nil {
		t.Fatal(err)
	} else if lock != nil {
		defer lock.Unlock()
		t.Error(lock)
	}

	lock, err = WaitLock(file.Name(), time.Nanosecond)
	if err != nil {
		t.Fatal(err)
	} else if lock != nil {
		defer lock.Unlock()
		t.Error(lock)
	}

}

func TestWaitLock(t *testing.T) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}

	n := 10
	loop := 1000

	counter := 0

	timeout := time.After(time.Duration(int64(n*loop*1000) * int64(time.Microsecond)))

	ch := make(chan int)
	lockPath := file.Name()
	for i := 0; i < n; i++ {
		go func(id int) {
			defer func() { ch <- 0 }()

			for j := 0; j < loop; j++ {
				var lock *Locker
				for lock == nil {
					var err error
					lock, err = WaitLock(lockPath, time.Second)
					if err != nil {
						t.Error(id, j, err)
						return
					}
				}

				counter++
				if err := lock.Unlock(); err != nil {
					t.Error(id, j, err)
					return
				}
			}
		}(i)
	}

	// 終了待ち。
	for i := 0; i < n; i++ {
		select {
		case <-ch:
		case <-timeout:
			t.Fatal("Dead lock?")
		}
	}

	if counter != n*loop {
		t.Error(counter, n, loop, file.Name())
	}
}

func TestTimeout(t *testing.T) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}

	lock, err := Lock(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Unlock()

	start := time.Now()
	lock, err = WaitLock(file.Name(), 4*time.Millisecond+500*time.Microsecond)
	if err != nil {
		t.Fatal(err)
	} else if lock != nil {
		defer lock.Unlock()
		t.Error(lock)
	} else if dur := time.Since(start); dur <= 4*time.Millisecond || 5*time.Millisecond <= dur {
		t.Error(dur)
	}

}

func BenchmarkWaitLock(b *testing.B) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(file.Name())
	if err := file.Close(); err != nil {
		b.Fatal(err)
	}

	timeout := time.After(time.Duration(int64(b.N*1000) * int64(time.Microsecond)))
	ch := make(chan error)
	lockPath := file.Name()
	go func() {
		for j := 0; j < b.N; j++ {
			var lock *Locker

			for lock == nil {
				var err error
				lock, err = WaitLock(lockPath, time.Second)
				if err != nil {
					ch <- err
					return
				}
			}

			if err := lock.Unlock(); err != nil {
				ch <- err
				return
			}
		}

		ch <- nil
	}()

	// 終了待ち。
	select {
	case err := <-ch:
		if err != nil {
			b.Fatal(err)
		}
	case <-timeout:
		b.Fatal("Dead lock?")
	}
}
