package etcd

import (
	"github.com/astaxie/beego/logs"
	"net"
)

var (
	LocalIpArray []string
)

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logs.Error("get local addr failed")
		panic("get local addr failed")
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			ip := ipNet.IP.To4()
			if ip != nil {
				LocalIpArray = append(LocalIpArray, ip.String())
			}
		}
	}
	logs.Debug(LocalIpArray)
	logs.Info("Get LocalIp success!")
}
