# Research：aiaikit 001 — Go 项目初始化与 CLI 基础框架

**Phase**: 0（研究 / 决策）
**目标**：把 spec 中的开放决策收敛为可执行的技术选型；所有 NEEDS CLARIFICATION 在本文件中解决。

## D-001：Go 版本基线

- **决策**：`go.mod` 中 `go` 指令采用 Go 官方当前稳定主版本；最低支持版本与 Go 官方支持矩阵对齐（即官方同时维护的两个最近主版本中较低者）。
- **理由**：章程 IV 要求 aiaikit 通过 npm 分发，CLI 自身需要足够新的 Go 工具链支持；但锁定过旧版本会拖累后续阶段使用现代语法（`slices`、`maps`、`cmp` 等）。
- **替代**：
  - 锁最低可用版本（如 Go 1.21）：换范围更广，但错过近期语法糖。
  - 锁最新版本：编译期风险较高，老用户升级工具链成本大。
- **落地**：`go.mod` 的 `go` 指令与实际开发机一致，CI 三平台 runner 镜像锁定同一 Go 主版本。

## D-002：CLI 框架（Cobra vs stdlib flag）

- **决策**：使用 `github.com/spf13/cobra` 作为 CLI 命令树框架。
- **理由**：
  1. Cobra 是 Go 生态事实标准（kubectl / docker / hugo / gh 等），用户熟悉度高；
  2. Cobra 原生支持子命令、自动生成帮助、统一退出码管理，与 spec FR-002 / FR-003 强匹配；
  3. Cobra 的 `AddCommand` / `PersistentPreRun` / `SilenceUsage` 等扩展点为后续 Plugin / Runtime 阶段提供稳定接口，符合章程 VI / XII；
  4. 单文件代码体量可控（核心包约 6k 行），不显著影响 SC-006 二进制大小基线。
- **替代**：
  - Go 标准库 `flag`：可避免外部依赖，但子命令树、自动帮助、退出码管理需要自行实现，扩展成本高于直接复用 Cobra；
  - `urfave/cli`：API 较 Cobra 更“重”，生态略小。
- **落地**：`go.mod` 引入 `github.com/spf13/cobra`；root 命令与预留命令的实现严格走 Cobra API。

## D-003：版本信息注入

- **决策**：通过链接器 `-ldflags "-X 'github.com/aiaikit/aiaikit/internal/version.Version=...'"` 注入版本信息；版本信息集中在 `internal/version` 包，导出变量为 `Version`、`Commit`、`BuildTime`、`Platform`、`GoVersion`。
- **理由**：
  1. Go 标准做法，Kubernetes / Docker 等大型 Go CLI 普遍采用；
  2. 所有命令（包括未来新增命令）共享同一份版本信息源；
  3. ldflags 是字符串原子赋值，无需代码改动；
  4. 默认值（`dev` / `unknown`）保证本地 `go run` 与缺失构建参数时不会崩溃（spec Edge Cases）。
- **替代**：
  - 在编译期通过 `runtime.ReadBuildInfo()` 读取：可读出 vcs 信息，但自定义字段（构建时间、平台）需要额外机制，不如 ldflags 直接；
  - 嵌入文件（`go:embed` + `json`）：引入额外构建步骤，不必要。
- **落地**：`Makefile` 中的 `build` 目标统一注入五个字段；开发模式保留 `dev` / `unknown`。

## D-004：跨平台构建脚本

- **决策**：以 GNU Make 为本地主入口；CI 在三类 runner 上各自调用同一份 `Makefile` 目标；不引入 GoReleaser（推迟到发布阶段）。
- **理由**：
  1. Makefile 是 Unix 工具链标配；Windows 通过 Git for Windows / WSL 上的 `make` 即可覆盖；
  2. 跨平台 `GOOS` / `GOARCH` 矩阵由 Make 变量驱动，与 Go 官方推荐方式一致；
  3. 不引入额外工具链，避免推迟本阶段交付。
- **替代**：
  - shell + PowerShell 双脚本：维护成本高；
  - GoReleaser / Task / Just：推迟到发布与 CI 统一化阶段。
- **落地**：Makefile 提供 `build`、`build-all`、`test`、`lint`、`fmt`、`cover`、`clean` 七个目标；跨平台构建以变量 `GOOS` / `GOARCH` 覆盖。

## D-005：测试策略

