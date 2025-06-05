package kernel

// BaseMachineInfo 定义机器基础信息接口
type BaseMachineInfo interface {
	ReleaseVersion() string
	KernelVersion() string
	Uptime() string
	CpuArch() string
	CpuCores() int
	CpuLoad() float32
	CpuPercent() string
	PidNum() string
	Memory() int
	MemoryPercent() string
	SwapTotal() string
	SwapUsed() string
	HugePage() string
}
