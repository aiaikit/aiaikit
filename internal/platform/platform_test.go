package platform

import (
	"runtime"
	"testing"
)

// TestCurrent_ReturnsRuntimeValues 断言 Current 返回值与 runtime 包一致。
func TestCurrent_ReturnsRuntimeValues(t *testing.T) {
	info := Current()

	if info.OS != runtime.GOOS {
		t.Errorf("OS = %q, expect %q", info.OS, runtime.GOOS)
	}
	if info.Arch != runtime.GOARCH {
		t.Errorf("Arch = %q, expect %q", info.Arch, runtime.GOARCH)
	}

	// IsUnixLike / IsWindows 互斥。
	if info.IsUnixLike == info.IsWindows {
		t.Errorf("IsUnixLike (%v) 与 IsWindows (%v) 应互斥", info.IsUnixLike, info.IsWindows)
	}

	if runtime.GOOS == "windows" {
		if !info.IsWindows {
			t.Errorf("runtime.GOOS=windows 时 IsWindows 必须为 true")
		}
		if info.IsUnixLike {
			t.Errorf("runtime.GOOS=windows 时 IsUnixLike 必须为 false")
		}
	} else {
		if info.IsWindows {
			t.Errorf("runtime.GOOS=%s 时 IsWindows 必须为 false", runtime.GOOS)
		}
		// 非 Windows 时，IsUnixLike 在 unix 编译标签下为 true，否则未必。
		// 接受任意非 Windows 下的非 true 情形，但记录以便观察。
		if !info.IsUnixLike {
			t.Logf("runtime.GOOS=%s 时 IsUnixLike=false（不常见，仅记录）", runtime.GOOS)
		}
	}
}
