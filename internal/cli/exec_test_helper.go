package cli

import (
	"bytes"
	"testing"
)

// runCmd 在测试中执行 root 命令的便捷封装：
//   - 注入 stdout / stderr 缓冲；
//   - 自动安装 localized HelpFunc；
//   - 返回 exit code 而非调用 os.Exit。
//
// 返回值：
//   - stdout 捕获的全部字节；
//   - stderr 捕获的全部字节；
//   - 退出码（由 FromError(err) 计算）；
//   - root.Execute 返回的错误（可能为 nil）。
func runCmd(t *testing.T, args ...string) (stdout string, stderr string, exitCode int, err error) {
	t.Helper()

	var outBuf, errBuf bytes.Buffer

	// 构造 root 并应用本地化 HelpFunc。
	root := NewRootCmd()
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	installLocalizedHelpFunc(root)

	root.SetArgs(args)
	err = root.Execute()
	if err != nil {
		writeErrorTemplate(&errBuf, err.Error())
	}
	exitCode = FromError(err)
	return outBuf.String(), errBuf.String(), exitCode, err
}
