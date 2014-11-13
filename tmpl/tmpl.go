package tmpl

import (
	"bytes"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

const dirPerm = 0755

func Generate(destPath, tmplPath string, data interface{}) error {

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return erro.Wrap(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return erro.Wrap(err)
	}

	// 変換に成功したので書き込む。

	if err := os.MkdirAll(filepath.Dir(destPath), dirPerm); err != nil {
		return erro.Wrap(err)
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return erro.Wrap(err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, &buf); err != nil {
		return erro.Wrap(err)
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
		if err := dest.Chmod(tmplFi.Mode()); err != nil {
			return erro.Wrap(err)
		}
		log.Debug("chmod ", tmplFi.Mode(), " ", destPath)
	}

	return nil
}
