package net

import (
	"fmt"

	"github.com/heroyf/node-diag-tool/pkg/consts"
	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/sirupsen/logrus"

	"runtime"
	"strconv"
	"strings"
)

type TcpPressureCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &TcpPressureCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *TcpPressureCheck) RunCheck() *plugin.CheckPluginResult {
	tcpMem := util.KernelValue(consts.TcpMem)
	s := strings.Fields(tcpMem)
	if s == nil || len(s) != 3 {
		return c.BuildPassResult(c, consts.ScriptFailed, "")
	}

	// 内存压力阈值, 单位: 页
	tcpPressure, _ := strconv.Atoi(s[1])
	logrus.Debugf("tcp pressure: %dMB \n", tcpPressure)

	// 读取当前tcp内存内存占用(单位: 页)
	result, err := util.ExecCmd("cat /proc/net/sockstat|grep TCP|grep -i mem|awk '{print $11}'")
	if err != nil {
		logrus.Debugf("parse socketstat failed: %v \n", err)
		return c.BuildPassResult(c, consts.ScriptFailed, "")
	}

	totalMem, _ := strconv.Atoi(result)
	logrus.Debugf("parse /proc/net/sockstat tcp mem: %dMB \n", totalMem)

	stdout := fmt.Sprintf("TCP压力阈值: [%.2f]MB, 当前TCP内存占用: [%.2f]MB",
		float64(tcpPressure*consts.Page)/float64(1024.0), float64(totalMem*consts.Page)/float64(1024.0))
	// 当前内存使用量(页)
	if totalMem > tcpPressure {
		return c.BuildBlockResult(c,
			fmt.Sprintf("检测TCP内存压力异常(共占用%.2fMB)", float64(totalMem*consts.Page)/1024.0),
			stdout,
		)
	}
	return c.BuildPassResult(c,
		fmt.Sprintf("检测TCP内存压力正常(共占用%.2fMB)", float64(totalMem*consts.Page)/1024.0),
		stdout)
}

func (c *TcpPressureCheck) PluginName() string {
	return c.BuildPluginName(c.FileName)
}

func (c *TcpPressureCheck) Author() string {
	return "heroyf"
}
