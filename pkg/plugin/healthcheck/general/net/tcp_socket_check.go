package net

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/heroyf/node-diag-tool/pkg/consts"
	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/sirupsen/logrus"
)

var (
	// tcpSocketConnMaxThrehold 需要检测的TCP链接最大阈值, 超过最大阈值, 需要报异常
	tcpSocketConnMaxThrehold = map[string]int{
		"ESTAB":      300000,
		"TIME-WAIT":  30000,
		"CLOSE-WAIT": 10000,
	}
)

type TcpSocketCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &TcpSocketCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *TcpSocketCheck) RunCheck() *plugin.CheckPluginResult {
	result, err := util.ExecCmd("ss -nat | grep -E -v 'grep|State' | awk '{print $1}' | sort | uniq -c | sort -n")
	if err != nil {
		logrus.Debugf("tcp socket check failed: %v \n", err)
		return c.BuildPassResult(c, consts.ScriptFailed, "")
	}

	lines := c.ReadLines(result)
	for _, line := range lines {
		arr := strings.Split(strings.TrimSpace(line), " ")
		curValue, err := strconv.Atoi(arr[0])
		if err != nil {
			logrus.Debugf("tcp socket parse number failed: %v \n", err)
			continue
		}

		// skt链接名, 移除结尾的'0'
		sktConnName := strings.TrimSuffix(strings.TrimSpace(arr[1]), "0")
		// 当前值大于临界值, 则返回检测异常
		if threhold, exits := tcpSocketConnMaxThrehold[sktConnName]; exits && curValue > threhold {
			return c.BuildBlockResult(c,
				fmt.Sprintf("检测tcp链接异常: [%s]当前链接数: [%d], 超过阈值: [%d]", arr[1], curValue, threhold),
				"")
		}
	}

	return c.BuildPassResult(c, "检测tcp链接无明显异常", result)
}

func (c *TcpSocketCheck) PluginName() string {
	return c.BuildPluginName(c.FileName)
}

func (c *TcpSocketCheck) Author() string {
	return "heroyf"
}
