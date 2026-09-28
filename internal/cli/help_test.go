package cli

import (
	"bytes"
	"testing"
)

// TestHelp_Golden 比对 `aiaikit --help` 与 golden 文件。
func TestHelp_Golden(t *testing.T) {
	stdout, _, _, err := runCmd(t, "--help")
	if err != nil {
		t.Fatalf("执行 --help 失败: %v", err)
	}
	assertGolden(t, "help", stdout)
}

// TestHelp_NoArgs 比对无子命令调用与 golden 文件。
func TestHelp_NoArgs(t *testing.T) {
	stdout, _, _, err := runCmd(t)
	if err != nil {
		t.Fatalf("执行 无参 失败: %v", err)
	}
	assertGolden(t, "help", stdout)
}

// TestHelpSubcommand 验证 `aiaikit help init` 输出 init 子命令帮助。
func TestHelpSubcommand(t *testing.T) {
	stdout, _, _, err := runCmd(t, "help", "init")
	if err != nil {
		t.Fatalf("执行 help init 失败: %v", err)
	}
	if !bytes.Contains([]byte(stdout), []byte("初始化 aiaikit 工作区")) {
		t.Errorf("help init 输出缺少中文描述: %s", stdout)
	}
	if !bytes.Contains([]byte(stdout), []byte("Usage:")) {
		t.Errorf("help init 输出缺少 Usage 段: %s", stdout)
	}
}

// TestHelpExitCode 断言 --help 退出码为 0。
func TestHelpExitCode(t *testing.T) {
	_, _, code, err := runCmd(t, "--help")
	if err != nil || code != ExitSuccess {
		t.Fatalf("--help 退出码 = %d, err = %v, 期望 %d", code, err, ExitSuccess)
	}
}
