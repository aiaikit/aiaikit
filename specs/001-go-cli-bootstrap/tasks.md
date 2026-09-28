---
description: "Task list for aiaikit 001 — Go 项目初始化与 CLI 基础框架"
---

# Tasks: aiaikit 001 — Go 项目初始化与 CLI 基础框架

**Input**: Design documents from `/specs/001-go-cli-bootstrap/`
- plan.md（required）
- spec.md（required；4 条 User Story，2×P1 + 2×P2）
- research.md（9 项决策 D-001 ~ D-009）
- data-model.md（6 个核心实体）
- contracts/cli.md（CLI 命令契约 7 节）
- quickstart.md（验证脚本 10 节）

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli.md, quickstart.md

**Tests**: spec FR-007 / FR-008 要求基础单元测试与 CLI 行为测试；本任务清单纳入测试任务，但不强求 TDD 红绿循环，按 quickstart.md §6 验收。

**Organization**: 任务按 user story 分组，便于独立实现与独立测试。

## Format: `[ID] [P?] [Story] Description with file path`

- **[P]**：可并行执行（不同文件，无依赖）
- **[Story]**：所属用户故事（US1 / US2 / US3 / US4）
- 描述 MUST 包含精确文件路径

---

## Phase 1: Setup（项目初始化）

**Purpose**: 建立仓库骨架与基础元数据。

- [ ] T001 初始化 Go 模块：在文件 `go.mod` 中声明模块路径 `github.com/aiaikit/aiaikit`，锁定 Go 主版本与 `go` 指令（决策 D-001）
- [ ] T002 [P] 创建目录骨架：`cmd/aiaikit/`、`internal/cli/`、`internal/version/`、`internal/platform/`、`testdata/cli/`、`bin/`
- [ ] T003 [P] 写入中文 `README.md`：项目简介（章程 XV 中文优先）、开发与构建说明、跨平台编译命令、引用 `quickstart.md`
- [ ] T004 [P] 写入 `LICENSE`：MIT 或 Apache-2.0（项目初始化阶段占位，PR 中确认最终许可）
- [ ] T005 [P] 写入 `.gitignore`：覆盖 `bin/`、`coverage.out`、`*.test`、`*.out`、本地构建产物与 IDE 临时文件

**Checkpoint**: `go mod download` 成功；目录布局与 plan.md §“Project Structure” 一致。

---

## Phase 2: Foundational（阻塞型前置）

**Purpose**: 所有 user story 共享的核心包；任何用户故事 MUST 在本阶段完成后才能开始。

**⚠️ CRITICAL**: 本阶段未完成前，任何 user story 工作不得启动。

- [ ] T006 实现版本信息包：在文件 `internal/version/version.go` 中声明 `Version`、`Commit`、`BuildTime`、`Platform`、`GoVersion` 五个字符串变量，默认值 `dev` / `unknown`（data-model.md “版本信息” 实体，决策 D-003，contract/章程 XI）
- [ ] T007 [P] 实现平台信息包：在文件 `internal/platform/platform.go` 中声明 `OS` / `Arch` / `IsUnixLike` / `IsWindows` 字段，运行时探测 `runtime.GOOS` / `runtime.GOARCH`
- [ ] T008 [P] 实现平台文件：在文件 `internal/platform/platform_unix.go` 中通过 `//go:build unix` 编译标签提供 Unix 实现
- [ ] T009 [P] 实现平台文件：在文件 `internal/platform/platform_windows.go` 中通过 `//go:build windows` 编译标签提供 Windows 实现
- [ ] T010 实现退出码常量：在文件 `internal/cli/exitcode.go` 中声明 `ExitSuccess=0`、`ExitGeneral=1`、`ExitUsage=2`、`ExitReservedStart=3` 常量集与 `FromError(err) int` 转换函数（决策 D-008，data-model.md “退出码” 实体）
- [ ] T011 [P] 实现中文文案集中模块：在文件 `internal/cli/messages.go` 中声明所有面向用户的中文字符串常量（`VersionTemplate`、`HelpIntro`、`ReservedMessage`、`ErrorTemplate` 等），使用 Go 原始字符串字面量（决策 D-007，章程 XV）
- [ ] T012 [P] 引入 Cobra 依赖：在文件 `go.mod` 中加入 `github.com/spf13/cobra`（决策 D-002）；`go mod tidy` 同步 `go.sum`
- [ ] T013 引入 golangci-lint 配置：在文件 `.golangci.yml` 中启用 `govet`、`errcheck`、`staticcheck`、`unused`、`gofmt`、`goimports` 六个 linter（决策 D-006）

