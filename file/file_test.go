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
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	if ok, err := IsExist(file.Name()); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Error(file.Name())
	}

	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}

	if ok, err := IsExist(file.Name()); err != nil {
		t.Fatal(err)
	} else if ok {
		t.Error(file.Name())
	}
}

func TestIsDir(t *testing.T) {
	file, err := ioutil.TempFile("", "test_file")
	if err != nil {
		t.Fatal(err)
	} else if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	dir, err := ioutil.TempDir("", "test_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)

	if ok, err := IsDir(file.Name()); err != nil {
		t.Fatal(err)
	} else if ok {
		t.Error(file.Name())
	}

	if ok, err := IsDir(dir); err != nil {
		t.Fatal(err)
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
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}

	if cmp, err := Compare(files[0].Name(), files[1].Name()); err != nil {
		t.Fatal(err)
	} else if cmp != 0 {
		t.Error(files, cmp)
	}
	for i := 2; i < len(files); i++ {
		if cmp, err := Compare(files[0].Name(), files[i].Name()); err != nil {
			t.Fatal(err)
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
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}

	if cmp, err := Compare(files[0].Name(), files[1].Name()); err != nil {
		t.Fatal(err)
	} else if cmp == 0 {
		t.Error(files, cmp)
	}

	if err := Copy(files[1].Name(), files[0].Name()); err != nil {
		t.Fatal(err)
	}

	if cmp, err := Compare(files[0].Name(), files[1].Name()); err != nil {
		t.Fatal(err)
	} else if cmp != 0 {
		t.Error(files, cmp)
	}
}

func TestAppend(t *testing.T) {
	file, err := ioutil.TempFile("", "test_file")
	if err != nil {
		t.Fatal(err)
	} else if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if err := Append(file.Name(), []byte("あいうえお")); err != nil {
		t.Fatal(err)
	}

	if b, err := ioutil.ReadFile(file.Name()); err != nil {
		t.Fatal(err)
	} else if string(b) != "あいうえお" {
		t.Error(string(b))
	}

	if err := Append(file.Name(), []byte("かきくけこ")); err != nil {
		t.Fatal(err)
	}

	if b, err := ioutil.ReadFile(file.Name()); err != nil {
		t.Fatal(err)
	} else if string(b) != "あいうえおかきくけこ" {
		t.Error(string(b))
	}
}

func TestAppendLines(t *testing.T) {
	file, err := ioutil.TempFile("", "test_file")
	if err != nil {
		t.Fatal(err)
	} else if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if err := ioutil.WriteFile(file.Name(), []byte("あいうえお"), filePerm); err != nil {
		t.Fatal(err)
	}

	if b, err := ioutil.ReadFile(file.Name()); err != nil {
		t.Fatal(err)
	} else if string(b) != "あいうえお" {
		t.Fatal(string(b))
	}

	if err := AppendLines(file.Name(), []string{"かきくけこ"}); err != nil {
		t.Fatal(err)
	}

	if b, err := ioutil.ReadFile(file.Name()); err != nil {
		t.Fatal(err)
	} else if string(b) != "あいうえお\nかきくけこ\n" {
		t.Error(string(b))
	}

	if err := AppendLines(file.Name(), []string{"さしすせそ", "たちつてと"}); err != nil {
		t.Fatal(err)
	}

	if b, err := ioutil.ReadFile(file.Name()); err != nil {
		t.Fatal(err)
	} else if string(b) != "あいうえお\nかきくけこ\nさしすせそ\nたちつてと\n" {
		t.Error(string(b))
	}
}
