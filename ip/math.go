package ip

import (
	"bytes"
	"encoding/binary"
	"net"
)

// IPv4 アドレスに整数値を足す。
func Add(ip net.IP, diff uint32) net.IP {
	return fromVal(val(ip) + diff)
}

// IPv4 アドレスの差を整数値で返す。
func Diff(from, to net.IP) uint32 {
	return val(to) - val(from)
}

func val(ip net.IP) uint32 {
	buff := []byte(ip)
	buff = buff[len(buff)-4:] // 基本 16 バイト (IPv6) で 4 バイト (IPv4) の場合は後ろ詰めのようだから。

	var val uint32
	if e := binary.Read(bytes.NewReader(buff), binary.BigEndian, &val); e != nil {
		panic(e) // 実装ミス。
	}

	return val
}

func fromVal(val uint32) net.IP {
	var buff bytes.Buffer
	if e := binary.Write(&buff, binary.BigEndian, &val); e != nil {
		panic(e) // 実装ミス。
	}

	return net.IP(buff.Bytes())
}
