package tmpl

import (
	"bytes"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

var log rglog.Logger

func init() {
	log = rglog.GetLogger("github.com/realglobe-Inc/go-lib-rg/tmpl")
}

const dirPerm = 0755

func Generate(destPath, tmplPath string, data interface{}) error {

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return erro.Wrap(err)
	}

	var buf bytes.Buffer
	if e := tmpl.Execute(&buf, data); e != nil {
		return erro.Wrap(e)
	}

	// 変換に成功したので書き込む。

	if e := os.MkdirAll(filepath.Dir(destPath), dirPerm); e != nil {
		return erro.Wrap(e)
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return erro.Wrap(err)
	}
	defer dest.Close()

	if _, e := io.Copy(dest, &buf); e != nil {
		return erro.Wrap(e)
	}

	log.Debug("tmpl ", tmplPath, " ", destPath)

	// パーミッションを合わせる。
	tmplFi, err := os.Stat(tmplPath)
	if err != nil {
		return erro.Wrap(err)
	}

	destFi, err := dest.Stat()
	if err != nil {
		return erro.Wrap(err)
	}

	if destFi.Mode() != tmplFi.Mode() {
		if e := dest.Chmod(tmplFi.Mode()); e != nil {
			return erro.Wrap(e)
		}
		log.Debug("chmod ", tmplFi.Mode(), " ", destPath)
	}

	return nil
}
