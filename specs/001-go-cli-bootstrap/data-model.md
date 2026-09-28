# Data Model：aiaikit 001 — Go 项目初始化与 CLI 基础框架

本文件定义本阶段核心实体的字段、关系与不变式。所有实体仅作为本阶段 CLI Core 的内部模型；不涉及持久化存储。

## 实体清单

| 实体 | 包路径 | 类别 |
|---|---|---|
| CLI 入口（Entry） | `cmd/aiaikit/main.go` | 进程级 |
| 命令节点（Command Node） | `internal/cli` | 注册时实例化 |
| 版本信息（Version Info） | `internal/version` | 链接期注入 |
| 平台信息（Platform Info） | `internal/platform` | 运行时探测 |
| 退出码（Exit Code） | `internal/cli` | 常量集 |
| 构建产物（Build Artifact） | —（Make 目标产出物） | 文件系统 |

---

## CLI 入口（Entry）

进程级单例；对应 `aiaikit` 二进制的 `main()`。

### 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| `os.Args` | `[]string` | 原始命令行参数 |
| `stdout` | `io.Writer` | 标准输出（默认 `os.Stdout`，测试可注入） |
| `stderr` | `io.Writer` | 标准错误（默认 `os.Stderr`） |
| `exitFn` | `func(int)` | 退出函数（默认 `os.Exit`，测试可注入） |

### 不变式

- 入口 MUST NOT 在执行任何 Cobra 命令前向 stdout / stderr 写入任何字节（避免污染 `--help` / `--version` 输出快照）。
- 入口 MUST 在 `cli.Execute()` 返回 `nil` 时通过 `exitFn(0)` 退出；返回非 `nil` 错误时通过 `exitFn(exitcode.FromError(err))` 退出。
- 入口 MUST NOT 捕获 panic；panic 沿 Go 默认机制传播，避免在 CLI Core 中实现额外 recovery。

---

## 命令节点（Command Node）

通过 Cobra 表达的子命令单元。本阶段包含以下实例：

| 命令 | 路径 | 状态 | 说明 |
|---|---|---|---|
| `root` | `aiaikit`（无参） | 已实现 | 默认输出帮助文本 |
| `--version` | `aiaikit --version` / `aiaikit version` | 已实现 | 输出版本信息 |
| `--help` | `aiaikit --help` / `aiaikit help` | 已实现 | 输出帮助文本 |
| `init` | `aiaikit init` | 预留 / 未实现 | 后续阶段：兼容 `specify init` |
| `update` | `aiaikit update` | 预留 / 未实现 | 后续阶段：升级 aiaikit 自身 |
| `runtime` | `aiaikit runtime` | 预留 / 未实现 | 后续阶段：Runtime 管理 |
| `plugin` | `aiaikit plugin` | 预留 / 未实现 | 后续阶段：插件管理 |
| `web` | `aiaikit web` | 预留 / 未实现 | 后续阶段：Web 管理中心 |

### 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| `Use` | `string` | 命令名（英文小写、连字符分隔） |
| `Short` | `string` | 中文短描述（≤ 40 字） |
| `Long` | `string` | 中文长描述（可选） |
| `Run` | `func(cmd *cobra.Command, args []string)` | 业务入口（本阶段仅预留命令的 Run 输出友好提示） |
| `Annotations` | `map[string]string` | 扩展元数据（如 `reserved: true`） |

### 关系

- 一个 `root` 节点 MUST 通过 `AddCommand` 持有所有其他命令节点；
- 预留命令 MUST 在 `Annotations` 中标注 `"reserved": "true"`，供 `internal/cli/reserved.go` 集中遍历；
- `root.Run` 在未提供子命令时 MUST 等价于 `root.Help()`，并以退出码 `0` 退出。

### 不变式

- 命令节点的 `Short` / `Long` MUST 使用中文（章程 XV）；
- 命令节点的 `Use` MUST 使用英文小写命令名（与 Spec-Kit 命名空间无冲突）；
- 预留命令的 `Run` MUST 输出中文友好提示并以退出码 `2` 退出（参数错误类语义）。

---

## 版本信息（Version Info）

链接期注入的全局只读快照；定义于 `internal/version` 包。

### 字段

