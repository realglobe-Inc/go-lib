package ip

import (
	"fmt"
	"math"
	"net"
	"testing"
)

func ExampleAdd() {
	ad := Add(net.ParseIP("172.16.0.0"), 1)
	fmt.Println(ad)
	ad = Add(net.ParseIP("172.16.0.0"), 65535)
	fmt.Println(ad)
	ad = Add(net.IP([]byte{172, 16, 0, 0}), 256)
	fmt.Println(ad)
	//Output: 172.16.0.1
	//172.16.255.255
	//172.16.1.0
}

func TestDiff(t *testing.T) {
	ip1, ip2 := net.ParseIP("172.16.0.0"), net.ParseIP("172.16.0.0")
	if d := Diff(ip1, ip2); d != 0 {
		t.Error(d, 0, ip1, ip2)
	}
	ip2 = net.ParseIP("172.16.0.1")
	if d := Diff(ip1, ip2); d != 1 {
		t.Error(d, 1, ip1, ip2)
	}
	ip1 = net.ParseIP("172.16.0.2")
	if d := Diff(ip1, ip2); d != math.MaxUint32 {
		t.Error(d, math.MaxUint32, ip1, ip2)
	}
}
