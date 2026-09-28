package e2e

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

// TestHelpCommand 通过真实二进制调用 --help 并验证核心字段。
func TestHelpCommand(t *testing.T) {
	bin := binaryPath(t)
	cmd := exec.Command(bin, "--help")
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		t.Fatalf("执行 %s --help 失败: %v, stderr=%s", bin, err, errBuf.String())
	}

	got := out.String()
	// 八条 Available Commands 全部必须出现。
	want := []string{
		"completion",
		"help",
		"init",
		"plugin",
		"runtime",
		"update",
		"version",
		"web",
	}
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("--help 输出缺少命令 %q", w)
		}
	}

	// 五个预留命令描述应包含“（预留 / 未实现）”。
	reserved := []string{"init", "plugin", "runtime", "update", "web"}
	for _, r := range reserved {
		// 在原文中找到以命令名开头的行，确认后缀。
		if !strings.Contains(got, r+" ") {
			t.Errorf("--help 输出缺少命令 %s 行", r)
			continue
		}
		// 行级断言：以命令名为锚点开始，到下一个换行为止。
		idx := strings.Index(got, r+" ")
		end := strings.Index(got[idx:], "\n")
		if end < 0 {
			end = len(got) - idx
		}
		line := got[idx : idx+end]
		if !strings.Contains(line, "（预留 / 未实现）") {
			t.Errorf("--help 中命令 %s 行缺少预留标记: %q", r, line)
		}
	}

	if errBuf.Len() != 0 {
		t.Errorf("--help 不应向 stderr 写入任何内容，实际: %s", errBuf.String())
	}
}

// TestHelpExitCode 断言 --help 退出码为 0。
func TestHelpExitCode(t *testing.T) {
	bin := binaryPath(t)
	cmd := exec.Command(bin, "--help")
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("期望退出码 0，实际 %d", ee.ExitCode())
		}
		t.Fatalf("执行失败: %v", err)
	}
}

// TestReservedCommand_ExitCode 验证预留命令退出码为 2。
func TestReservedCommand_ExitCode(t *testing.T) {
	bin := binaryPath(t)
	for _, name := range []string{"init", "update", "runtime", "plugin", "web"} {
		t.Run(name, func(t *testing.T) {
			cmd := exec.Command(bin, name)
			err := cmd.Run()
			if err == nil {
				t.Fatalf("期望非零退出码，实际为 0")
			}
			ee, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatalf("期望 ExitError，实际 %T", err)
			}
			if ee.ExitCode() != 2 {
				t.Errorf("退出码 = %d，期望 2", ee.ExitCode())
			}
		})
	}
}
