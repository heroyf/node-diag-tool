package register

import "github.com/heroyf/node-diag-tool/pkg/plugin"

// PluginRegisters 定义全局插件注册器
var PluginRegisters = make(map[string]plugin.BaseCheckPlugin)

// EnabledPlugins 启用的插件(默认是启用所有)
var EnabledPlugins = []string{}

// DiabledPlugins 禁用的插件
var DiabledPlugins = []string{}
