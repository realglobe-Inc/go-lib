package lock

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	"os"
	"syscall"
)

var log rglog.Logger

func init() {
	log = rglog.GetLogger("github.com/realglobe-Inc/go-lib-rg/lock")
}

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

	if e := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|flag); e != nil {
		file.Close()

		if e == syscall.EWOULDBLOCK {
			return nil, nil
		}

		return nil, erro.Wrap(err)
	}

	log.Debug("Locked ", file.Name(), ".")
	return (*Locker)(file), nil
}

func (lock *Locker) Unlock() error {
	file := (*os.File)(lock)
	err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	if err != nil {
		return erro.Wrap(err)
	}
	file.Close()

	log.Debug("Unlocked ", file.Name(), ".")
	return nil
}