**Checkpoint**: `go build ./...` 成功；`golangci-lint run ./...` 在空实现下零报错。

---

## Phase 3: User Story 1 — 基础 CLI 自检（Priority: P1）🎯 MVP

**Goal**: 用户可在终端执行 `aiaikit --version` 与 `aiaikit --help`，得到稳定的版本信息与帮助列表。

**Independent Test**: `go run ./cmd/aiaikit --version` 与 `go run ./cmd/aiaikit --help` 各自输出符合 `contracts/cli.md §1` 与 `§2` 的内容；退出码为 0。

### Implementation for User Story 1

- [ ] T014 [US1] 实现 CLI 入口：在文件 `cmd/aiaikit/main.go` 中调用 `internal/cli.Execute(stdout, stderr, exitFn)`，按退出码常量退出（章程 I / V 不污染 Core 入口）
- [ ] T015 [US1] 实现 root 命令：在文件 `internal/cli/root.go` 中通过 Cobra 注册 root 命令名 `aiaikit`，定义 `-v, --version` 与 `-h, --help` PersistentFlag，默认 `SilenceUsage=true` 与 `SilenceErrors=true`
- [ ] T016 [P] [US1] 实现版本子命令：在文件 `internal/cli/version.go` 中通过 Cobra 注册 `version` 子命令，使用 `internal/version` 包变量填充 `contracts/cli.md §1` 规定的五字段模板
- [ ] T017 [P] [US1] 实现 `--version` 标志：在文件 `internal/cli/root.go` 中通过 Cobra `PersistentFlags` 注册 `-v, --version`，与 `version` 子命令共用同一渲染函数
- [ ] T018 [US1] 实现 `--help` 行为：在文件 `internal/cli/root.go` 中设置 `Run` 字段为空时默认执行 help 行为（无子命令时输出帮助文本，退出码 0）
- [ ] T019 [US1] 写入版本 golden 文件：在文件 `testdata/cli/version.golden` 中写入占位快照（首行 `aiaikit 版本 dev`），供后续 golden 测试比对
- [ ] T020 [US1] 写入帮助 golden 文件：在文件 `testdata/cli/help.golden` 中写入完整的 `contracts/cli.md §2` 期望输出（含 `Available Commands` 八条与“（预留 / 未实现）”标注）
- [ ] T021 [US1] 实现 CLI 行为测试：在文件 `internal/cli/version_test.go` 中使用 `cmd.SetArgs` / `cmd.SetOut` / `cmd.Execute()` 模式断言 `aiaikit --version` 与 `aiaikit version` 输出匹配 `testdata/cli/version.golden`
- [ ] T022 [P] [US1] 实现 CLI 行为测试：在文件 `internal/cli/help_test.go` 中断言 `aiaikit --help` / `aiaikit`（无参）/ `aiaikit help` 三种入口输出匹配 `testdata/cli/help.golden`
- [ ] T023 [US1] 实现版本包单元测试：在文件 `internal/version/version_test.go` 中断言五个字段默认值为 `dev` / `unknown` / `runtime.GOOS + "/" + runtime.GOARCH` / `runtime.Version()`（data-model.md “版本信息” 不变式）

**Checkpoint**: `go run ./cmd/aiaikit --version` 与 `go run ./cmd/aiaikit --help` 在三类平台上输出一致；`go test ./...` 全通过。

---

## Phase 4: User Story 2 — 跨平台一致体验（Priority: P1）

**Goal**: 同源代码在 Windows / Linux / macOS 上构建产出一致的命令、选项与退出码。

**Independent Test**: CI 三平台 runner 各自 `make ci` 全通过；本地 `make build-all` 产出三类平台二进制，且每个产物 `aiaikit --version` 字段集合与顺序一致。

### Implementation for User Story 2

