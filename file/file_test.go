package file

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIsExist(t *testing.T) {
	file, err := ioutil.TempFile("", "test_file")
	if err != nil {
		t.Fatal(err)
	}
	if e := file.Close(); e != nil {
		t.Fatal(e)
	}

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

func TestIsDir(t *testing.T) {
	file, err := ioutil.TempFile("", "test_file")
	if err != nil {
		t.Fatal(err)
	} else if e := file.Close(); e != nil {
		t.Fatal(e)
	}
	defer os.Remove(file.Name())

	dir, err := ioutil.TempDir("", "test_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)

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

func TestCompare(t *testing.T) {
	files := make([]*os.File, 5)
	for i := 0; i < len(files); i++ {
		var err error
		files[i], err = ioutil.TempFile("", "test_file")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(files[i].Name())
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
}

func TestCopy(t *testing.T) {
	files := make([]*os.File, 2)
	for i := 0; i < len(files); i++ {
		var err error
		files[i], err = ioutil.TempFile("", "test_file")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(files[i].Name())
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
}

func TestAppend(t *testing.T) {
	file, err := ioutil.TempFile("", "test_file")
	if err != nil {
		t.Fatal(err)
	} else if e := file.Close(); e != nil {
		t.Fatal(e)
	}
	defer os.Remove(file.Name())

	if e := Append(file.Name(), []byte("あいうえお")); e != nil {
		t.Fatal(e)
	}

	if b, e := ioutil.ReadFile(file.Name()); e != nil {
		t.Fatal(e)
	} else if string(b) != "あいうえお" {
		t.Error(string(b))
	}

	if e := Append(file.Name(), []byte("かきくけこ")); e != nil {
		t.Fatal(e)
	}

	if b, e := ioutil.ReadFile(file.Name()); e != nil {
		t.Fatal(e)
	} else if string(b) != "あいうえおかきくけこ" {
		t.Error(string(b))
	}
}

func TestAppendLines(t *testing.T) {
	file, err := ioutil.TempFile("", "test_file")
	if err != nil {
		t.Fatal(err)
	} else if e := file.Close(); e != nil {
		t.Fatal(e)
	}
	defer os.Remove(file.Name())

	if e := ioutil.WriteFile(file.Name(), []byte("あいうえお"), filePerm); e != nil {
		t.Fatal(e)
	}

	if b, e := ioutil.ReadFile(file.Name()); e != nil {
		t.Fatal(e)
	} else if string(b) != "あいうえお" {
		t.Fatal(string(b))
	}

	if e := AppendLines(file.Name(), []string{"かきくけこ"}); e != nil {
		t.Fatal(e)
	}

	if b, e := ioutil.ReadFile(file.Name()); e != nil {
		t.Fatal(e)
	} else if string(b) != "あいうえお\nかきくけこ\n" {
		t.Error(string(b))
	}

	if e := AppendLines(file.Name(), []string{"さしすせそ", "たちつてと"}); e != nil {
		t.Fatal(e)
	}

	if b, e := ioutil.ReadFile(file.Name()); e != nil {
		t.Fatal(e)
	} else if string(b) != "あいうえお\nかきくけこ\nさしすせそ\nたちつてと\n" {
		t.Error(string(b))
	}
}
