// plugin/general/cpu/cpu_load_check.go
package cpu

import (
	"fmt"
	"runtime"
	"time"

	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
)

// CpuSysUserCheck CPU内核态、用户态时间检测
type CpuSysUserCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &CpuSysUserCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *CpuSysUserCheck) RunCheck() *plugin.CheckPluginResult {
	// 第一次采样
	times1, _ := cpu.Times(false)
	time.Sleep(1 * time.Second) // 等待间隔
	// 第二次采样
	times2, _ := cpu.Times(false)

	// 计算增量
	userDelta := times2[0].User - times1[0].User
	sysDelta := times2[0].System - times1[0].System
	idleDelta := times2[0].Idle - times1[0].Idle
	totalDelta := times2[0].Total() - times1[0].Total()

	// 计算比例（百分比）
	userPercent := (userDelta / totalDelta) * 100
	sysPercent := (sysDelta / totalDelta) * 100
	idlePercent := (idleDelta / totalDelta) * 100

	logrus.Debugf("用户态占用: %.2f%%, 内核态占用: %.2f%%, 空闲率占用: %.2f%% \n",
		userPercent, sysPercent, idlePercent)

	stdout := fmt.Sprintf("内核态占用: [%.2f], 用户态占用: [%.2f]", sysPercent, userPercent)
	// 内核态大于用户态, 且内核态大于10%, 则认为异常
	if sysPercent > userPercent && sysPercent > 10 {
		return c.BuildBlockResult(c, "检测CPU内核态占用异常",
			stdout)
	}
	return c.BuildPassResult(c, "检测CPU内核态占用正常", stdout)
}

func (c *CpuSysUserCheck) PluginName() string {
	return c.BuildPluginName(c.FileName)
}

func (c *CpuSysUserCheck) Author() string {
	return "heroyf"
}
