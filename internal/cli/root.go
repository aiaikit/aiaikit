package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aiaikit/aiaikit/internal/version"
)

// ExecuteOptions 控制命令执行的输出与退出行为。
type ExecuteOptions struct {
	Stdout io.Writer
	Stderr io.Writer
	Exit   func(int)
}

// NewRootCmd 构建并返回 root 命令。
//
// 本函数为后续 Plugin / Runtime / Web 阶段保留稳定的扩展点；
// 当前阶段实现 --version / --help 与五个预留命令占位。
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   AppName,
		Short: "aiaikit —— 基于 Go 的 AI/SDD 工具生态平台 CLI",
		Long:  AppLongDescription,
		// 关闭 Cobra 内置错误打印；改由 ExecuteWith 统一走 ErrorTemplate。
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// 全局 --version / -v 标志。
	root.Flags().BoolP("version", "v", false, "查看版本信息")

	// 注册 version 子命令。
	root.AddCommand(newVersionCmd())

	// 注册五个预留命令（init / update / runtime / plugin / web）。
	RegisterReserved(root)

	// 根命令的 RunE：响应 --version 标志；未识别时输出帮助文本。
	root.RunE = func(cmd *cobra.Command, args []string) error {
		if v, err := cmd.Flags().GetBool("version"); err == nil && v {
			return renderVersionTo(cmd.OutOrStdout())
		}
		return cmd.Help()
	}

	return root
}

// Execute 在本默认位置串起 root 命令、输出与退出码转换。
func Execute() error {
	return ExecuteWith(nil)
}

// ExecuteWith 在测试或嵌入场景下允许注入 stdout / stderr / exit。
// opts 为 nil 时使用默认值（stdout / stderr / os.Exit）。
//
// 流程：
//  1. 构造 root 命令；
//  2. 覆盖输出与 HelpFunc（本地化内置命令）；
//  3. 执行并按错误码退出。
func ExecuteWith(opts *ExecuteOptions) error {
	if opts == nil {
		opts = &ExecuteOptions{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Exit:   os.Exit,
		}
	}

	root := NewRootCmd()
	root.SetOut(opts.Stdout)
	root.SetErr(opts.Stderr)
	installLocalizedHelpFunc(root)

	err := root.Execute()
	if err != nil {
		// 渲染统一的中文错误模板（章程 XV）。
		writeErrorTemplate(opts.Stderr, err.Error())
	}
	opts.Exit(FromError(err))
	return err
}

// installLocalizedHelpFunc 安装覆盖 root 的 HelpFunc，
// 使 --help 渲染前先把 Cobra 内置命令的 Short 本地化为中文。
//
// 工作原理：
//  1. 保存一份 fallback 闭包，调用原 root.HelpFunc() 走 Cobra 默认渲染；
//  2. 替换 root 的 HelpFunc 为"先本地化 → 调用 fallback"；
//  3. fallback 通过 sync.Once 风格的标志位防止递归。
func installLocalizedHelpFunc(root *cobra.Command) {
	var entered bool
	fallback := func(c *cobra.Command, args []string) {
		// 已进入 fallback：直接渲染，避免无限递归。
		_ = c.Help()
	}
	_ = fallback

	// 通过 SetHelpFunc + 内部标志位避免递归。
	root.SetHelpFunc(func(c *cobra.Command, args []string) {
		if entered {
			// 二次进入：c.Help() 触发后再走 HelpFunc，已清空自己，
			// 应返回 Cobra 默认渲染。
			_ = c.Help()
			return
		}
		entered = true
		localizeBuiltinCommands(c.Root())
		// 清空自身 HelpFunc 后调用 Help()，c.Help() 会走 Cobra 默认逻辑。
		c.Root().SetHelpFunc(nil)
		_ = c.Help()
	})
}

// localizeBuiltinCommands 将 Cobra 自动注册的 completion / help 子命令
// 的 Short 描述本地化为中文。
func localizeBuiltinCommands(root *cobra.Command) {
	for _, c := range root.Commands() {
		switch c.Name() {
		case "completion":
			c.Short = "生成 shell 自动补全脚本"
		case "help":
			c.Short = "查看任意命令的帮助"
		}
	}
}

// HandleError 是测试场景下复刻 ExecuteWith 的错误处理路径：
// 将 err 渲染到 ErrorTemplate 后返回 ExitCode。
//
// 测试可使用此函数以避免直接调用 ExecuteWith 中的 os.Exit。
func HandleError(stderr io.Writer, err error) int {
	if err != nil {
		writeErrorTemplate(stderr, err.Error())
	}
	return FromError(err)
}

// renderVersionTo 将当前版本快照写入 w。
func renderVersionTo(w io.Writer) error {
	_, err := fmt.Fprint(w, version.Snapshot().String())
	return err
}

// writeErrorTemplate 向 w 写入 ErrorTemplate 渲染结果。
func writeErrorTemplate(w io.Writer, msg string) {
	if w == nil {
		return
	}
	out := strings.ReplaceAll(ErrorTemplate, "{{.Error}}", msg)
	_, _ = io.WriteString(w, out)
}
