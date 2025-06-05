package plugin

import "github.com/heroyf/node-diag-tool/pkg/consts"

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