- [ ] T024 [US2] 实现 Makefile 基础目标：在文件 `Makefile` 中提供 `build`、`test`、`fmt`、`lint`、`cover`、`clean`、`ci`、`build-all` 八个目标；显式 `CGO_ENABLED=0`（quickstart.md §5，决策 D-004）
- [ ] T025 [US2] Makefile 注入版本信息：在文件 `Makefile` 的 `build` 目标中通过 `ldflags` 注入 `internal/version` 包的五个变量；接受 `VERSION / COMMIT / BUILD_TIME` 环境变量覆盖（决策 D-003）
- [ ] T026 [US2] Makefile 跨平台矩阵：在文件 `Makefile` 的 `build-all` 目标中遍历 `GOOS ∈ {linux, darwin, windows}` 与 `GOARCH ∈ {amd64, arm64}`，产出命名含 `GOOS-GOARCH` 的二进制
- [ ] T027 [US2] Makefile `ci` 目标：在文件 `Makefile` 中将 `ci` 串联为 `fmt + lint + test + build`，与本地开发者一致（决策 D-009）
- [ ] T028 [US2] 端到端测试：在文件 `e2e/version_test.go` 中通过 `go build` 产出的真实二进制调用 `--version`，断言退出码 0 且字段集合完整
- [ ] T029 [P] [US2] 端到端测试：在文件 `e2e/help_test.go` 中通过真实二进制调用 `--help`，断言八个 Available Commands 全部出现
- [ ] T030 [US2] 平台包单元测试：在文件 `internal/platform/platform_test.go` 中断言 `OS` / `Arch` / `IsUnixLike` / `IsWindows` 在当前构建目标下的取值符合预期

**Checkpoint**: 三类平台本地或 CI 构建成功；端到端测试在三类 runner 上全部通过；`make build-all` 产出 ≥ 6 个二进制（`linux/amd64`、`linux/arm64`、`darwin/amd64`、`darwin/arm64`、`windows/amd64`、`windows/arm64`）。

---

## Phase 5: User Story 3 — 可扩展命令与目录骨架（Priority: P2）

**Goal**: `init`、`update`、`runtime`、`plugin`、`web` 五个预留命令已注册到 root 命令下，`--help` 列出且各自输出“未实现”友好提示，退出码 2。

**Independent Test**: `aiaikit init` / `update` / `runtime` / `plugin` / `web` 各自输出 `contracts/cli.md §3` 的友好提示文案，退出码 2，且无网络请求、无磁盘写入、无 Spec-Kit Runtime 加载。

### Implementation for User Story 3

- [ ] T031 [P] [US3] 注册 `init` 预留命令：在文件 `internal/cli/reserved.go` 中通过 `RegisterReserved(root *cobra.Command, name, short string)` 注册 `init` 预留命令，描述含“（预留 / 未实现）”（章程 VI / XII 扩展点）
- [ ] T032 [P] [US3] 注册 `update` 预留命令：在文件 `internal/cli/reserved.go` 中注册 `update` 预留命令
- [ ] T033 [P] [US3] 注册 `runtime` 预留命令：在文件 `internal/cli/reserved.go` 中注册 `runtime` 预留命令
- [ ] T034 [P] [US3] 注册 `plugin` 预留命令：在文件 `internal/cli/reserved.go` 中注册 `plugin` 预留命令
- [ ] T035 [P] [US3] 注册 `web` 预留命令：在文件 `internal/cli/reserved.go` 中注册 `web` 预留命令
- [ ] T036 [US3] 预留命令统一行为：在文件 `internal/cli/reserved.go` 中定义共享 `runReserved(cmd *cobra.Command, args []string)`：输出 `contracts/cli.md §3` 模板，调用 `exitcode.ExitUsage` 退出（data-model.md “命令节点” 不变式）
- [ ] T037 [US3] 注册入口接入：在文件 `internal/cli/root.go` 的 `init()` 中调用 `RegisterReserved` 五次，确保自动出现在 `--help` 与 `Available Commands` 列表（顺序固定）
- [ ] T038 [US3] 注册预留命令机制：在 `internal/cli/reserved.go` 中标注 `Annotations["reserved"] = "true"`，供后续 Plugin 系统按注解挂载真实实现（章程 VI 扩展点）
- [ ] T039 [US3] 预留命令行为测试：在文件 `internal/cli/reserved_test.go` 中对五个预留命令各断言一次：输出文案、退出码 `ExitUsage=2`、无写入磁盘副作用
- [ ] T040 [P] [US3] 预留命令隔离测试：在文件 `internal/cli/reserved_isolation_test.go` 中通过 fsnotify 或临时目录 mock 验证预留命令不会修改任何路径下的文件（章程 XII / spec FR-004）

**Checkpoint**: `aiaikit --help` 列出 8 条命令（`completion`、`help`、`init`、`plugin`、`runtime`、`update`、`version`、`web`），预留命令描述以“（预留 / 未实现）”结尾；五个预留命令各自输出 §3 文案并以退出码 2 退出。

---

## Phase 6: User Story 4 — 基础测试与质量门禁（Priority: P2）

**Goal**: 仓库根目录一键执行格式化、静态检查、单元测试、CLI 行为测试、覆盖率输出，并在 CI 中通过；golden 文件契约稳定。

