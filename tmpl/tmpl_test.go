package tmpl

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	tmpl, err := ioutil.TempFile("", "daiku_tmpl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpl.Name())

	if _, err := tmpl.WriteString("" +
		"私は{{.セリフ1}}と言った。\n" +
		"彼は{{.セリフ2}}と答えた。\n" +
		"私は怒った ({{.オチ}})。\n" +
		""); err != nil {
		t.Fatal(err)
	}
	if err := tmpl.Close(); err != nil {
		t.Fatal(err)
	}

	dest, err := ioutil.TempFile("", "daiku_tmpl")
	if err != nil {
		t.Fatal(err)
	}
	if err := dest.Close(); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dest.Name())

	data := map[string]string{
		"セリフ1": "「やあやあ我こそは。」",
		"セリフ2": "「ぬるぽ。」",
		"オチ":   "性的な意味で"}

	if err := Generate(dest.Name(), tmpl.Name(), data); err != nil {
		t.Fatal(err)
	}

	output, err := ioutil.ReadFile(dest.Name())
	if err != nil {
		t.Fatal(err)
	}

	answer := "" +
		"私は「やあやあ我こそは。」と言った。\n" +
		"彼は「ぬるぽ。」と答えた。\n" +
		"私は怒った (性的な意味で)。\n"
	if string(output) != answer {
		t.Error("[", string(output), "][", answer, "]")
	}
}
