package cli

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// updateGoldenFlag 控制是否在测试期间更新 golden 文件。
// 通过环境变量 UPDATE_GOLDEN=1 启用；默认 false。
var updateGoldenFlag = flag.Bool("update-golden", false, "更新 golden 文件")

// isUpdateGoldenEnabled 检查环境变量；用于非测试入口。
func isUpdateGoldenEnabled() bool {
	if v := os.Getenv("UPDATE_GOLDEN"); v != "" && v != "0" && v != "false" {
		return true
	}
	if updateGoldenFlag != nil && *updateGoldenFlag {
		return true
	}
	return false
}

// assertGolden 比对 actual 与 testdata/cli/<name>.golden 的内容。
//
// 当 UPDATE_GOLDEN=1（或 -update-golden 标志）时，写入新快照而非断言失败。
//
// 返回值：
//   - 当 actual 与 golden 内容一致时，PASS；
//   - 当启用更新时，覆盖 golden 并 PASS；
//   - 否则，调用 t.Errorf 报告差异并继续执行（不 FailNow）。
func assertGolden(t *testing.T, name string, actual string) {
	t.Helper()

	goldenPath := filepath.Join("..", "..", "testdata", "cli", name+".golden")
	if isUpdateGoldenEnabled() {
		// 写盘：包含结尾换行。
		if !strings.HasSuffix(actual, "\n") {
			actual += "\n"
		}
		if err := os.WriteFile(goldenPath, []byte(actual), 0o644); err != nil {
			t.Fatalf("更新 golden 失败: %v", err)
		}
		t.Logf("已更新 golden: %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("读取 golden 失败 %s: %v", goldenPath, err)
	}

	// 去除 actual 末尾多余换行（来自 fmt.Print）。
	got := strings.TrimRight(actual, "\n")
	wanted := strings.TrimRight(string(want), "\n")

	if got != wanted {
		t.Errorf("golden 不匹配: %s\n--- got ---\n%s\n--- want ---\n%s", name, got, wanted)
	}
}
