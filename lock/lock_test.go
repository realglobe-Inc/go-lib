package lock

import (
	"github.com/realglobe-Inc/go-lib/erro"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestLockConcurrency(t *testing.T) {
	file, err := ioutil.TempFile("", "go-lib-test")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(file.Name())
	path := file.Name()

	if err := file.Close(); err != nil {
		t.Error(err)
		return
	} else if err := os.Remove(path); err != nil {
		t.Error(err)
		return
	}

	counter := 0
	errCh := make(chan error)

	n := 100
	loop := 100
	for i := 0; i < n; i++ {
		go func(id int) {
			for j := 0; j < loop; j++ {
				lock, err := Lock(path)
				if err != nil {
					errCh <- erro.New(err)
					return
				}

				counter++

				if err := lock.Unlock(); err != nil {
					errCh <- erro.New(err)
					return
				}
			}
			errCh <- nil
		}(i)
	}

	// 終了待ち。
	for i := 0; i < n; i++ {
		if err := <-errCh; err != nil {
			t.Error(err)
		}
	}

	if counter != n*loop {
		t.Error(counter, n*loop)
	}
}

func TestTryLockConcurrency(t *testing.T) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	path := file.Name()

	if err := file.Close(); err != nil {
		t.Fatal(err)
	} else if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	n := 100
	loop := 100

	counter := 0

	succCh := make(chan int)
	errCh := make(chan error)
	for i := 0; i < n; i++ {
		go func(id int) {
			succ := 0
			defer func() { succCh <- succ }()

			for j := 0; j < loop; j++ {
				lock, err := TryLock(path)
				if err != nil {
					errCh <- erro.Wrap(err)
					return
				} else if lock == nil {
					continue
				}

				counter++
				succ++

				if err := lock.Unlock(); err != nil {
					errCh <- erro.Wrap(err)
					return
				}
			}

			errCh <- nil
		}(i)
	}

	// 終了待ち。
	for i := 0; i < n; i++ {
		if err := <-errCh; err != nil {
			t.Error(err)
		}
	}

	sum := 0
	for i := 0; i < n; i++ {
		sum += <-succCh
	}

	if counter == 0 {
		t.Error(counter)
	} else if counter != sum {
		t.Error(counter, sum)
	}
}

func TestWaitLockConcurrency(t *testing.T) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	path := file.Name()

	if err := file.Close(); err != nil {
		t.Fatal(err)
	} else if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	n := 100
	loop := 100
	wait := 10 * time.Millisecond

	counter := 0

	succCh := make(chan int)
	errCh := make(chan error)
	for i := 0; i < n; i++ {
		go func(id int) {
			succ := 0
			defer func() { succCh <- succ }()

			for j := 0; j < loop; j++ {
				lock, err := WaitLock(path, wait)
				if err != nil {
					errCh <- erro.Wrap(err)
					return
				} else if lock == nil {
					continue
				}

				counter++
				succ++

				if err := lock.Unlock(); err != nil {
					errCh <- erro.Wrap(err)
					return
				}
			}

			errCh <- nil
		}(i)
	}

	// 終了待ち。
	for i := 0; i < n; i++ {
		if err := <-errCh; err != nil {
			t.Error(err)
		}
	}

	sum := 0
	for i := 0; i < n; i++ {
		sum += <-succCh
	}

	if counter == 0 {
		t.Error(counter)
	} else if counter != sum {
		t.Error(counter, sum)
	}
}

func BenchmarkLock(b *testing.B) {
	file, err := ioutil.TempFile("", "go-lib-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(file.Name())
	path := file.Name()

	if err := file.Close(); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		lock, err := Lock(path)
		if err != nil {
			b.Error(err)
			return
		}

		if err := lock.Unlock(); err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkTryLock(b *testing.B) {
	file, err := ioutil.TempFile("", "go-lib-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(file.Name())
	path := file.Name()

	if err := file.Close(); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		lock, err := TryLock(path)
		if err != nil {
			b.Error(err)
			return
		} else if lock == nil {
			b.Error("lock failed")
			return
		}

		if err := lock.Unlock(); err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkWaitLock(b *testing.B) {
	file, err := ioutil.TempFile("", "go-lib-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(file.Name())
	path := file.Name()

	if err := file.Close(); err != nil {
		b.Fatal(err)
	}

	wait := 10 * time.Millisecond
	for i := 0; i < b.N; i++ {
		lock, err := WaitLock(path, wait)
		if err != nil {
			b.Error(err)
			return
		} else if lock == nil {
			b.Error("lock failed")
			return
		}

		if err := lock.Unlock(); err != nil {
			b.Error(err)
			return
		}
	}
}

func TestReentrant(t *testing.T) {
	file, err := ioutil.TempFile("", "go-lib-test")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(file.Name())
	path := file.Name()

	if err := file.Close(); err != nil {
		t.Error(err)
		return
	} else if err := os.Remove(path); err != nil {
		t.Error(err)
		return
	}

	lock, err := Lock(path)
	if err != nil {
		t.Error(err)
		return
	}
	defer lock.Unlock()

	lock, err = TryLock(path)
	if err != nil {
		t.Error(err)
		return
	} else if lock != nil {
		lock.Unlock()
		t.Error(lock)
		return
	}

	lock, err = WaitLock(path, time.Millisecond)
	if err != nil {
		t.Error(err)
		return
	} else if lock != nil {
		lock.Unlock()
		t.Error(lock)
		return
	}
}

func TestTimeout(t *testing.T) {
	file, err := ioutil.TempFile("", "test_lock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	path := file.Name()

	if err := file.Close(); err != nil {
		t.Fatal(err)
	} else if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	lock, err := Lock(path)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Unlock()

	wait := 10 * time.Millisecond

	start := time.Now()
	lock, err = WaitLock(path, wait)
	dur := time.Since(start)
	if err != nil {
		t.Fatal(err)
	} else if lock != nil {
		lock.Unlock()
		t.Error(lock)
	} else if dur < wait || wait+time.Second < dur {
		t.Error(dur)
	}
}
