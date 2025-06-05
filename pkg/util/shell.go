package util

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ShellResult struct {
	Result    string `json:"result"`
	ResultMsg string `json:"result_msg"`
}

// ExecCmd 执行shell
func ExecCmd(shell string, args ...string) (string, error) {
	// 兼容embed超长文本内容, 减少错误输出
	if shell != "" && len(shell) < 100 {
		log.Debugf("cmd: [%s]", shell)
	}
	cmd := exec.Command("/bin/sh", "-c", shell)
	// exit code: 0, err为nil
	// exit code: 1, err不为nil
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// ParseShellResult 解析shell执行结果
func ParseShellResult(result string) *ShellResult {
	re := regexp.MustCompile(`OUTPUT_RESULT\s+(.*?)\s+OUTPUT_END`)
	match := re.FindStringSubmatch(result)
	if len(match) > 1 {
		s := &ShellResult{}
		json.Unmarshal([]byte(match[1]), s)
		return s
	}
	return nil
}
