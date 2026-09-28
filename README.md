# aiaikit

> 基于 Go 的 AI/SDD 工具生态平台 CLI。

`aiaikit` 是一个跨平台（Windows / Linux / macOS）Go CLI，用于分发、安装、版本管理与生态扩展 Spec-Kit Runtime 及其插件生态。本仓库当前为 **001 阶段**：CLI Core 骨架，仅交付 `aiaikit --version`、`aiaikit --help` 与五个预留命令入口；本阶段不实现任何 Spec-Kit 业务。

## 快速开始

```bash
# 开发模式（无需构建产物）
go run ./cmd/aiaikit --version
go run ./cmd/aiaikit --help

# 标准本地构建
make build
./bin/aiaikit --version
```

更多本地构建、跨平台构建、质量门禁与故障排查的完整步骤，请参见 [`specs/001-go-cli-bootstrap/quickstart.md`](specs/001-go-cli-bootstrap/quickstart.md)。

## 架构分层

本项目遵循 [`specs/001-go-cli-bootstrap/plan.md`](specs/001-go-cli-bootstrap/plan.md) 中的分层约束；当前仓库仅包含 CLI Core 层：

| 层 | 路径 | 本阶段状态 |
|---|---|---|
| CLI Core | `cmd/aiaikit/`、`internal/cli/`、`internal/version/`、`internal/platform/` | 已实现骨架 |
| Spec-Kit Runtime | — | 不属于本阶段（`runtime` 命令仅占位） |
| Plugin 系统 | — | 不属于本阶段（`plugin` 命令仅占位） |
| Web 管理中心 | — | 不属于本阶段（`web` 命令仅占位） |

## 质量门禁

```bash
make ci      # fmt + lint + test + build
make cover   # 生成覆盖率报告
make build-all  # 跨平台矩阵构建
```

CI 与本地使用同一入口，避免"本地过 / CI 红"。质量门禁覆盖：

- **格式化**：`gofmt -s -w` 与 `goimports -w`；
- **静态检查**：`golangci-lint`（启用 `govet` / `errcheck` / `staticcheck` / `unused` / `gofmt` / `goimports`）；
- **单元与 CLI 行为测试**：基于 Go 内置 `testing`；包含 golden 文件快照对比；
- **端到端测试**：调用真实二进制，验证 `--version`、`--help` 与五个预留命令；
- **构建**：单平台与跨平台矩阵。

## 合规

本项目遵循仓库根目录 `.specify/memory/constitution.md` 的章程，包括但不限于：

- **章程 I**：不重复实现 Spec-Kit 已有能力；本阶段不引用、不下载、不执行 Spec-Kit 任何代码；
- **章程 V**：Spec-Kit Runtime 与 aiaikit CLI 解耦；本阶段 CLI Core 不引入 Runtime 业务；
- **章程 VI**：插件系统是核心扩展机制；预留命令通过 `Annotations["reserved"]="true"` 标注，供后续 Plugin 系统按注解挂载；
- **章程 X**：跨平台一致；Windows / Linux / macOS 三平台构建通过统一 Makefile 入口；
- **章程 XV**：中文优先；所有面向用户的文案默认使用简体中文，版本号、命令名、路径、URL 等标识符保留原始形式。

章程一致性自检脚本：[`scripts/constitution-audit.sh`](scripts/constitution-audit.sh)。

## 贡献

提交前必须执行 [`make ci`](Makefile) 并确保本地全绿；详细约定参见 [`CONTRIBUTING.md`](CONTRIBUTING.md)。

## 许可

本项目使用 [MIT License](LICENSE)。