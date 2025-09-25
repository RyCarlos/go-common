package network

import "net"

// GetLocalIP 获取本地IP的函数
func GetLocalIP() string {
	defaultIp := "127.0.0.1"
	address, err := net.InterfaceAddrs()
	if err != nil {
		return defaultIp
	}

	for _, address := range address {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return defaultIp
}
