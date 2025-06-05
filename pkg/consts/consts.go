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
