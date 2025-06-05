package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/heroyf/node-diag-tool/pkg/consts"
	"github.com/heroyf/node-diag-tool/pkg/information"
	"github.com/heroyf/node-diag-tool/pkg/register"
	"github.com/heroyf/node-diag-tool/pkg/util"
	"github.com/heroyf/node-diag-tool/version"
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

	// 验证配置
	// if err := cfg.Validate(); err != nil {
	// 	return nil, err
	// }

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

	printlnInfo(fmt.Sprint(strings.Repeat("#", 40), "机器基础属性", strings.Repeat("#", 40)))
	printfInfo("诊断IP:\t\t\t[%s]\n", cfg.CurrentIP)
	printfInfo("发行版本:\t\t[%s]\n", information.ReleaseVersion())
	printfInfo("内核版本:\t\t[%s]\n", information.KernelVersion())
	printfInfo("OS运行时间:\t\t[%s]\n", information.Uptime())
	printfInfo("CPU架构:\t\t[%s]\n", information.CpuArch())
	printfInfo("CPU核心:\t\t[%d]\n", information.CpuCores())
	printfInfo("内存大小/使用率:\t[%dGi]/[%s]\n", information.Memory(), information.MemoryPercent())
	printfInfo("透明大页:\t\t[%sKB]\n", information.HugePage())
	printfInfo("进程个数:\t\t[%s]\n", information.PidNum())
	printfInfo("Swap内存(used/total):\t[%s/%s] \n", information.SwapUsed(), information.SwapTotal())
	printlnInfo(strings.Repeat("#", 94))
	printfInfo("总检测项:\t\t[%d]\n", len(register.PluginRegisters))
}

func printfInfo(format string, args ...any) {
	fmt.Fprintf(diagLog, format, args...)
	fmt.Fprintf(os.Stdout, format, args...)
}

func printlnInfo(content string) {
	fmt.Fprintln(diagLog, content)
	fmt.Fprintln(os.Stdout, content)
}

func timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("time cost = %v\n", tc)
	}
}
