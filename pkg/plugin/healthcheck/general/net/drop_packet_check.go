package net

import (
	"runtime"

	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/sirupsen/logrus"
)

type DropPackageCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &DropPackageCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *DropPackageCheck) RunCheck() *plugin.CheckPluginResult {
	result, err := util.ExecCmd("ip -s -s link show eth1")
	if err != nil {
		logrus.Debugf("drop packet check failed: %v \n", result)
		o := util.ParseShellResult(result)
		return c.BuildBlockResult(c, o.ResultMsg, result)
	}
	return c.BuildPassResult(c, "检测无明显网络丢包", result)
}

func (c *DropPackageCheck) PluginName() string {
	return c.BuildPluginName(c.FileName)
}

func (c *DropPackageCheck) Author() string {
	return "heroyf"
}
