package run

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"os"
	"testing"
)

func tempFileName() (string, error) {
	temp, err := ioutil.TempFile("", "run_test")
	if err != nil {
		return "", erro.Wrap(err)
	}
	if err := temp.Close(); err != nil {
		return "", erro.Wrap(err)
	}
	if err := os.Remove(temp.Name()); err != nil {
		return "", erro.Wrap(err)
	}
	return temp.Name(), nil
}

func isExist(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, erro.Wrap(err)
		}
	}
	return true, nil
}

func TestRun(t *testing.T) {
	temp, err := tempFileName()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(temp)
	if err := Run("touch", temp); err != nil {
		t.Fatal(err)
	}
	if exist, err := isExist(temp); err != nil {
		t.Fatal(err)
	} else if !exist {
		t.Error(temp)
	}
}

func TestNonInteractive(t *testing.T) {
	temp, err := tempFileName()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(temp)
	if err := NonInteractive("touch", temp); err != nil {
		t.Fatal(err)
	}
	if exist, err := isExist(temp); err != nil {
		t.Fatal(err)
	} else if !exist {
		t.Error(temp)
	}
}

func TestNeglect(t *testing.T) {
	temp, err := tempFileName()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(temp)
	Neglect("touch", temp)
	if exist, err := isExist(temp); err != nil {
		t.Fatal(err)
	} else if !exist {
		t.Error(temp)
	}
}

func TestOutput(t *testing.T) {
	stdout, stderr, err := Output("echo", "aho")
	if err != nil {
		t.Fatal(err)
	}
	if stdout != "aho" || stderr != "" {
		t.Error("stdout[", stdout, "], stderr[", stderr, "]")
	}
}

func TestStdout(t *testing.T) {
	stdout, err := Stdout("echo", "aho")
	if err != nil {
		t.Fatal(err)
	}
	if stdout != "aho" {
		t.Error(stdout)
	}
}

func TestQuiet(t *testing.T) {
	temp, err := tempFileName()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(temp)
	if err := Quiet("touch", temp); err != nil {
		t.Fatal(err)
	}
	if exist, err := isExist(temp); err != nil {
		t.Fatal(err)
	} else if !exist {
		t.Error(temp)
	}
}
