// Package version 集中声明 aiaikit 自身的版本元数据。
//
// 这些变量在链接期通过 -ldflags "-X" 注入；缺失注入参数时使用默认值，
// 保证开发模式（go run）与缺失构建信息时不崩溃。
package version

import "runtime"

// 默认值常量；任何 ldflags 注入都必须覆盖到这些默认。
const (
	defaultVersion   = "dev"
	defaultCommit    = "unknown"
	defaultBuildTime = "unknown"
)

// 版本信息变量。运行时由 ldflags 注入覆盖。
//
// 注入路径：
//
//	-X github.com/aiaikit/aiaikit/internal/version.Version=...
//	-X github.com/aiaikit/aiaikit/internal/version.Commit=...
//	-X github.com/aiaikit/aiaikit/internal/version.BuildTime=...
//	-X github.com/aiaikit/aiaikit/internal/version.Platform=...
//	-X github.com/aiaikit/aiaikit/internal/version.GoVersion=...
var (
	// Version 语义化版本号，例如 "0.1.0"。
	Version = defaultVersion

	// Commit 构建时的 Git Commit Hash。
	Commit = defaultCommit

	// BuildTime 构建时间，RFC3339 格式，例如 "2026-09-28T10:15:30Z"。
	BuildTime = defaultBuildTime

	// Platform 目标平台，例如 "linux/amd64"。运行时填入默认值，ldflags 可覆盖。
	Platform = runtime.GOOS + "/" + runtime.GOARCH

	// GoVersion 编译期 Go 版本，例如 "go1.23.0"。运行时填入，ldflags 可覆盖。
	GoVersion = runtime.Version()
)

// Info 描述一次构建的完整版本元数据快照。
type Info struct {
	Version   string
	Commit    string
	BuildTime string
	Platform  string
	GoVersion string
}

// Snapshot 返回当前全局版本变量的快照。
func Snapshot() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
		Platform:  Platform,
		GoVersion: GoVersion,
	}
}

// String 渲染版本信息为 contracts/cli.md §1 规定的五字段输出格式。
//
// 字段顺序固定：
//  1. Version
//  2. Commit
//  3. BuildTime
//  4. Platform
//  5. GoVersion
//
// 使用纯字符串拼接以避免引入 text/template 与上游包产生依赖循环。
func (i Info) String() string {
	// 每行单独构造，便于未来字段变更最小化测试影响。
	return "aiaikit 版本 " + i.Version + "\n" +
		"  Commit   : " + i.Commit + "\n" +
		"  构建时间 : " + i.BuildTime + "\n" +
		"  目标平台 : " + i.Platform + "\n" +
		"  Go 版本  : " + i.GoVersion + "\n"
}
