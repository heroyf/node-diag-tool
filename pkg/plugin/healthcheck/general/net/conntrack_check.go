package net

import (
	"fmt"

	"os"
	"runtime"
	"strings"

	"github.com/heroyf/node-diag-tool/pkg/consts"
	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/sirupsen/logrus"
)

type ConntrackCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &ConntrackCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *ConntrackCheck) RunCheck() *plugin.CheckPluginResult {
	conntrackCount := util.ToInt(util.KernelValue(consts.ConntrackCount))
	conntrackMax := util.ToInt(util.KernelValue(consts.ConntrackMax))

	stdout := fmt.Sprintf("conntrack当前值: [%d], 最大值: [%d]", conntrackCount, conntrackMax)
	if float32(conntrackCount/conntrackMax) > 0.85 {
		return c.BuildBlockResult(c, "conntrack分配率>85%", stdout)
	}

	_, err := os.Stat("/usr/sbin/conntrack")
	if err != nil {
		_, yumErr := util.ExecCmd("yum install conntrack -y")
		if yumErr != nil {
			logrus.Debug("yum install conntrack failed.")
			return c.BuildPassResult(c, "conntrack命令安装失败(跳过检查)", stdout)
		}
	}

	// 检测conntrack表是否存在insert failed
	output, _ := util.ExecCmd("conntrack -S|grep -i 'insert_failed'|awk '{print $6}'|awk -F '=' '{print $2}'")
	if util.IsNotEmpty(output) {
		lines := c.ReadLines(output)
		for _, line := range lines {
			// 计数不等于0, 则认为异常
			if strings.TrimSpace(line) != "0" {
				return c.BuildBlockResult(c, "检测conntrack异常(存在insert_failed,会导致网络丢包)", "")
			}
		}
	}
	return c.BuildPassResult(c, "检测conntrack正常", stdout)
}

func (c *ConntrackCheck) PluginName() string {
	return c.BuildPluginName(c.FileName)
}

func (c *ConntrackCheck) Author() string {
	return "heroyf"
}