| 字段 | 类型 | 默认值 | ldflags 注入键 |
|---|---|---|---|
| `Version` | `string` | `"dev"` | `github.com/aiaikit/aiaikit/internal/version.Version` |
| `Commit` | `string` | `"unknown"` | `github.com/aiaikit/aiaikit/internal/version.Commit` |
| `BuildTime` | `string` | `"unknown"`（RFC3339） | `github.com/aiaikit/aiaikit/internal/version.BuildTime` |
| `Platform` | `string` | `runtime.GOOS + "/" + runtime.GOARCH` | `github.com/aiaikit/aiaikit/internal/version.Platform` |
| `GoVersion` | `string` | `runtime.Version()` | `github.com/aiaikit/aiaikit/internal/version.GoVersion` |

### 关系

- `aiaikit --version` MUST 消费此实体并按 contracts/cli.md §1 输出格式；
- `Platform` 与 `GoVersion` 由程序运行时填入（不通过 ldflags 注入），避免跨平台构建时手工维护。

### 不变式

- 五个字段 MUST 全部导出（小写变量无法被 ldflags 覆盖）；
- 字段值 MUST 为字符串原子赋值兼容（`ldflags -X` 仅支持字符串标量）；
- 默认值 MUST NOT 为空字符串；任何字段缺失构建参数时显示 `dev` / `unknown`（spec Edge Case）。

---

## 平台信息（Platform Info）

运行时探测结果；定义于 `internal/platform` 包。

### 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| `OS` | `string` | `runtime.GOOS` 取值（`windows` / `linux` / `darwin` / ...） |
| `Arch` | `string` | `runtime.GOARCH` 取值（`amd64` / `arm64` / ...） |
| `IsUnixLike` | `bool` | Unix 类系统判定（`unix` 编译标签与 `darwin` / `linux` / `freebsd` 等） |
| `IsWindows` | `bool` | Windows 判定 |

### 关系

- 本阶段仅用于 `--version` 输出与未来跨平台路径 / 进程适配；
- `platform_unix.go` / `platform_windows.go` 通过 build tags 隔离实现，避免 `runtime.GOOS` 分支污染调用方。

### 不变式

- 本阶段 MUST NOT 引入任何依赖平台信息的具体行为分支（避免无谓复杂度）；所有字段仅供信息展示与后续阶段使用。

---

## 退出码（Exit Code）

`internal/cli/exitcode.go` 中定义的常量集。

### 字段

| 常量 | 值 | 语义 |
|---|---|---|
| `ExitSuccess` | `0` | 成功 |
| `ExitUsage` | `2` | 用法 / 参数错误（含未实现命令） |
| `ExitGeneral` | `1` | 一般运行时错误（本阶段保留） |
| `ExitReservedStart` | `3` | 后续阶段分配起点 |

### 关系

- `cmd.Execute()` 错误 MUST 经 `exitcode.FromError(err)` 转换后调用 `exitFn`；
- Cobra 默认 `SilenceUsage = true` + `SilenceErrors = true`，确保错误信息走中文包装函数。

### 不变式

- 退出码值 MUST 落在 `[0, 125]` 区间（保留 `126` / `127` 给 shell 自身使用）。

---

## 构建产物（Build Artifact）

Make 目标产出的文件系统对象。

### 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| `Path` | `string` | 产物路径（默认 `./bin/aiaikit` 或 `./bin/aiaikit-{GOOS}-{GOARCH}`） |
| `GOOS` / `GOARCH` | `string` | 目标平台 |
| `VersionInfo` | 嵌入二进制 | 由 ldflags 注入的版本信息 |

### 关系

- 本阶段不维护产物清单文件；CI 通过 `bin/` 目录扫描获得。

### 不变式

- 产物 MUST 在三类平台上均可执行（同源代码 + 同 Makefile）；
- 产物 MUST 内嵌完整版本信息（含 `Commit`、`BuildTime`、`Platform`、`GoVersion`）。

---

## 实体关系总览

```text
Entry (main.go)
  └── cli.Root (Command Node, root)
       ├── version 子命令 (Command Node, version) ──consumes──> Version Info
       ├── help 子命令 (Command Node, help)
       ├── init / update / runtime / plugin / web (Command Node, reserved)
       └── --version / --help (Cobra PersistentFlag)
       
Version Info
  ├── consumed by version 子命令
  └── consumed by --version 标志

Platform Info
  ├── consumed by Version Info (Platform 字段)
  └── consumed by 未来跨平台适配层（占位）

Exit Code
  └── consumed by Entry (exitFn)

Build Artifact
  └── embeds Version Info (via ldflags)
```