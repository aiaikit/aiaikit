package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// binaryPath 返回 bin/aiaikit 的绝对路径。
// 测试假定已通过 make build 在仓库根目录构建出二进制。
func binaryPath(t *testing.T) string {
	t.Helper()
	root, err := repoRoot()
	if err != nil {
		t.Fatalf("查找仓库根目录失败: %v", err)
	}
	candidates := []string{
		filepath.Join(root, "bin", "aiaikit"),
		filepath.Join(root, "bin", "aiaikit.exe"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	t.Skipf("未找到构建产物，请先执行 make build；尝试路径: %v", candidates)
	return ""
}

// repoRoot 通过 git 命令获取仓库根目录。
func repoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// TestVersionCommand 通过真实二进制调用 --version 并断言字段集合。
func TestVersionCommand(t *testing.T) {
	bin := binaryPath(t)

	cmd := exec.Command(bin, "--version")
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		t.Fatalf("执行 %s --version 失败: %v, stderr=%s", bin, err, errBuf.String())
	}

	got := out.String()
	want := []string{
		"aiaikit 版本",
		"Commit   :",
		"构建时间 :",
		"目标平台 :",
		"Go 版本  :",
	}
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("--version 输出缺少字段 %q\n实际输出:\n%s", w, got)
		}
	}

	if errBuf.Len() != 0 {
		t.Errorf("--version 不应向 stderr 写入任何内容，实际: %s", errBuf.String())
	}
}

// TestVersionExitCode 断言 --version 退出码为 0。
func TestVersionExitCode(t *testing.T) {
	bin := binaryPath(t)
	cmd := exec.Command(bin, "--version")
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("期望退出码 0，实际 %d", ee.ExitCode())
		}
		t.Fatalf("执行失败: %v", err)
	}
}