**Independent Test**: `make ci` 在 CI 三平台 runner 全通过；`go test -cover` 输出 `coverage.out`；`UPDATE_GOLDEN=1 make test` 仅在契约变更后用于更新 golden。

### Implementation for User Story 4

- [ ] T041 [US4] 格式化目标：在文件 `Makefile` 中实现 `fmt` 目标，调用 `gofmt -s -w` 与 `goimports -w`（quickstart.md §6.1）
- [ ] T042 [US4] 静态检查目标：在文件 `Makefile` 中实现 `lint` 目标，调用 `golangci-lint run ./...`（quickstart.md §6.2）
- [ ] T043 [US4] 测试目标：在文件 `Makefile` 中实现 `test` 目标，调用 `go test -race -count=1 ./...`；支持 `UPDATE_GOLDEN=1` 环境变量
- [ ] T044 [US4] 覆盖率目标：在文件 `Makefile` 中实现 `cover` 目标，调用 `go test -coverprofile=coverage.out`，生成 `coverage.html` 与摘要（quickstart.md §6.4）
- [ ] T045 [US4] golden 更新机制：在文件 `internal/cli/golden_test.go` 中实现 `assertGolden(t, name, got string)` 辅助函数：当 `UPDATE_GOLDEN=1` 时写入 `testdata/cli/<name>.golden`，否则做精确比对
- [ ] T046 [US4] golden 文件更新文档：在文件 `README.md` 与 `quickstart.md` 中补充“如何安全更新 golden 文件”的步骤（含 PR 描述模板）
- [ ] T047 [P] [US4] 静态检查豁免注释：在文件 `.golangci.yml` 中配置 `issues.exclude-default` 仅排除与本阶段无关的规则；避免后续阶段套用同一配置时引入豁免蔓延
- [ ] T048 [US4] 质量门禁冒烟测试：在文件 `internal/cli/quality_smoke_test.go` 中编写 `TestMakeCi` 占位测试：直接调用 `exec.Command("make", "ci")`，作为本地与 CI 一致性的最后一道闸（quickstart.md §6.5）

**Checkpoint**: `make ci` 本地与 CI 全绿；golden 文件覆盖 `--version` / `--help` / 五个预留命令共七个快照；覆盖率报告可用。

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: 跨用户故事的整理与一致性收尾。

- [ ] T049 [P] 文档：在文件 `README.md` 中补充章节 “架构分层”、`compliance`、`quality-gates`，与 `contracts/cli.md` / `data-model.md` 互相引用（章程 XIII）
- [ ] T050 [P] 文档：在文件 `specs/001-go-cli-bootstrap/quickstart.md` 中补充“CI 集成建议”章节，列出 GitHub Actions / GitLab CI / Jenkins 三类平台的最小 `make ci` 片段（占位，本阶段不强制启用 CI）
- [ ] T051 代码整理：在仓库根目录执行 `gofmt -s -w .` 与 `goimports -w .`，消除 Phase 2 ~ Phase 6 中累积的格式差异
- [ ] T052 章程一致性自检：在文件 `scripts/constitution-audit.sh`（占位脚本，可仅打印命令清单）中固化 quickstart.md §8 的三条自检命令（章程 I / VI / XV）
- [ ] T053 产物大小基线检查：在文件 `Makefile` 的 `build` 目标中追加 `du -h bin/aiaikit` 输出，并在 README 中记录当前实际值，便于后续修订章程 SC-024 的 30MB 基线
- [ ] T054 提交前自检：在文件 `CONTRIBUTING.md` 中固化“提交前必须执行 `make ci`”的开发约定（章程 XIV 治理可审计）
- [ ] T055 验收回归：执行 `quickstart.md §2 ~ §8` 全流程，验证 spec FR-001 ~ FR-010 全部通过

**Checkpoint**: `make ci` 三类平台全绿；`quickstart.md` 全流程可复现；README 与 contracts/cli.md 互相引用一致。

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖，可立即开始。
- **Foundational (Phase 2)**: 依赖 Setup 完成；**BLOCKS** 所有用户故事。
- **User Stories (Phase 3 ~ Phase 6)**: 全部依赖 Foundational 完成；P1 故事（US1 / US2）先行，P2 故事（US3 / US4）随后或与 P1 并行（仅当团队容量允许）。
- **Polish (Phase 7)**: 依赖所有目标用户故事完成。

### User Story Dependencies

