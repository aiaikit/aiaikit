package cli

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/spf13/cobra"
)

// ErrReserved 是预留命令未实现时返回的哨兵错误。
//
// ExecuteWith 通过 exitcode.FromError 将其映射为 ExitUsage。
var ErrReserved = errors.New("命令尚未实现")

// RegisterReserved 在 root 命令下注册五个预留命令。
//
// 预留命令的当前阶段行为：
//   - 输出 ReservedMessageTemplate 提示用户尚未实现；
//   - 退出码为 ExitUsage（章程 SC 与 spec FR 一致）；
//   - 不得下载或执行任何 Spec-Kit Runtime 代码；
//   - 不得修改磁盘上的任何文件；
//   - 在 Annotations 中标注 reserved=true，供后续 Plugin 系统按注解挂载真实实现。
//
// 本函数是章程 VI / XII 要求的扩展点入口；后续 Plugin 阶段将通过同
// 一机制把真实实现挂载到这些预留命令上。
func RegisterReserved(root *cobra.Command) {
	commands := []struct {
		name  string
		short string
	}{
		{"init", "初始化 aiaikit 工作区"},
		{"update", "升级 aiaikit 自身"},
		{"runtime", "Spec-Kit Runtime 管理"},
		{"plugin", "插件管理"},
		{"web", "Web 管理中心"},
	}

	for _, c := range commands {
		cmd := buildReservedCmd(c.name, c.short)
		root.AddCommand(cmd)
	}
}

// buildReservedCmd 构造单个预留命令实例。
func buildReservedCmd(name, short string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short + " " + ReservedSuffix,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := renderReserved(cmd.OutOrStdout(), name); err != nil {
				return err
			}
			return ErrReserved
		},
		Annotations: map[string]string{
			ReservedAnnotationKey: ReservedAnnotationValue,
		},
	}
}

// renderReserved 向 w 写入预留命令的中文友好提示。
func renderReserved(w interface {
	Write(p []byte) (n int, err error)
}, name string) error {
	tmpl, err := template.New("reserved").Parse(ReservedMessageTemplate)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		Command string
	}{Command: name}); err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}
