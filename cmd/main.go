package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/heroyf/node-diag-tool/pkg/consts"
	_ "github.com/heroyf/node-diag-tool/pkg/imports"
	"github.com/heroyf/node-diag-tool/pkg/information"
	"github.com/heroyf/node-diag-tool/pkg/plugin"
	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/heroyf/node-diag-tool/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Config 应用配置结构体
type Config struct {
	// 基础配置
	Help      bool   // 帮助信息
	Version   bool   // 版本信息
	Verbose   bool   // 详细日志模式
	Debug     bool   // 调试日志级别
	CurrentIP string // 当前IP地址

	// 日志配置
	DiagLogFile string // 诊断日志文件名

	// 文件配置
	DefaultConfigPath string // 默认配置文件路径
	DefaultConfigName string // 默认配置文件名

	// 功能开关
	OnlyBlock bool // 仅显示阻塞检查项
}

func parseFlags() (*Config, error) {
	cfg := &Config{}

	// 获取当前IP
	var err error
	cfg.CurrentIP, err = util.GetLocalIp()
	if err != nil {
		return nil, fmt.Errorf("获取本地IP失败: %w", err)
	}

	// 基础配置
	pflag.BoolVarP(&cfg.Version, "version", "V", false, "显示版本信息")
	pflag.BoolVarP(&cfg.Verbose, "verbose", "v", false, "详细日志模式")
	pflag.BoolVarP(&cfg.Debug, "debug", "d", false, "调试日志级别")

	// 日志配置
	pflag.StringVar(&cfg.DiagLogFile, "diagFileName", "diag.log", "诊断日志文件名")
	pflag.CommandLine.MarkHidden("diagFileName")

	// 文件配置
	pflag.StringVar(&cfg.DefaultConfigPath, "defaultPath", "./", "默认配置文件路径")
	pflag.StringVar(&cfg.DefaultConfigName, "defaultConfig", "application", "默认配置文件名")

	// 功能开关
	pflag.BoolVar(&cfg.OnlyBlock, "onlyBlock", false, "仅显示阻塞检查项")

	pflag.Parse()

	return cfg, nil
}

func main() {
	// 设置内存限制
	debug.SetMemoryLimit(consts.MaxMemoryLimit)

	// 解析命令行参数
	cfg, err := parseFlags()
	if err != nil {
		fmt.Printf("配置错误: %v\n", err)
		return
	}

	// 处理版本信息
	if cfg.Version {
		fmt.Println(version.Version())
		return
	}

	// 处理帮助信息
	if cfg.Help {
		pflag.Usage()
		return
	}

	cfg.runDiag()
}

var (
	diagLog *os.File
)

func (cfg *Config) runDiag() {
	defer timeCost()
	os.RemoveAll(cfg.DiagLogFile)
	var err error
	diagLog, err = os.OpenFile(cfg.DiagLogFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("create diag log file failed", err)
	}
	util.SetDiagLog(diagLog)

	util.PrintlnInfo(fmt.Sprint(strings.Repeat("#", 40), "机器基础属性", strings.Repeat("#", 40)))
	util.PrintfInfo("诊断IP:\t\t\t[%s]\n", cfg.CurrentIP)
	util.PrintfInfo("发行版本:\t\t[%s]\n", information.ReleaseVersion())
	util.PrintfInfo("内核版本:\t\t[%s]\n", information.KernelVersion())
	util.PrintfInfo("OS运行时间:\t\t[%s]\n", information.Uptime())
	util.PrintfInfo("CPU架构:\t\t[%s]\n", information.CpuArch())
	util.PrintfInfo("CPU核心:\t\t[%d]\n", information.CpuCores())
	util.PrintfInfo("内存大小/使用率:\t[%dGi]/[%s]\n", information.Memory(), information.MemoryPercent())
	util.PrintfInfo("透明大页:\t\t[%sKB]\n", information.HugePage())

	util.PrintlnInfo(strings.Repeat("#", 94))
	util.PrintfInfo("总检测项:\t\t[%d]\n", len(plugin.PluginRegisters))

	plugin.DynamicLoadPlugins(cfg.DefaultConfigPath, cfg.DefaultConfigName)

	var mutex sync.Mutex
	var wg sync.WaitGroup
	checkResults := make([]*plugin.CheckPluginResult, 0, len(plugin.PluginRegisters))
	for _, p := range plugin.PluginRegisters {
		if len(plugin.EnabledPlugins) > 0 {
			// 如果不包含启用的插件, 则跳过
			if !util.Contains(plugin.EnabledPlugins, p.PluginName()) {
				log.Debugf("skip enabled plugin: %s", p.PluginName())
				continue
			}
		}

		// 跳过禁用的插件
		if util.Contains(plugin.DiabledPlugins, p.PluginName()) {
			log.Debugf("skip disabled plugin: %s", p.PluginName())
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			result := p.RunCheck()

			mutex.Lock()
			checkResults = append(checkResults, result)
			mutex.Unlock()
		}()
	}

	wg.Wait()
	util.RenderResult(checkResults, cfg.OnlyBlock, cfg.Verbose)
}

func timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("time cost = %v\n", tc)
	}
}