- **决策**：
  1. 单元测试：Go 内置 `testing`，包级别覆盖 `internal/cli` / `internal/version` / `internal/platform`；
  2. CLI 行为测试：基于 Cobra 的 `cmd.SetArgs` + `cmd.Execute()` + `cmd.SetOut/SetErr` 模式，对 `aiaikit --version`、`aiaikit --help`、预留子命令进行输出快照对比（golden 文件位于 `testdata/cli/`）；
  3. 端到端测试：`e2e/` 目录，调用 `go build` 产出的真实二进制，断言退出码与关键字段；本阶段可仅在 Linux runner 上跑，其他平台由 CI 矩阵触发。
- **理由**：
  1. stdlib `testing` 零依赖；
  2. golden 文件对“输出格式稳定”提供最直接保障（章程 SC-002 一致性要求）；
  3. 端到端测试覆盖真实二进制，避免 `cmd.Execute()` 路径与 `main()` 路径不一致。
- **替代**：
  - `testify`：断言库风格偏重；本阶段规模不需要；
  - `ginkgo` / `gocheck`：BDD 风格，对 CLI 场景不必要。
- **落地**：每个内部包同名 `_test.go`；`testdata/cli/*.golden` 为只读快照，更新通过 `UPDATE_GOLDEN=1 go test ...` 显式触发。

## D-006：静态检查工具链

- **决策**：`golangci-lint`（latest stable）作为静态检查入口；启用 `govet`、`errcheck`、`staticcheck`、`unused`、`gofmt`、`goimports` 五个 linter（默认启用集合的子集）。
- **理由**：
  1. `golangci-lint` 是 Go 社区静态检查事实标准；
  2. 启用集合聚焦“正确性 + 一致性”，避免引入风格偏好争议；
  3. 与 `go vet` / `gofmt` 兼容，不冲突。
- **替代**：
  - 仅 `go vet`：覆盖度不够；
  - 直接引入 `staticcheck`：CI 集成成本更高。
- **落地**：`.golangci.yml` 显式列出启用的 linter；Makefile `lint` 目标调用 `golangci-lint run ./...`。

## D-007：中文文案在 CLI 中的落地

- **决策**：所有面向用户的字符串（`--help`、`--version`、错误提示、usage）集中在 `internal/cli` 的中文常量与模板；字符串统一使用 Go 原始字符串字面量（反引号），避免转义引入 bug。
- **理由**：
  1. 章程 XV 要求中文优先；同时保留版本号、命令名、路径、URL、错误码为原始形式；
  2. 集中常量便于后续 i18n 演进（即使本阶段不交付多语言）；
  3. 原始字符串字面量天然支持中文，避免转义错误。
- **替代**：
  - `golang.org/x/text/message`：i18n 库，本阶段不交付多语言；
  - 外部 `.po` / `.json` 资源文件：增加分发复杂度。
- **落地**：`internal/cli/messages.go` 集中声明所有面向用户的中文字符串；代码注释与 commit message 同样遵守章程 XV。

## D-008：退出码与错误约定

- **决策**：
  - `0`：成功；
  - `2`：参数或用法错误（含 `--help` / `--version` 等纯信息命令之外的命令参数错误）；
  - `1`：一般运行时错误（本阶段不触发，仅占位）；
  - 预留 `3` ~ `127`：供后续阶段分配。
- **理由**：
  1. 退出码语义与 CLI 行业惯例（sysexits.h、Go `os.Exit` 通用约定）对齐；
  2. 留出大量未分配空间，便于后续 Runtime / Plugin / Web 阶段扩展；
  3. 章程 SC 不强制退出码数值，但后续 CI 与运维脚本需要稳定约定。

## D-009：CI 与本地构建一致性

- **Decision**：CI 在 Windows / Linux / macOS runner 上各自执行 `make ci`（包含 `fmt + lint + test + build`）；本地开发者同样执行 `make ci` 作为提交前自检。
- **理由**：
  1. 单入口降低跨平台差异；
  2. CI 与本地使用同一 Makefile，避免“本地过 / CI 红”问题。
- **落地**：本阶段 CI 配置文件（`.github/workflows/*.yml` 或等价）不在 plan 范围内；后续阶段由独立 spec 覆盖。

## Open Questions

- **None**：本阶段无遗留 NEEDS CLARIFICATION；所有 spec 中保留的假设均已在 plan 与 research 中收敛。

## 章程一致性再确认（Phase 0 之后）

- I / II / V：CLI Core 与 Spec-Kit Runtime / Plugin / Web 单向解耦，依赖图中无 Spec-Kit 模块 → 通过。
- VI / XII：扩展点（Cobra `AddCommand`）已明确 → 通过。
- X：跨平台适配占位 + Makefile 矩阵 → 通过。
- XV：中文优先 + 标识符保留原始形式 → 通过。