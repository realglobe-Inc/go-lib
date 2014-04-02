package file

import (
	"bytes"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	"github.com/realglobe-Inc/go-lib-rg/run"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"
)

var log rglog.Logger

func init() {
	log = rglog.GetLogger("github.com/realglobe-Inc/go-lib-rg/file")
}

// ファイル・ディレクトリに関するユーティリティ。

const (
	dirPerm  = 0755
	filePerm = 0644
)

func IsExist(path string) (bool, error) {
	_, err := os.Lstat(path) // リンクも検知する。
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
		return 0, erro.Wrap(err)
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

// ディレクトリを登る方向で列挙。
// /a/b/c, /a/b, /a, /
func AscentDirList(path string) []string {
	return DirList(path)
}

// ディレクトリを下る方向で列挙。
func DescentDirList(path string) []string {
	buff := DirList(path)
	list := make([]string, len(buff))
	for i := 0; i < len(buff); i++ {
		list[i] = buff[len(buff)-1-i]
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

func RecursiveCopy(to, from string) error {
	if e := os.MkdirAll(filepath.Dir(to), dirPerm); e != nil {
		return erro.Wrap(e)
	}

	// ひどい手抜き。
	log.Debug("TE NU KI.")
	return erro.Wrap(run.Quiet("cp", "-r", from, to))
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

// ファイルの末尾に行を付け足す。
func AppendLines(path string, lines []string) error {
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

	buff := []byte{}
	if stat.Size() > 1 {
		// 改行で終わっていなければ改行する。
		tail := make([]byte, 1)
		if _, e := writer.ReadAt(tail, stat.Size()-1); e != nil {
			return erro.Wrap(e)
		} else if tail[0] != '\n' {
			buff = append(buff, '\n')
		}
	}
	for _, line := range lines {
		buff = append(buff, []byte(fmt.Sprintln(line))...)
	}

	if _, e := writer.WriteAt(buff, stat.Size()); e != nil {
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

	return ChownById(path, int(uid))
}

func ChownById(path string, uid int) error {
	return ChownByIdGid(path, uid, uid)
}

func ChownByIdGid(path string, uid, gid int) error {
	if e := filepath.Walk(path, func(curPath string, info os.FileInfo, err error) error {
		if err != nil {
			return erro.Wrap(err)
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if e := os.Lchown(curPath, uid, gid); e != nil {
				return erro.Wrap(e)
			}
		} else {
			if e := os.Chown(curPath, uid, gid); e != nil {
				return erro.Wrap(e)
			}
		}

		return nil
	}); e != nil {
		return erro.Wrap(e)
	}

	log.Debug("chown ", uid, ":", gid, " -R ", path)

	return nil
}
