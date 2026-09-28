// Package platform 提供运行时平台信息抽象。
//
// 本阶段仅用于 --version 输出与未来跨平台路径/进程适配占位，
// 不引入任何平台特定的行为分支。
package platform

import "runtime"

// Info 描述当前运行平台的元数据。
type Info struct {
	OS         string // runtime.GOOS 取值，例如 "windows" / "linux" / "darwin"
	Arch       string // runtime.GOARCH 取值，例如 "amd64" / "arm64"
	IsUnixLike bool   // Unix 类系统判定
	IsWindows  bool   // Windows 判定
}

// Current 返回当前进程的平台信息。
func Current() Info {
	return Info{
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		IsUnixLike: isUnixLike(),
		IsWindows:  runtime.GOOS == "windows",
	}
}
