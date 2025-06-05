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

type TLinuxKernelInfo struct {
}

func (k *TLinuxKernelInfo) ReleaseVersion() string {
	output, _ := util.ExecCmd("cat /etc/tlinux-release")
	return output
}

func (k *TLinuxKernelInfo) KernelVersion() string {
	kernelVer, _ := host.KernelVersion()
	return kernelVer
}

func (k *TLinuxKernelInfo) Uptime() string {
	uptime, _ := host.Uptime()
	return fmt.Sprintf("%d天", uptime/60/60/24)
}

func (k *TLinuxKernelInfo) CpuArch() string {
	arch, _ := host.KernelArch()
	return arch
}

func (k *TLinuxKernelInfo) CpuCores() int {
	counts, _ := cpu.Counts(true)
	return counts
}

func (k *TLinuxKernelInfo) CpuLoad() float32 {
	result, _ := util.ExecCmd("cat /proc/loadavg | awk '{print $2}'")
	loadaverage, _ := strconv.ParseFloat(result, 32)
	return float32(loadaverage)
}

func (k *TLinuxKernelInfo) CpuPercent() string {
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

func (k *TLinuxKernelInfo) PidNum() string {
	processes, err := process.Processes()
	if err != nil {
		return "-"
	}
	return fmt.Sprintf("%d", len(processes))
}

// Memory 单位: Gi
func (k *TLinuxKernelInfo) Memory() int {
	memory, _ := mem.VirtualMemory()
	return int(memory.Total / 1024 / 1024 / 1024)
}

func (k *TLinuxKernelInfo) MemoryPercent() string {
	memory, _ := mem.VirtualMemory()
	return fmt.Sprintf("%0.2f%%", memory.UsedPercent)
}

func (k *TLinuxKernelInfo) SwapTotal() string {
	swapMemory, _ := mem.SwapMemory()
	return fmt.Sprintf("%dKB", swapMemory.Total)
}

func (k *TLinuxKernelInfo) SwapUsed() string {
	swapMemory, _ := mem.SwapMemory()
	return fmt.Sprintf("%dKB", swapMemory.Used)
}

func (k *TLinuxKernelInfo) HugePage() string {
	hugePage, _ := util.ExecCmd("cat /proc/meminfo | grep -i HugePages_Total | awk '{print $2}'")
	return hugePage
}
