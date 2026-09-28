# Implementation Plan: aiaikit 001 — Go 项目初始化与 CLI 基础框架

**Branch**: `001-go-cli-bootstrap` | **Date**: 2026-09-28 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-go-cli-bootstrap/spec.md`

## Summary

为 aiaikit 建立一个可跨平台（Windows / Linux / macOS）构建的 Go CLI 骨架。本阶段交付：
1. 标准 Go 模块布局与依赖最小化的命令注册结构；
2. `aiaikit --version`、`aiaikit --help` 两个必交付能力；
3. `init`、`update`、`runtime`、`plugin`、`web` 五个预留命令入口（仅占位，不实现业务）；
4. 通过链接期 `-ldflags "-X"` 注入语义化版本号、构建时间、目标平台；
5. Makefile + golangci-lint + Go 内置 `testing` 组成的标准质量门禁；
6. 为后续 Runtime / Plugin / Web 阶段保留稳定的扩展接口（命令注册、版本包、跨平台适配层），本阶段不实现其业务。

明确不在本阶段交付：Spec-Kit Runtime 任何代码、Plugin 加载机制、Web 后端 / 前端、`aiaikit update` 真实下载、npm 发布链路。

## Technical Context

**Language/Version**: Go（采用 Go 官方当前稳定主版本，最低支持版本对齐 Go 官方支持矩阵；具体版本在 `go.mod` 中以 `go` 指令锁定，详见 research.md 决策 D-001）

**Primary Dependencies**:
- `github.com/spf13/cobra`：CLI 命令树与帮助生成（决策 D-002）
- `github.com/spf13/pflag`：Cobra 推荐的 POSIX 风格 flag（由 Cobra 间接引入）
- `github.com/spf13/pflag` 之外的运行依赖：刻意保持空；任何额外依赖进入前必须经过章程 XIII 评审

**Storage**: N/A（本阶段不涉及持久化）

**Testing**:
- Go 内置 `testing` 包作为单元测试基座
- `github.com/spf13/cobra` 自带 `cmd.SetArgs` + `cmd.Execute()` 模式用于 CLI 行为测试
- 集成 / 端到端测试：通过 `testscript`（Go 内置 `cmd/go` test 工具集中的 `Testscript` 实验性支持）或独立的 `e2e/` 目录调用真实构建产物
- 覆盖率：内置 `go test -cover` 输出 `coverage.out`

**Target Platform**: Windows（amd64 / arm64）、Linux（amd64 / arm64）、macOS（amd64 / arm64）；CI 上至少跑 amd64 三平台

**Project Type**: 单体 CLI 项目（`cmd` + `internal` 标准布局）

**Performance Goals**: 单次 `aiaikit --version` 启动至输出耗时 P95 ≤ 200ms（冷启动场景下，I/O 缓存命中时）；不强制，但作为 Go CLI 通用基线

**Constraints**:
- 产物二进制 ≤ 30MB（章程 SC-006 与 spec 一致）
- 不引入 CGO，避免跨平台 C 工具链依赖
- 不在二进制中捆绑 Spec-Kit 任何代码、模板、脚本（spec FR-008）
- 默认输出编码 UTF-8，跨平台控制台输出兼容当前平台默认编码

**Scale/Scope**: 单一可执行文件；本阶段代码量预期 ≤ 1500 行（含测试与 Makefile），不预设用户量级

## Constitution Check

*GATE: 必须先于 Phase 0 research 通过；Phase 1 设计后再次复核。*

| 章程条款 | 检查结论 | 说明 |
|---|---|---|
| I. 不重复实现 Spec-Kit | ✅ 通过 | 本阶段不引用、不下载、不执行 Spec-Kit 任何代码；`go.mod` 不得包含 Spec-Kit 自身模块 |
| II. Runtime 是核心分发对象 | ✅ 通过 | 本阶段仅做 CLI Core 骨架，Runtime 目录预留，不引入 Runtime 业务 |
| III. `aiaikit init` 兼容 `specify init` | ✅ 通过 | 兼容性属于后续阶段范围，本阶段仅预留 `init` 命令入口并在 `--help` 标注“预留 / 未实现” |
| IV. aiaikit 通过 npm 分发 | ✅ 通过 | 本阶段交付 Go 二进制与 Makefile；npm 发布链路推迟到后续阶段，并在 plan 中显式 Out of Scope |
| V. Runtime 与 CLI 解耦 | ✅ 通过 | 目录分层保证 `cmd/aiaikit` 与未来的 `internal/runtime` / `internal/plugin` 单向依赖 |
| VI. 插件系统是核心扩展机制 | ✅ 通过 | 通过 Cobra 子命令 + 显式 `RegisterReserved` 扩展点为后续 Plugin 预留接口 |
| VII. Prompt Generator 作为首个内置插件 | ✅ 通过 | 推迟到 Plugin 阶段；本阶段 `internal/cli/reserved.go` 不创建 `prompt` 子命令 |
| VIII. 保留原生斜杠命令 | ✅ 通过 | 本阶段不实现任何斜杠命令；后续接入 aiaikit 命名空间（`/speckit-` 等），不与 Spec-Kit 冲突 |
| IX. Web 复用 Core | ✅ 通过 | `web` 子命令仅占位，目录预留但不实现；后续复用 `internal/cli` 注册点 |
| X. 跨平台与安全基线 | ✅ 通过 | 构建脚本三类平台可执行；引入跨平台适配占位（`internal/platform`） |
| XI. 资源版本、依赖、兼容性管理 | ✅ 通过 | 版本信息集中注入；`-ldflags` 注入字段必须为字符串原子赋值 |
| XII. 新增能力优先扩展 | ✅ 通过 | 所有非 Core 能力（Runtime / Plugin / Web / Prompt Generator）以预留命令形式占位，不污染 Core |
| XIII. 全栈架构纳入范围 | ✅ 通过 | 本阶段聚焦 CLI；后端 / 前端 / DB 在后续阶段独立 spec |
| XIV. 与官方生态兼容 | ✅ 通过 | 不修改 Spec-Kit 任何代码；Cobra 为通用 Go CLI 生态，无 Spec-Kit 耦合 |
| XV. 中文优先 | ✅ 通过 | 所有错误信息、`--help` 输出、CLI 提示 MUST 使用中文；版本号、命令名、路径保留原始形式 |

复核结论：本阶段不触发任何章程违规项，不需要 Complexity Tracking 表。

## Project Structure

### Documentation (this feature)

```text
specs/001-go-cli-bootstrap/
├── plan.md              # 本文件
├── research.md          # Phase 0 产出
├── data-model.md        # Phase 1 产出
├── quickstart.md        # Phase 1 产出
├── contracts/           # Phase 1 产出
│   └── cli.md           # CLI 命令契约（版本输出、帮助输出、退出码）
└── tasks.md             # Phase 2 产出（由 /speckit-tasks 生成）
```

### Source Code (repository root)

```text
aiaikit/
├── cmd/
│   └── aiaikit/
│       └── main.go              # 入口：调用 internal/cli.Execute() 并处理退出码
├── internal/
│   ├── cli/
│   │   ├── root.go              # root command：注册 --version、--help、全局 flag
│   │   ├── reserved.go          # 预留命令：init / update / runtime / plugin / web
│   │   ├── version.go           # 版本子命令（--version 等价）
│   │   └── exitcode.go          # 退出码常量与转换
│   ├── version/
│   │   └── version.go           # 全局版本信息；变量由链接器覆盖
│   └── platform/
│       ├── platform.go          # 平台检测接口
│       ├── platform_unix.go     //build unix
│       └── platform_windows.go  //build windows
├── testdata/
│   └── cli/
│       ├── version.golden       # golden 文件：版本输出快照
│       └── help.golden          # golden 文件：帮助输出快照
├── .golangci.yml                # 静态检查配置
├── Makefile                     # build / test / lint / fmt / cover 目标
├── go.mod
├── go.sum
├── README.md                    # 中文 README；开发与构建说明
└── LICENSE
```

**Structure Decision**: 采用 “Single project (CLI)” 单项目布局。理由：
1. 本阶段只交付一个 CLI 二进制，无服务端 / 前端配套；
2. `cmd/` 入口与 `internal/` 私有包分层符合 Go 官方推荐；
3. `internal/cli`、`internal/version`、`internal/platform` 三个内部包严格遵循章程“Architecture Constraints”中 Core 层职责划分；
4. 后续 Runtime / Plugin / Web 阶段以新增 `internal/runtime`、`internal/plugin` 等同级包的方式增量扩展，不修改本阶段目录。

## Complexity Tracking

> *当前 Constitution Check 无违规项，本节留空。*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| —         | —          | —                                   |

## Phase 0 Artifacts

- `research.md`：包含 D-001 Go 版本选择、D-002 CLI 框架选型（Cobra vs stdlib flag）、D-003 版本注入方式、D-004 跨平台构建脚本形式、D-005 测试策略、D-006 静态检查工具链、D-007 中文文案在 CLI 中的落地方式。
- 所有 NEEDS CLARIFICATION 已在 research.md 收敛为决策，不遗留开放问题。

## Phase 1 Artifacts

- `data-model.md`：CLI 入口、命令节点、版本信息、构建产物、退出码五个核心实体的字段与关系。
- `contracts/cli.md`：CLI 命令契约（`aiaikit --version`、`aiaikit --help`、预留子命令行为的输入 / 输出 / 退出码规范）。
- `quickstart.md`：本地构建、跨平台构建、质量门禁、CLI 行为验证的可复现步骤。

## Phase 2 Handoff

- 由 `/speckit-tasks` 接管：基于本 plan 与 research.md 拆解为可执行任务清单并落 `tasks.md`。
- 实现阶段不重新打开架构决策；如需变更，须先升级本 plan 并重新通过章程检查。