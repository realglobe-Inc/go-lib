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
	if e := temp.Close(); e != nil {
		return "", erro.Wrap(e)
	}
	if e := os.Remove(temp.Name()); e != nil {
		return "", erro.Wrap(e)
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
	if e := Run("touch", temp); e != nil {
		t.Fatal(e)
	}
	if exist, e := isExist(temp); e != nil {
		t.Fatal(e)
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
	if e := NonInteractive("touch", temp); e != nil {
		t.Fatal(e)
	}
	if exist, e := isExist(temp); e != nil {
		t.Fatal(e)
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
	if exist, e := isExist(temp); e != nil {
		t.Fatal(e)
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
	if e := Quiet("touch", temp); e != nil {
		t.Fatal(e)
	}
	if exist, e := isExist(temp); e != nil {
		t.Fatal(e)
	} else if !exist {
		t.Error(temp)
	}
}