- **US1（基础 CLI 自检）**: Foundational 完成后即可；不依赖其他故事。
- **US2（跨平台一致）**: Foundational 完成后即可；与 US1 在命令树实现上正交，可并行。
- **US3（可扩展骨架）**: 依赖 US1（root 命令已注册）才能挂载预留子命令；不依赖 US2。
- **US4（质量门禁）**: 依赖 US1 + US3 产生的测试目标与 golden 文件；不依赖 US2。

### Within Each User Story

- 测试任务可以与实现并行；本清单按 spec FR 顺序整理，不强求红绿循环。
- 包级常量与变量（T006 / T010 / T011）先于命令注册（T015 ~ T018 / T031 ~ T038）。
- 命令注册先于该命令的测试（T021 ~ T025 / T039 / T049）。
- 跨用户故事的依赖收敛在帮助 golden（T020）必须出现在 US3 测试（T039）之前——具体由 `internal/cli/help_test.go` 与 `internal/cli/reserved_test.go` 的执行顺序在 Tighten 内串行解决。

### Parallel Opportunities

- T002 / T003 / T004 / T005：Phase 1 中可并行（不同文件）。
- T007 / T008 / T009：Phase 2 中可并行（不同编译单元）。
- T016 / T017：US1 中可并行（`internal/cli/version.go` 与 `internal/cli/root.go` 的不同字段）。
- T021 / T022：US1 中可并行（不同测试文件）。
- T031 ~ T035：US3 中可并行（同一文件 `reserved.go`，但函数级别独立；实际实现时建议一次性写完以减少冲突）。
- T028 / T029：US2 端到端测试可并行。
- T049 / T050：Polish 阶段可并行。

---

## Parallel Example: User Story 1

```bash
# 并行启动 US1 中的版本包相关任务
Task: "T016 [P] [US1] 实现版本子命令 internal/cli/version.go"
Task: "T017 [P] [US1] 实现 --version 标志 internal/cli/root.go"
Task: "T023 [US1] 实现版本包单元测试 internal/version/version_test.go"
```

## Parallel Example: User Story 3

```bash
# 五个预留命令注册可在 reserved.go 中一次性完成；如需独立 PR，可拆分：
Task: "T031 [P] [US3] 注册 init 预留命令"
Task: "T032 [P] [US3] 注册 update 预留命令"
Task: "T033 [P] [US3] 注册 runtime 预留命令"
Task: "T034 [P] [US3] 注册 plugin 预留命令"
Task: "T035 [P] [US3] 注册 web 预留命令"
```

---

## Implementation Strategy

### MVP First（User Story 1 + User Story 2）

1. Phase 1: Setup（仓库骨架）
2. Phase 2: Foundational（共享包 + Cobra 接入 + lint 配置）
3. Phase 3: US1（`--version` / `--help`）
4. Phase 4: US2（跨平台 Makefile + e2e）
5. **STOP and VALIDATE**: `make ci` 本地通过；`aiaikit --version` / `--help` 在三类平台输出一致。
7. 部署 / 演示 MVP：可执行的跨平台 CLI 二进制，已通过质量门禁。

### Incremental Delivery

1. Setup + Foundational → 基础就绪。
2. 添加 US1 → 独立验证 → MVP 演示。
3. 添加 US2 → 三类平台验证 → 跨平台 MVP。
4. 添加 US3 → 五个预留命令占位 → 扩展点验证。
5. 添加 US4 → 完整质量门禁 → 生产就绪。
6. Polish → 文档、CI 建议、章程自检。

### Parallel Team Strategy（多开发者）

1. 团队共同完成 Setup + Foundational。
2. Foundational 完成后：
   - 开发者 A：US1（版本与帮助）
   - 开发者 B：US2（跨平台构建）
   - 开发者 C：US3（预留命令）
   - 开发者 D：US4（质量门禁）
3. 各故事独立完成、独立测试、独立合并。

---

## Notes

- `[P]` 任务在不同文件、无依赖时可并行；同文件多任务建议合并提交以减少冲突。
- `[Story]` 标签用于追溯到 spec.md 的具体用户故事。
- 每个用户故事独立可验证；任何 checkpoint 可暂停以独立演示。
- 不在 Core 入口（`cmd/aiaikit/main.go`、`internal/cli/root.go`）引入 Spec-Kit / Plugin / Web 任何业务实现；扩展通过 `RegisterReserved` 暴露的扩展点交付（章程 VI / XII）。
- 中文文案一律通过 `internal/cli/messages.go` 集中（决策 D-007 / 章程 XV）；不得在命令文件内散落字符串。