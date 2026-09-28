# 贡献指南

本项目遵循 [`constitution.md`](.specify/memory/constitution.md) 中的章程；任何 PR 提交前必须满足以下条件。

## 提交前自检

在本地仓库根目录执行：

```bash
make ci
```

`make ci` 是 CI 与本地共用的质量门禁入口，等价于：

- `make fmt` — gofmt + goimports
- `make lint` — golangci-lint（若未安装则回退到 `go vet`）
- `make test` — 单元测试 + CLI 行为测试 + 端到端测试
- `make build` — 当前平台构建

任何步骤失败都禁止提交。

## 契约变更

如修改 `internal/cli` 的输出（`--version` / `--help` / 预留命令），必须遵循以下流程：

1. **先修订契约**：更新 [`specs/001-go-cli-bootstrap/contracts/cli.md`](specs/001-go-cli-bootstrap/contracts/cli.md)；
2. **更新 golden 文件**：通过 `UPDATE_GOLDEN=1 go test ./...` 一次性重新生成快照；
3. **PR 描述中说明**：明确指出契约变更的影响面（哪些 User Story、哪些测试、哪些文档需要同步）；
4. **章程复核**：章程 I（不重复 Spec-Kit）、V（解耦）、XV（中文优先）三项的合规性需要在 PR 中显式确认。

## 分支与提交

- 默认分支：`main`；
- 特性分支：`<特性编号>-<短名>` 格式（例如 `001-go-cli-bootstrap`）；
- 提交消息：中文优先；遵循章程 XV；标识符与命令名保留原始形式。

## 章程引用

- **章程 II（核心分发对象）**：aiaikit 主要负责 Runtime 分发，不重复实现 Runtime 业务；新增 Runtime 相关能力必须落入 `internal/runtime` 包（后续阶段）；
- **章程 XII（新增能力优先扩展）**：本阶段 Core 入口（`cmd/aiaikit/main.go`、`internal/cli/root.go`）的修改必须附带 ADR 说明动机与扩展点影响；
- **章程 XV（中文优先）**：所有用户可见文案、错误信息、提交消息、PR 描述必须使用简体中文。