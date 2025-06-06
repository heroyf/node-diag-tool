// plugin/general/cpu/cpu_load_check.go
package cpu

import (
	"fmt"
	"runtime"

	"github.com/heroyf/node-diag-tool/pkg/information"
	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/sirupsen/logrus"
)

type CpuLoadCheck struct {
	plugin.Helper
}

func init() {
	_, fullFilename, _, _ := runtime.Caller(0)
	c := &CpuLoadCheck{
		plugin.Helper{
			FileName: fullFilename,
		},
	}
	plugin.PluginRegisters[c.PluginName()] = c
}

func (c *CpuLoadCheck) RunCheck() *plugin.CheckPluginResult {
	loadaverage := information.CpuLoad()
	cpuCores := information.CpuCores()
	logrus.Debugf("load average: %v \n", loadaverage)
	if loadaverage > float32(cpuCores) {
		return c.BuildBlockResult(c, "检测CPU负载异常",
			fmt.Sprintf("当前负载: [%.2f], cpu核数: [%d]", loadaverage, cpuCores))
	}
	return c.BuildPassResult(c, "检测CPU负载正常", fmt.Sprintf("当前负载: %.2f", loadaverage))
}

func (c *CpuLoadCheck) PluginName() string {
	return c.BuildPluginName(c.FileName)
}

func (c *CpuLoadCheck) Author() string {
	return "heroyf"
}
