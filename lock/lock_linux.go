package lock

import (
	"github.com/realglobe-Inc/go-lib/erro"
	"os"
	"syscall"
	"time"
)

// ロックファイル式のロック。

type Locker os.File

// ロックするまで待つ。
func Lock(path string) (*Locker, error) {
	return lock(path, 0)
}

// ロックできなかったら nil を返す。
func TryLock(path string) (*Locker, error) {
	return lock(path, syscall.LOCK_NB)
}

func lock(path string, flag int) (*Locker, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|flag); err != nil {
		file.Close()

		if err == syscall.EWOULDBLOCK {
			// ロックできなかった。
			return nil, nil
		}

		return nil, erro.Wrap(err)
	}

	return (*Locker)(file), nil
}

// 解放する。
func (lock *Locker) Unlock() error {
	file := (*os.File)(lock)
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
		return erro.Wrap(err)
	}
	if err := file.Close(); err != nil {
		return erro.Wrap(err)
	}

	return nil
}

// ロックできるか指定した時間が経つまで待つ。
// ロックできずに指定した時間が経ったら nil を返す。
func WaitLock(path string, wait time.Duration) (*Locker, error) {

	timer := time.NewTimer(wait)
	defer timer.Stop()
	lockerCh := make(chan *Locker, 1)
	errCh := make(chan error, 1)
	ackCh := make(chan bool, 1)
	go func() {
		locker, err := Lock(path)
		if err != nil {
			errCh <- err
			return
		}

		lockerCh <- locker
		if <-ackCh {
			// 受け取ってもらえた。
			return
		}

		// 受け取ってもらえなかった。
		locker.Unlock()
	}()

	select {
	case err := <-errCh:
		return nil, err
	case locker := <-lockerCh:
		ackCh <- true
		return locker, nil
	case <-timer.C:
		ackCh <- false
		return nil, nil
	}
}
