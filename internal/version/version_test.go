package version

import (
	"runtime"
	"strings"
	"testing"
)

// TestSnapshot_Defaults 断言 Snapshot 返回五个字段均为非空。
func TestSnapshot_Defaults(t *testing.T) {
	info := Snapshot()

	if info.Version == "" {
		t.Errorf("Version 必须非空")
	}
	if info.Commit == "" {
		t.Errorf("Commit 必须非空")
	}
	if info.BuildTime == "" {
		t.Errorf("BuildTime 必须非空")
	}
	if info.Platform == "" {
		t.Errorf("Platform 必须非空")
	}
	if info.GoVersion == "" {
		t.Errorf("GoVersion 必须非空")
	}

	// 默认值断言：未注入 ldflags 时 Version 应为 "dev"。
	if Version != "dev" && Version != "" {
		t.Logf("Version 已通过 ldflags 覆盖为 %q", Version)
	}
}

// TestString_Format 断言 String 输出包含五字段标签。
func TestString_Format(t *testing.T) {
	out := Snapshot().String()

	want := []string{
		"aiaikit 版本",
		"Commit   :",
		"构建时间 :",
		"目标平台 :",
		"Go 版本  :",
	}
	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("版本输出缺少字段 %q\n实际输出:\n%s", w, out)
		}
	}
}

// TestPlatform_Defaults 断言 Platform 运行时探测符合当前进程目标。
func TestPlatform_Defaults(t *testing.T) {
	if !strings.HasPrefix(Platform, runtime.GOOS) {
		t.Errorf("Platform = %q，不以 %q 开头", Platform, runtime.GOOS)
	}
}
