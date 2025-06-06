package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/heroyf/node-diag-tool/pkg/consts"
	"github.com/heroyf/node-diag-tool/pkg/plugin"
	log "github.com/sirupsen/logrus"
)

var diagLog *os.File

// SetDiagLog sets the diagnostic log file
func SetDiagLog(logFile *os.File) {
	diagLog = logFile
}

// PrintfInfo prints formatted info to both log file and stdout
func PrintfInfo(format string, args ...any) {
	fmt.Fprintf(diagLog, format, args...)
	fmt.Fprintf(os.Stdout, format, args...)
}

// PrintlnInfo prints a line to both log file and stdout
func PrintlnInfo(content string) {
	fmt.Fprintln(diagLog, content)
	fmt.Fprintln(os.Stdout, content)
}

// MaxLenWhiteSpace calculates the maximum length for whitespace
func MaxLenWhiteSpace(maxNameLen int, maxStateLen int) int {
	return maxNameLen + maxStateLen + 80
}

// RenderResult renders the check results
func RenderResult(checkResults []*plugin.CheckPluginResult, onlyBlock bool, verbose bool) {
	// 计算最大长度
	maxNameLen := 0
	maxStateLen := 0
	for _, result := range checkResults {
		if result == nil {
			log.Errorf("check result is nil")
			continue
		}
		if len(result.PluginName()) > maxNameLen {
			maxNameLen = len(result.PluginName())
		}

		if len(result.CheckState) > maxStateLen {
			maxStateLen = len(result.CheckState)
		}
	}

	// 计算总长度
	totalLen := MaxLenWhiteSpace(maxNameLen, maxStateLen)
	PrintlnInfo(strings.Repeat("#", totalLen))

	// 格式化输出
	for _, result := range checkResults {
		if result == nil {
			log.Errorf("check result is nil, skip plugin.")
			continue
		}

		// 跳过检测正常的插件
		if onlyBlock && result.CheckState != consts.Blocked {
			continue
		}

		name := result.PluginName()
		state := strings.ToUpper(string(result.CheckState))
		namePadding := strings.Repeat(" ", maxNameLen-len(name))
		statePadding := strings.Repeat(" ", maxStateLen-len(state))
		PrintfInfo("检测项: [%s]%s    |    判定状态: [%s]%s    |    判定结果: [%s]\n\n",
			name, namePadding,
			state, statePadding,
			result.CheckResult)

		// 打印标准输出
		if result.Stdout != "" && verbose {
			PrintlnInfo("=================诊断详细说明 Begin()=================")
			PrintlnInfo(result.Stdout)
			PrintlnInfo("=================诊断详细说明 End()===================")
			PrintfInfo("\n")
		}
	}
	PrintlnInfo(strings.Repeat("#", totalLen))
}
