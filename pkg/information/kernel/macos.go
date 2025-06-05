package kernel

import (
	"fmt"
	"strconv"
	"time"

	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type MacOSKernelInfo struct {
}

func (k *MacOSKernelInfo) ReleaseVersion() string {
	productName, _ := util.ExecCmd("sw_vers -productName")
	productVersion, _ := util.ExecCmd("sw_vers -productVersion")
	return fmt.Sprintf("%s %s", productName, productVersion)
}

func (k *MacOSKernelInfo) KernelVersion() string {
	kernelVer, _ := host.KernelVersion()
	return kernelVer
}

func (k *MacOSKernelInfo) Uptime() string {
	uptime, _ := host.Uptime()
	return fmt.Sprintf("%d天", uptime/60/60/24)
}

func (k *MacOSKernelInfo) CpuArch() string {
	arch, _ := host.KernelArch()
	return arch
}

func (k *MacOSKernelInfo) CpuCores() int {
	counts, _ := cpu.Counts(true)
	return counts
}

func (k *MacOSKernelInfo) CpuLoad() float32 {
	result, _ := util.ExecCmd("sysctl -n vm.loadavg | awk '{print $2}'")
	loadaverage, _ := strconv.ParseFloat(result, 32)
	return float32(loadaverage)
}

func (k *MacOSKernelInfo) CpuPercent() string {
	// 参数：采样间隔1秒，汇总所有核心
	percentages, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return "-"
	}

	if len(percentages) > 0 {
		return fmt.Sprintf("%.2f%%", percentages[0])
	}
	return "-"
}

func (k *MacOSKernelInfo) PidNum() string {
	processes, err := process.Processes()
	if err != nil {
		return "-"
	}
	return fmt.Sprintf("%d", len(processes))
}

// Memory 单位: Gi
func (k *MacOSKernelInfo) Memory() int {
	memory, _ := mem.VirtualMemory()
	return int(memory.Total / 1024 / 1024 / 1024)
}

func (k *MacOSKernelInfo) MemoryPercent() string {
	memory, _ := mem.VirtualMemory()
	return fmt.Sprintf("%0.2f%%", memory.UsedPercent)
}

func (k *MacOSKernelInfo) SwapTotal() string {
	swapMemory, _ := mem.SwapMemory()
	return fmt.Sprintf("%dKB", swapMemory.Total)
}

func (k *MacOSKernelInfo) SwapUsed() string {
	swapMemory, _ := mem.SwapMemory()
	return fmt.Sprintf("%dKB", swapMemory.Used)
}

func (k *MacOSKernelInfo) HugePage() string {
	// MacOS 不支持透明大页，返回 N/A
	return "N/A"
}
