package util

import (
	"errors"
	"fmt"
	"net"
)

func GetLocalIp() (string, error) {
	// 优先检查bond1接口
	if ip, err := getInterfaceIp("bond1"); err == nil {
		return ip, nil
	}

	// 其次检查eth1接口
	if ip, err := getInterfaceIp("eth1"); err == nil {
		return ip, nil
	}

	// 回退到遍历所有接口
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("获取网络接口失败: %w", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("未找到符合条件的IPv4地址")
}

func getInterfaceIp(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", fmt.Errorf("接口 %s 不存在: %w", interfaceName, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("获取接口地址失败: %w", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipv4 := ipnet.IP.To4(); ipv4 != nil {
				return ipv4.String(), nil
			}
		}
	}

	return "", fmt.Errorf("接口 %s 未配置有效IPv4地址", interfaceName)
}
