package information

var kernelInfo = NewKernelInfo()

// ReleaseVersion 获取发行版本
func ReleaseVersion() string {
	return kernelInfo.ReleaseVersion()
}

// KernelVersion 获取内核版本
func KernelVersion() string {
	return kernelInfo.KernelVersion()
}

// Uptime 获取系统运行时间
func Uptime() string {
	return kernelInfo.Uptime()
}

// CpuArch 获取CPU架构
func CpuArch() string {
	return kernelInfo.CpuArch()
}

// CpuCores 获取CPU核心数
func CpuCores() int {
	return kernelInfo.CpuCores()
}

// CpuLoad 获取CPU负载
func CpuLoad() float32 {
	return kernelInfo.CpuLoad()
}

// CpuPercent 获取CPU使用率
func CpuPercent() string {
	return kernelInfo.CpuPercent()
}

// PidNum 获取进程数
func PidNum() string {
	return kernelInfo.PidNum()
}

// Memory 获取内存大小（单位：Gi）
func Memory() int {
	return kernelInfo.Memory()
}

// MemoryPercent 获取内存使用率
func MemoryPercent() string {
	return kernelInfo.MemoryPercent()
}

// SwapTotal 获取Swap总量
func SwapTotal() string {
	return kernelInfo.SwapTotal()
}

// SwapUsed 获取Swap使用量
func SwapUsed() string {
	return kernelInfo.SwapUsed()
}

// HugePage 获取透明大页信息
func HugePage() string {
	return kernelInfo.HugePage()
}
