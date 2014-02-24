package ip

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	ipmath "github.com/realglobe-Inc/go-lib-rg/ip"
	"github.com/realglobe-Inc/go-lib-rg/log"
	"net"
	"os/user"
	"strconv"
)

// コンテナと IP アドレスの対応。
// コンテナ名と同じユーザーがシステムに存在するとして、IP アドレスの始点にその uid を足してコンテナの IP アドレスとする。
// 通常は uid を 60000 以内でうまく割り当ててくれるようなので (/etc/login.defs 参照)、IP アドレスの範囲は 16 ビットあると良い。

// コンテナの IP アドレスを返す。
// コンテナユーザーが存在しないときは nil を返す。
func FromName(name string, ips *net.IPNet) (net.IP, error) {

	usr, err := user.Lookup(name)
	if err != nil {
		if _, ok := err.(user.UnknownUserError); ok {
			log.Debug("User ", name, " was not found.")
			return nil, nil
		}
		return nil, erro.Wrap(err)
	}

	maskSize, _ := ips.Mask.Size()
	uid, err := strconv.ParseUint(usr.Uid, 10, 32)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if uid >= 1<<uint(maskSize) {
		return nil, erro.New("Uid ", uid, " is larger than ip range ", 1<<uint(maskSize), ".")
	}

	ip := ipmath.Add(ips.IP, uint32(uid))

	log.Debug("IP ", ip, "/", maskSize, " was selected.")
	return ip, nil
}

// IP アドレスからコンテナ名を返す。
// コンテナユーザーが存在しないときは空文字列を返す。
func Name(ip net.IP, ips *net.IPNet) (string, error) {
	uid := ipmath.Diff(ips.IP, ip)

	usr, err := user.LookupId(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		if _, ok := err.(user.UnknownUserIdError); ok {
			log.Debug("User was not found for ", ip, ".")
			return "", nil
		}
		return "", erro.Wrap(err)
	}

	log.Debug("User ", usr.Username, " was found for ", ip, ".")
	return usr.Username, nil
}
