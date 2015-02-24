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
			return nil, nil
		}

		return nil, erro.Wrap(err)
	}

	log.Debug("Locked ", file.Name(), ".")
	return (*Locker)(file), nil
}

func (lock *Locker) Unlock() error {
	file := (*os.File)(lock)
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
		return erro.Wrap(err)
	}
	if err := file.Close(); err != nil {
		return erro.Wrap(err)
	}

	log.Debug("Unlocked ", file.Name(), ".")
	return nil
}

// ロックできるか指定した時間が経つまで待つ。
// ロックできずに指定した時間が経ったら nil を返す。
func WaitLock(path string, waittime time.Duration) (*Locker, error) {

	timeoutCh := time.After(waittime)
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
	case <-timeoutCh:
		ackCh <- false
		return nil, nil
	}
}
