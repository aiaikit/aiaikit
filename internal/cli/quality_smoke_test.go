package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// repoRootFromCwd 通过查找 Makefile 确定仓库根目录。
func repoRootFromCwd(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "Makefile")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("未找到仓库根目录（向上查找无 Makefile）")
		}
		dir = parent
	}
}

// TestMakeCi 验证 `make ci` 的核心目标（fmt + test + build）在本地可通过，
// 作为质量门禁的最后一道闸。
//
// 通过环境变量 AIK_TEST_SKIP_MAKE_CI=1 避免递归：
// TestMakeCi 本身会触发 `make test`，而 test 目标中再次运行 go test
// 时若仍执行 TestMakeCi 会无限递归。
func TestMakeCi(t *testing.T) {
	if os.Getenv("AIK_TEST_SKIP_MAKE_CI") != "" {
		t.Skip("AIK_TEST_SKIP_MAKE_CI 已设置，跳过以避免递归")
	}
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make 未安装，跳过此测试")
	}
	repo := repoRootFromCwd(t)
	cmd := exec.Command("make", "ci")
	cmd.Dir = repo
	cmd.Env = append(os.Environ(), "AIK_TEST_SKIP_MAKE_CI=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("make ci 失败: %v\n%s", err, out)
	}
}
