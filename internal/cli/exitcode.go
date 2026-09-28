package cli

import "errors"

// 退出码常量集。
//
// 数值落在 [0, 125] 区间；126 / 127 保留给 shell 自身使用。
// 3 ~ 125 留作后续 Runtime / Plugin / Web 阶段分配。
const (
	// ExitSuccess 表示命令执行成功。
	ExitSuccess = 0
	// ExitGeneral 表示一般运行时错误（本阶段保留，未来使用）。
	ExitGeneral = 1
	// ExitUsage 表示用法、参数错误或调用未实现命令。
	ExitUsage = 2
	// ExitReservedStart 表示后续阶段分配退出码的起点。
	ExitReservedStart = 3
)

// FromError 将任意错误映射到退出码。
//
// 本阶段策略：
//   - err == nil                       → ExitSuccess
//   - errors.Is(err, ErrReserved)       → ExitUsage（未实现预留命令）
//   - 其他错误                          → ExitUsage（参数错误 / 用法错误）
//
// 后续阶段可在此函数中按错误类型路由到不同退出码。
func FromError(err error) int {
	if err == nil {
		return ExitSuccess
	}
	if errors.Is(err, ErrReserved) {
		return ExitUsage
	}
	return ExitUsage
}
