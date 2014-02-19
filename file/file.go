package file

import (
	"bytes"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/log"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"
)

// ファイル・ディレクトリに関するユーティリティ。

const (
	dirPerm  = 0755
	filePerm = 0644
)

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, erro.Wrap(err)
		}
	}
	return true, nil
}

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, erro.Wrap(err)
	}
	return fi.IsDir(), nil
}

func Size(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return -1, erro.Wrap(err)
	}
	return fi.Size(), nil
}

// /a/b/c から [ /a/b /a / ] をつくる
func DirList(path string) []string {
	list := []string{}
	for i := 0; i < 10; i++ {
		newPath := filepath.Dir(path)
		if newPath == path {
			break
		}
		path = newPath
		list = append(list, path)
	}
	return list
}

// ファイルの中身を丸ごと比較する。
func Compare(path1, path2 string) (int, error) {
	fi1, err := os.Stat(path1)
	if err != nil {
		return 0, erro.Wrap(err)
	}
	fi2, err := os.Stat(path2)
	if err != nil {
		return 0, erro.Wrap(err)
	}
	if fi1.Size() < fi2.Size() {
		return -1, nil
	} else if fi1.Size() > fi2.Size() {
		return 1, nil
	}

	// 読み込みつつ逐次比較するように改善できるが、めんどう。
	bytes1, err := ioutil.ReadFile(path1)
	if err != nil {
		return 0, erro.Wrap(err)
	}
	bytes2, err := ioutil.ReadFile(path2)
	if err != nil {
		return 0, erro.Wrap(err)
	}
	return bytes.Compare(bytes1, bytes2), nil
}

// ファイルまたはディレクトリの名前をてきとうに変更する。ディレクトリ間の移動はしない。
// 返り値は新しい名前。ディレクトリは含まない。
func Escape(path, suffix string) (newName string, err error) {
	name := filepath.Base(path)
	date := time.Now()
	tag := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%09d", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond())
	newName = name + "." + tag + suffix
	newPath := filepath.Join(filepath.Dir(path), newName)

	if e := os.Rename(path, newPath); e != nil {
		return "", erro.Wrap(e)
	}

	log.Debug("mv ", path, " ", newPath)

	return newName, nil
}

// 引数の順番は io.Copy にならった。
func Copy(to, from string) error {
	buff, err := ioutil.ReadFile(from)
	if err != nil {
		return erro.Wrap(err)
	}

	// パーミッションも取得。
	fi, err := os.Stat(from)
	if err != nil {
		return erro.Wrap(err)
	}

	if e := os.MkdirAll(filepath.Dir(to), dirPerm); e != nil {
		return erro.Wrap(e)
	}

	if e := ioutil.WriteFile(to, buff, fi.Mode()); e != nil {
		return erro.Wrap(e)
	}

	log.Debug("cp ", from, " ", to)

	return nil
}

// ファイルの末尾に付け足す。
func Append(path string, data []byte) error {
	if e := os.MkdirAll(filepath.Dir(path), dirPerm); e != nil {
		return erro.Wrap(e)
	}

	writer, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, filePerm)
	if err != nil {
		return erro.Wrap(err)
	}
	defer writer.Close()

	stat, err := writer.Stat()
	if err != nil {
		return erro.Wrap(err)
	}

	if _, e := writer.WriteAt(data, stat.Size()); e != nil {
		return erro.Wrap(e)
	}

	return erro.Wrap(err)
}

// path 以下を全部 owner のものにする。
func Chown(path, owner string) error {
	user, err := user.Lookup(owner)
	if err != nil {
		return erro.Wrap(err)
	}

	uid, err := strconv.ParseUint(user.Uid, 10, 32)
	if err != nil {
		return erro.Wrap(err)
	}

	if e := filepath.Walk(path, func(curPath string, info os.FileInfo, err error) error {
		if err != nil {
			return erro.Wrap(err)
		}

		if e := os.Chown(curPath, int(uid), int(uid)); e != nil {
			return erro.Wrap(e)
		}

		return nil
	}); e != nil {
		return erro.Wrap(err)
	}

	return nil
}