package etcd

import (
	"github.com/astaxie/beego/logs"
	"net"
)

var (
	localIpArray []string
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
				localIpArray = append(localIpArray, ip.String())
			}
		}
	}
	logs.Debug(localIpArray)
	logs.Warn("Get LocalIp success!")
}
