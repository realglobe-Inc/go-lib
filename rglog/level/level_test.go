package level

import (
	"flag"
	"testing"
)

func TestConvert(t *testing.T) {
	for _, lv := range Values() {
		if lv2, err := ValueOf(lv.String()); err != nil {
			t.Error(lv)
		} else if lv2 != lv {
			t.Error(lv)
		}
	}
}

func TestVar(t *testing.T) {
	for _, lv := range Values() {
		flags := flag.NewFlagSet("test", flag.ExitOnError)
		var lv2 Level
		flags.Var(Var(&lv2, INFO), "lv", "Log level.")
		flags.Parse([]string{"-lv", lv.String()})
		if lv2 != lv {
			t.Error(lv)
		}
	}
}
