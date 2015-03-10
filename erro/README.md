<!--
Copyright 2015 realglobe, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->


# erro

スタックトレース付きエラー。


## 使い方

エラーを投げる方では、

```Go
func g() error {
	...
	v, err := f()
	if err != nil {
		return Wrap(err)
	}
	...
}
```

受ける方では、

```Go
	err := g()
	if err != nil {
		switch e := erro.Unwrap(err).(type) {
		case *net.OpError:
			...
		default:
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
```

こんな感じ。
