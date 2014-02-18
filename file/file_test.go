package file

import (
	"github.com/realglobe-Inc/daiku/common"
	"io/ioutil"
	"os"
	"testing"
)

func TestIsExist(t *testing.T) {
	file, err := ioutil.TempFile("", common.TestLabel)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	if ok, e := IsExist(file.Name()); e != nil {
		t.Fatal(e)
	} else if !ok {
		t.Error(file.Name())
	}

	if e := os.Remove(file.Name()); e != nil {
		t.Fatal(e)
	}

	if ok, e := IsExist(file.Name()); e != nil {
		t.Fatal(e)
	} else if ok {
		t.Error(file.Name())
	}
}

func TestCompare(t *testing.T) {
	files := make([]*os.File, 5)
	for i := 0; i < len(files); i++ {
		var err error
		files[i], err = ioutil.TempFile("", common.TestLabel)
		if err != nil {
			t.Fatal(err)
		}
	}

	files[0].WriteString("1234567890")
	files[1].WriteString("1234567890")
	files[2].WriteString("123456789")
	files[3].WriteString("123456789a")
	files[4].WriteString("1234567890a")

	for _, file := range files {
		if e := file.Close(); e != nil {
			t.Fatal(e)
		}
	}

	if cmp, e := Compare(files[0].Name(), files[1].Name()); e != nil {
		t.Fatal(e)
	} else if cmp != 0 {
		t.Error(files, cmp)
	}
	for i := 2; i < len(files); i++ {
		if cmp, e := Compare(files[0].Name(), files[i].Name()); e != nil {
			t.Fatal(e)
		} else if cmp == 0 {
			t.Error(files, cmp)
		}
	}

	for _, file := range files {
		if e := os.Remove(file.Name()); e != nil {
			t.Fatal(e)
		}
	}

}

func TestCopy(t *testing.T) {
	files := make([]*os.File, 2)
	for i := 0; i < len(files); i++ {
		var err error
		files[i], err = ioutil.TempFile("", common.TestLabel)
		if err != nil {
			t.Fatal(err)
		}
	}

	files[0].WriteString("1234567890abcdefg")

	for _, file := range files {
		if e := file.Close(); e != nil {
			t.Fatal(e)
		}
	}

	if cmp, e := Compare(files[0].Name(), files[1].Name()); e != nil {
		t.Fatal(e)
	} else if cmp == 0 {
		t.Error(files, cmp)
	}

	if e := Copy(files[1].Name(), files[0].Name()); e != nil {
		t.Fatal(e)
	}

	if cmp, e := Compare(files[0].Name(), files[1].Name()); e != nil {
		t.Fatal(e)
	} else if cmp != 0 {
		t.Error(files, cmp)
	}

	for _, file := range files {
		if e := os.Remove(file.Name()); e != nil {
			t.Fatal(e)
		}
	}
}

func TestDirList(t *testing.T) {
	path := "/a/b/c"
	if dirs := DirList(path); len(dirs) != 3 || dirs[0] != "/a/b" || dirs[1] != "/a" || dirs[2] != "/" {
		t.Error(path, dirs)
	}
	path = "a/b/c"
	if dirs := DirList(path); len(dirs) != 3 || dirs[0] != "a/b" || dirs[1] != "a" || dirs[2] != "." {
		t.Error(path, dirs)
	}
}

func TestIsDir(t *testing.T) {
	file, err := ioutil.TempFile("", common.TestLabel)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	dir, err := ioutil.TempDir("", common.TestLabel)
	if err != nil {
		t.Fatal(err)
	}

	if ok, e := IsDir(file.Name()); e != nil {
		t.Fatal(e)
	} else if ok {
		t.Error(file.Name())
	}

	if ok, e := IsDir(dir); e != nil {
		t.Fatal(e)
	} else if !ok {
		t.Error(dir)
	}

	if e := os.Remove(file.Name()); e != nil {
		t.Fatal(e)
	}
	if e := os.Remove(dir); e != nil {
		t.Fatal(e)
	}
}
