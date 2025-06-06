package plugin

import (
	"embed"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/heroyf/node-diag-tool/pkg/consts"
	log "github.com/sirupsen/logrus"
)

type BaseCheckPlugin interface {
	RunCheck() *CheckPluginResult
	PluginName() string
	Author() string
}

type CheckPluginResult struct {
	// 持有插件对象
	BaseCheckPlugin
	CheckState  consts.DiagnosisState
	CheckResult string
	// 插件执行完整的标准输出
	Stdout string
}

type Helper struct {
	FileName string
}

//go:embed healthcheck/general/*/script/*.sh
var scripts embed.FS

func (h *Helper) LookupMatchedScript() string {
	log.Debugf("dynamic lookup script use embed!!!")
	// 先获取base.go的完整路径
	_, fullFilename, _, _ := runtime.Caller(0)

	// 获取插件相对路径(例如: healthcheck/general/cpu/cpu_load_check.go)
	pluginRelativePath, _ := filepath.Rel(filepath.Dir(fullFilename), h.FileName)

	// 拼接shell脚本父路径
	scriptPath := filepath.Join(filepath.Dir(pluginRelativePath), consts.ScriptDir)

	// 拼接最后的shell脚本全路径(例如: general/net/script/eth_cpuaffinity_check.sh)
	scriptName := filepath.Join(scriptPath, filepath.Base(h.FileName))
	scriptName = strings.Replace(scriptName, ".go", ".sh", -1)

	log.Debugf("pluginRelativePath: [%s], scriptPath: [%s], scriptName: [%s] \n", pluginRelativePath,
		scriptPath, scriptName)

	// 使用embed读取ELF内嵌的脚本内容
	content, err := scripts.ReadFile(scriptName)
	if err != nil {
		log.Debugf("script not found: %s", scriptName)
		return ""
	}
	return string(content)
}

// ReadLines 将shell输出转换成多行
func (h *Helper) ReadLines(result string) []string {
	return strings.Split(result, "\n")
}

func (h *Helper) BuildPluginName(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), ".go")
}

func (h *Helper) BuildPassResult(c BaseCheckPlugin,
	checkResult, stdout string) *CheckPluginResult {
	return h.buildCheckResult(c, consts.Pass, checkResult, stdout)
}

func (h *Helper) BuildBlockResult(c BaseCheckPlugin,
	checkResult, stdout string) *CheckPluginResult {
	return h.buildCheckResult(c, consts.Blocked, checkResult, stdout)
}

func (h *Helper) buildCheckResult(c BaseCheckPlugin,
	state consts.DiagnosisState,
	checkResult, stdout string) *CheckPluginResult {
	return &CheckPluginResult{
		BaseCheckPlugin: c,
		CheckState:      state,
		CheckResult:     checkResult,
		Stdout:          stdout,
	}
}

type ShellResult struct {
	Result    string `json:"result"`
	ResultMsg string `json:"result_msg"`
}
