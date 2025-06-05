package util

import (
	"fmt"
	"os"
	"strings"
)

func KernelValue(key string) string {
	value, _ := os.ReadFile(key)
	return strings.TrimSpace(string(value))
}

// SysctlToProcPath 将sysctl参数名转换为/proc文件路径
// 示例：net.ipv4.tcp_timestamps → /proc/sys/net/ipv4/tcp_timestamps
func SysctlToProcPath(param string) (string, error) {
	if param == "" {
		return "", fmt.Errorf("参数不能为空")
	}

	// 替换点号为路径分隔符
	pathPart := strings.Replace(param, ".", "/", -1)
	// 拼接完整路径
	fullPath := fmt.Sprintf("/proc/sys/%s", pathPart)
	return fullPath, nil
}
