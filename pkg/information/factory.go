package information

import (
	"runtime"

	"github.com/heroyf/node-diag-tool/pkg/information/kernel"
)

// NewKernelInfo 根据操作系统返回对应的内核信息实现
func NewKernelInfo() kernel.BaseMachineInfo {
	switch runtime.GOOS {
	case "darwin":
		return &kernel.MacOSKernelInfo{}
	case "linux":
		return &kernel.TLinuxKernelInfo{}
	default:
		return &kernel.TLinuxKernelInfo{} // 默认使用 Linux 实现
	}
}
