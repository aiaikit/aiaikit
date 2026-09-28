package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestReservedCommands_NoDiskWrites 验证预留命令不会修改磁盘。
//
// 方法：在临时目录中执行预留命令，校验前后目录快照（仅一级）
// 与初始状态完全一致。允许目录本身被创建（go test 默认行为），
// 但不允许创建任何新文件。
func TestReservedCommands_NoDiskWrites(t *testing.T) {
	tmpDir := t.TempDir()

	// 切换到临时目录，确保即使有磁盘写入也会发生在 tmpDir 内。
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("切换到临时目录失败: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	before, err := listDir(tmpDir)
	if err != nil {
		t.Fatalf("列出临时目录失败: %v", err)
	}

	for _, name := range reservedCommands {
		name := name
		t.Run(name, func(t *testing.T) {
			cmd := NewRootCmd()
			cmd.SetArgs([]string{name})

			var out, errBuf bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&errBuf)

			_ = cmd.Execute()

			// 仅校验当前目录下的文件集合不变（深度为 1）。
			after, err := listDir(tmpDir)
			if err != nil {
				t.Fatalf("列出临时目录失败: %v", err)
			}
			if !equalFileSet(before, after) {
				t.Errorf("预留命令 %s 修改了磁盘文件集合\nbefore=%v\nafter=%v",
					name, before, after)
			}
		})
	}
}

// listDir 列出一级目录下的所有条目（文件 + 子目录）。
func listDir(dir string) (map[string]bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(entries))
	for _, e := range entries {
		out[filepath.Join(dir, e.Name())] = true
	}
	return out, nil
}

// equalFileSet 比较两个文件集合是否一致。
func equalFileSet(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}
