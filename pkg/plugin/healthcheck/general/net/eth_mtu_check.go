package net

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/heroyf/node-diag-tool/pkg/consts"
	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/sirupsen/logrus"
)

// EthMtuCheck 网卡最大传输单位检测
type EthMtuCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &EthMtuCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *EthMtuCheck) RunCheck() *plugin.CheckPluginResult {
	result, err := util.ExecCmd("ip -s -s link|grep -i eth|grep -i mtu|awk '{print $5}'")
	if err != nil {
		logrus.Debugf("run check failed: %v \n", err)
		return c.BuildBlockResult(c, consts.ScriptFailed, "")
	}

	lines := c.ReadLines(result)
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.TrimSpace(line) != "1500" {
			return c.BuildBlockResult(c, "检测网卡MTU默认值被篡改", fmt.Sprintf("当前MTU: %s", line))
		}
	}

	return c.BuildPassResult(c, "检测网卡MTU正常", "当前MTU: 1500")
}

func (c *EthMtuCheck) PluginName() string {
	return c.BuildPluginName(c.FileName)
}

func (c *EthMtuCheck) Author() string {
	return "heroyf"
}
