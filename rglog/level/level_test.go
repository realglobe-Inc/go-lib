package level

import (
	"testing"
)

func TestConvert(t *testing.T) {
	for _, level := range Values() {
		if level2, e := ValueOf(level.String()); e != nil {
			t.Error(level)
		} else if level2 != level {
			t.Error(level)
		}
	}
}
