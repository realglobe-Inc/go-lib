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
