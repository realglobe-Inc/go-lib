// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package level

import (
	"flag"
	"testing"
)

func TestCompare(t *testing.T) {
	lvs := Values()
	for i := 0; i < len(lvs); i++ {
		for j := 0; j < len(lvs); j++ {
			if i == j {
				if lvs[i].Higher(lvs[j]) {
					t.Fatal(lvs[i], lvs[j])
				} else if lvs[i].Lower(lvs[j]) {
					t.Fatal(lvs[i], lvs[j])
				}
			} else if i < j {
				if !lvs[i].Higher(lvs[j]) {
					t.Fatal(lvs[i], lvs[j])
				} else if lvs[i].Lower(lvs[j]) {
					t.Fatal(lvs[i], lvs[j])
				}
			} else {
				if lvs[i].Higher(lvs[j]) {
					t.Fatal(lvs[i], lvs[j])
				} else if !lvs[i].Lower(lvs[j]) {
					t.Fatal(lvs[i], lvs[j])
				}
			}
		}
	}
}

func TestConvert(t *testing.T) {
	for _, lv := range Values() {
		if lv2, err := ValueOf(lv.String()); err != nil {
			t.Fatal(lv)
		} else if lv2 != lv {
			t.Fatal(lv)
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
			t.Fatal(lv)
		}
	}
}
