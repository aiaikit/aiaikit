package cli

import (
	"github.com/spf13/cobra"

	"github.com/aiaikit/aiaikit/internal/version"
)

// newVersionCmd 构建 `aiaikit version` 子命令。
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "查看 aiaikit 版本信息",
		Long:  "输出当前 aiaikit 二进制的版本号、Commit、构建时间、目标平台与 Go 版本。",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cmd.OutOrStdout().Write([]byte(version.Snapshot().String()))
			return err
		},
	}
}
