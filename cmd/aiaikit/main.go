// aiaikit CLI 入口。
//
// 本入口 MUST 仅完成三件事：
//  1. 调用 internal/cli.Execute()；
//  2. 按返回错误码退出；
//  3. 不向 stdout / stderr 写入任何字节（避免污染 --version / --help 快照）。
//
// 章程 I / V / XII 要求 Core 入口保持稳定，Runtime / Plugin / Web 业务一律通过
// internal/cli 暴露的扩展点接入，不在此处实现。
package main

import (
	"os"

	"github.com/aiaikit/aiaikit/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		// Execute() 已完成中文错误模板渲染与退出码转换；此处仅返回非零状态。
		os.Exit(cli.FromError(err))
	}
}
