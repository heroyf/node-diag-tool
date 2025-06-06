package net

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/heroyf/node-diag-tool/pkg/consts"
	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/heroyf/node-diag-tool/pkg/util"
	log "github.com/sirupsen/logrus"
)

type ReuseportCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &ReuseportCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *ReuseportCheck) RunCheck() *plugin.CheckPluginResult {
	result, err := util.ExecCmd("ss -nlt|grep -i -E -v 'grep|Local'|awk '{print $4}'|uniq -c|sort -nr -k 1|head -n 10")
	if err != nil {
		log.Debugf("reuseport check failed, err: %v", err)
		return c.BuildPassResult(c, consts.ScriptFailed, "")
	}
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			row := strings.Split(line, " ")
			count, _ := strconv.Atoi(row[0])
			if count > 1 {
				return c.BuildBlockResult(c,
					fmt.Sprintf("存在两条LISTEN[%s]链接, 怀疑reuseport问题", row[1]),
					"")
			}
		}
	}
	return c.BuildPassResult(c, "检测reuseport正常(未发现reuseport)", "")
}

func (c *ReuseportCheck) PluginName() string {
	return "reuseport_check"
}

func (c *ReuseportCheck) Author() string {
	return "heroyf"
}
