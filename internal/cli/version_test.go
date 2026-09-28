package cli

import (
	"testing"

	"github.com/aiaikit/aiaikit/internal/version"
)

// TestVersionCommand_Golden 比对 `aiaikit version` 与 golden 文件。
func TestVersionCommand_Golden(t *testing.T) {
	stdout, _, _, err := runCmd(t, "version")
	if err != nil {
		t.Fatalf("执行 version 命令失败: %v", err)
	}
	assertGolden(t, "version", stdout)
}

// TestVersionFlag_Golden 比对 `aiaikit --version` 与 golden 文件。
func TestVersionFlag_Golden(t *testing.T) {
	stdout, _, _, err := runCmd(t, "--version")
	if err != nil {
		t.Fatalf("执行 --version 失败: %v", err)
	}
	assertGolden(t, "version", stdout)
}

// TestVersionShortFlag_Golden 比对 `aiaikit -v` 与 golden 文件。
func TestVersionShortFlag_Golden(t *testing.T) {
	stdout, _, _, err := runCmd(t, "-v")
	if err != nil {
		t.Fatalf("执行 -v 失败: %v", err)
	}
	assertGolden(t, "version", stdout)
}

// TestVersionSnapshot_Defaults 断言默认版本变量符合 data-model.md 不变式。
func TestVersionSnapshot_Defaults(t *testing.T) {
	info := version.Snapshot()

	if info.Version == "" {
		t.Errorf("Version MUST NOT 为空字符串")
	}
	if info.Commit == "" {
		t.Errorf("Commit MUST NOT 为空字符串")
	}
	if info.BuildTime == "" {
		t.Errorf("BuildTime MUST NOT 为空字符串")
	}
	if info.Platform == "" {
		t.Errorf("Platform MUST NOT 为空字符串")
	}
	if info.GoVersion == "" {
		t.Errorf("GoVersion MUST NOT 为空字符串")
	}
}

// TestVersionExitCode 断言 --version 退出码为 0。
func TestVersionExitCode(t *testing.T) {
	_, _, code, err := runCmd(t, "--version")
	if err != nil || code != ExitSuccess {
		t.Fatalf("--version 退出码 = %d, err = %v, 期望 %d", code, err, ExitSuccess)
	}
}
