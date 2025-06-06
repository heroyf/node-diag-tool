package net

import (
	"runtime"

	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/sirupsen/logrus"
)

type IptablesCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &IptablesCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *IptablesCheck) RunCheck() *plugin.CheckPluginResult {
	result, err := util.ExecCmd(c.LookupMatchedScript())
	if err != nil {
		logrus.Debugf("iptables check failed: %v \n", result)
		o := util.ParseShellResult(result)
		return c.BuildBlockResult(c, o.ResultMsg, result)
	}
	return c.BuildPassResult(c, "检测iptables正常", result)
}

func (c *IptablesCheck) PluginName() string {
	return c.BuildPluginName(c.FileName)
}

func (c *IptablesCheck) Author() string {
	return "heroyf"
}
