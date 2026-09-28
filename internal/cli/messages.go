package cli

// 本文件集中声明所有面向用户的中文字符串（章程 XV）。
//
// 规则：
//   - 任何面向 stdout / stderr 的中文字符串 MUST 通过本文件提供；
//   - 命令名、flag、路径、URL、版本号、错误码等标识符保留原始形式；
//   - 使用 Go 原始字符串字面量（反引号）保留换行与排版。

const (
	// AppName 二进制名；保留原始英文以匹配可执行文件名。
	AppName = "aiaikit"

	// AppLongDescription 是 root 命令的长描述，出现在 --help 顶部。
	AppLongDescription = "aiaikit —— 基于 Go 的 AI/SDD 工具生态平台 CLI"

	// ReservedMessageTemplate 渲染预留命令的友好提示。
	// 占位 {{.Command}} 替换为实际命令名。
	ReservedMessageTemplate = `命令 "{{.Command}}" 尚未实现，将在后续阶段提供。

提示：本命令当前为预留状态，不会下载或执行任何 Spec-Kit Runtime 代码。
`

	// ErrorTemplate 渲染通用错误信息（带用法提示）。
	// 占位 {{.Error}} 替换为具体错误描述。
	ErrorTemplate = `错误：{{.Error}}
用法：aiaikit [command]
查看帮助：aiaikit --help
`

	// ReservedSuffix 标记命令为预留 / 未实现，统一附加到 Short 描述末尾。
	ReservedSuffix = "（预留 / 未实现）"

	// ReservedAnnotationKey 在 Cobra Annotations 中标注预留命令，供后续 Plugin 系统识别。
	ReservedAnnotationKey = "reserved"
	// ReservedAnnotationValue 预留命令的标注值。
	ReservedAnnotationValue = "true"
)
