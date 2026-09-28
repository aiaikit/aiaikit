package cli

import (
	"strings"
	"testing"
)

// reservedCommands 列出五个预留命令名；与 RegisterReserved 保持一致。
var reservedCommands = []string{"init", "update", "runtime", "plugin", "web"}

// TestReservedCommands_All 验证五个预留命令均已注册到 root。
func TestReservedCommands_All(t *testing.T) {
	root := NewRootCmd()

	for _, name := range reservedCommands {
		found := false
		for _, sub := range root.Commands() {
			if sub.Name() == name {
				found = true
				if v, ok := sub.Annotations[ReservedAnnotationKey]; !ok || v != ReservedAnnotationValue {
					t.Errorf("命令 %s 缺少 reserved=true 标注", name)
				}
				break
			}
		}
		if !found {
			t.Errorf("root 缺少预留命令 %s", name)
		}
	}
}

// TestReservedCommands_Run 验证每个预留命令输出友好提示且退出码符合 spec FR-004。
func TestReservedCommands_Run(t *testing.T) {
	for _, name := range reservedCommands {
		t.Run(name, func(t *testing.T) {
			stdout, stderr, code, err := runCmd(t, name)
			if code != ExitUsage {
				t.Errorf("退出码 = %d，期望 %d", code, ExitUsage)
			}
			if err == nil {
				t.Fatalf("期望非 nil 错误，实际 nil")
			}
			if !strings.Contains(err.Error(), "命令尚未实现") {
				t.Errorf("错误文案不符合预期: %v", err)
			}
			if !strings.Contains(stdout, `命令 "`+name+`" 尚未实现`) {
				t.Errorf("stdout 缺少命令名友好提示: %q", stdout)
			}
			if !strings.Contains(stdout, "Spec-Kit Runtime") {
				t.Errorf("stdout 缺少 Spec-Kit 提示: %q", stdout)
			}
			if !strings.Contains(stderr, "错误：") {
				t.Errorf("stderr 缺少中文错误模板: %q", stderr)
			}
		})
	}
}
