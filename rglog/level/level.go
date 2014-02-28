package level

import ()

type Level int

const (
	// 初期値 (0) なら何も出力しない。
	ERR Level = iota + 1
	WARN
	INFO
	DEBUG
)

func (level Level) String() string {
	switch level {
	case ERR:
		return "ERR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
