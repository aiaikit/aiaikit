# CLI 命令契约：aiaikit 001

本文件定义本阶段所有命令的输入、输出与退出码契约。CLI 行为测试与 golden 文件 MUST 以本契约为唯一真源；任何对契约的偏离 MUST 先升级本文件再改实现。

> 字符集约定：所有面向用户的可读字段 MUST 使用简体中文（章程 XV）。版本号、命令名、路径、URL、错误码、`GOOS` / `GOARCH`、Commit Hash 等标识符保留原始形式。

---

## §1 `aiaikit --version`

### 输入

```
aiaikit --version
aiaikit version        # 等价
```

### 输出（stdout）

```
aiaikit 版本 <Version>
  Commit   : <Commit>
  构建时间 : <BuildTime RFC3339>
  目标平台 : <Platform>
  Go 版本  : <GoVersion>
```

示例：

```
aiaikit 版本 0.1.0
  Commit   : a1b2c3d4e5f6
  构建时间 : 2026-09-28T10:15:30Z
  目标平台 : linux/amd64
  Go 版本  : go1.23.0
```

缺省值（开发模式 / 未注入）：

```
aiaikit 版本 dev
  Commit   : unknown
  构建时间 : unknown
  目标平台 : <runtime.GOOS>/<runtime.GOARCH>
  Go 版本  : <runtime.Version()>
```

### 输出（stderr）

无（命令成功路径下 MUST 不向 stderr 写入任何字节）。

### 退出码

- `0`：成功。

### 不变式

- 五个字段 MUST 全部输出，且字段顺序固定；
- 字段值 MUST 保留原始大小写与字符（不得做 trim / 转义 / 大小写转换）；
- 默认值 `dev` / `unknown` 与运行时探测的 `Platform` / `GoVersion` 不得为空字符串。

---

## §2 `aiaikit --help`

### 输入

```
aiaikit --help
aiaikit help           # 等价
aiaikit               # 无子命令时等价
```

### 输出（stdout）

Cobra 自动生成；结构固定：

```
aiaikit —— 基于 Go 的 AI/SDD 工具生态平台 CLI

Usage:
  aiaikit [command]

Available Commands:
  completion  生成 shell 自动补全脚本
  help        查看任意命令的帮助
  init        初始化 aiaikit 工作区（预留 / 未实现）
  plugin      插件管理（预留 / 未实现）
  runtime     Spec-Kit Runtime 管理（预留 / 未实现）
  update      升级 aiaikit 自身（预留 / 未实现）
  version     查看 aiaikit 版本信息
  web         Web 管理中心（预留 / 未实现）

Flags:
  -h, --help      查看命令帮助
  -v, --version   查看版本信息

Use "aiaikit [command] --help" for more information about a command.
```

### 输出（stderr）

无（`--help` / 无子命令路径 MUST NOT 向 stderr 写入字节）。

### 退出码

- `0`：成功。

### 不变式

- 八个 Available Commands MUST 全部出现：`completion`、`help`、`init`、`plugin`、`runtime`、`update`、`version`、`web`；
- 预留命令的简短描述 MUST 以“（预留 / 未实现）”结尾；
- 顺序固定（completion 在前，web 在后），便于快照测试稳定；
- Cobra 生成的 “Usage:”、“Available Commands:”、“Flags:”、“Use ...” 等小标题 MUST 保持英文（Cobra 框架默认行为，避免影响快照稳定性）。

---

## §3 预留命令（init / update / runtime / plugin / web）

### 输入

```
aiaikit init
aiaikit update
aiaikit runtime [...]
aiaikit plugin [...]
aiaikit web [...]
```

可附带任意参数 / flag；本阶段 MUST 接受全部输入并在最后输出固定友好提示。

### 输出（stdout）

```
命令 "<cmd>" 尚未实现，将在后续阶段提供。

提示：本命令当前为预留状态，不会下载或执行任何 Spec-Kit Runtime 代码。
```

`<cmd>` 为用户实际调用的命令名。

### 输出（stderr）

无（本阶段保留为空）。

### 退出码

- `2`：退出码语义为“用法 / 参数错误 / 未实现”，符合 §3 不变式。

### 不变式

- MUST NOT 触发网络请求、文件下载、进程派生、Spec-Kit Runtime 加载；
- MUST NOT 修改磁盘上的任何文件；
- MUST NOT 在 stdout / stderr 之外的渠道（IPC、日志文件）写入与本次调用相关的信息（除非未来扩展明确允许）。

---

## §4 `aiaikit help <subcommand>`

### 输入

```
aiaikit help <subcommand>
aiaikit <subcommand> --help
```

### 输出（stdout）

Cobra 自动生成：每个子命令的 Usage / Flags / 示例（若 Long 描述中存在）。

### 输出（stderr）

无。

### 退出码

- `0`：成功；
- `2`：未知子命令（Cobra 默认行为）。

### 不变式

- 预留子命令的 help MUST 出现“（预留 / 未实现）”字样。

---

## §5 错误路径（占位）

> 本阶段不引入除 §3 之外的错误路径；任何运行时错误（如配置文件缺失、网络异常）推迟到后续阶段。本节作为占位，便于后续契约扩展。

### 输入

任意未涵盖的非法输入。

### 输出（stderr）

```
错误：<中文错误描述>
用法：aiaikit [command]
查看帮助：aiaikit --help
```

### 退出码

- `2`：用法 / 参数错误；
- `1`：一般运行时错误（本阶段保留）。

### 不变式

- 错误信息 MUST 使用中文；
- MUST NOT 暴露内部包路径、堆栈、敏感信息。

---

## §6 全局约定

| 项 | 约定 |
|---|---|
| 字符编码 | UTF-8 |
| 行尾 | `\n`（LF），与 Go 默认 fmt.Println 一致 |
| 表头标签 | 中文为主（Cobra 框架默认英文小标题保留） |
| 标识符 | 命令名、flag、版本号、路径、URL、Commit Hash 保留原始形式 |
| 帮助输出稳定性 | 顺序、字段集合、字段顺序固定（snapshot 测试以此为准） |

---

## §7 版本演进

新增命令 / 修订现有命令时：
1. MUST 先修订本契约；
2. 更新 `testdata/cli/*.golden` 通过 `UPDATE_GOLDEN=1 go test ./...` 触发；
3. 在 PR 描述中说明契约变更与影响面；
4. 章程 I / V / XII 复核：避免 Spec-Kit 命名空间冲突；避免绕过扩展点直接修改 Core。