package net

import (
	"runtime"

	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/sirupsen/logrus"
)

// EthCpuAffinityCheck 网卡多队列绑核检测
type EthCpuAffinityCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &EthCpuAffinityCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *EthCpuAffinityCheck) RunCheck() *plugin.CheckPluginResult {
	result, err := util.ExecCmd(c.LookupMatchedScript())
	if err != nil {
		logrus.Debugf("eth cpu affinity check failed: %v \n", result)
		o := util.ParseShellResult(result)
		return c.BuildBlockResult(c, o.ResultMsg, result)
	}
	return c.BuildPassResult(c, "检测网卡多队列绑核正常", result)
}

func (c *EthCpuAffinityCheck) PluginName() string {
	return c.BuildPluginName(c.FileName)
}

func (c *EthCpuAffinityCheck) Author() string {
	return "heroyf"
}
