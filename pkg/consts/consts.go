package consts

type DiagnosisState string

const (
	Pass    DiagnosisState = "Pass"
	Blocked DiagnosisState = "Blocked"
	Unknown DiagnosisState = "Unknown"
)

const (
	// MaxMemoryLimit 最大内存使用(1Gi)
	MaxMemoryLimit = 1024 * 1 << 20
	Page           = 4
)

const (
	ScriptFailed = "执行脚本异常"
	ScriptDir    = "script"
)

const (
	MemMinFreeKBytes  = "/proc/sys/vm/min_free_kbytes"
	DefaultQdisc      = "/proc/sys/net/core/default_qdisc"
	TcpMem            = "/proc/sys/net/ipv4/tcp_mem"
	ConntrackCount    = "/proc/sys/net/netfilter/nf_conntrack_count"
	ConntrackMax      = "/proc/sys/net/netfilter/nf_conntrack_max"
	FsFileNr          = "/proc/sys/fs/file-nr"
	BugQdiscAlgorithm = "pfifo_fast"
)
