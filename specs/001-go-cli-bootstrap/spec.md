# Feature Specification: aiaikit 001 — Go 项目初始化与 CLI 基础框架

**Feature Branch**: `001-go-cli-bootstrap`

**Created**: 2026-09-28

**Status**: Draft

**Input**: User description: "@.aiaikit/prompts/001/spec.md"

## User Scenarios & Testing *(mandatory)*

### User Story 1 — 基础 CLI 自检（Priority: P1）

作为 aiaikit 用户或贡献者，我期望在终端执行 `aiaikit --version` 与 `aiaikit --help`，从而在数秒内确认 CLI 已正确安装、版本号可读、可用命令列表完整。

**Why this priority**: 这是 CLI 可被识别的最小可验证信号；缺少它则后续所有命令（init、update、runtime、plugin、web）都无法被用户或脚本正确调用，也难以在 CI/CD 中嵌入自检。

**Independent Test**: 在 Windows / Linux / macOS 任一平台上克隆源码后执行标准构建脚本，再运行 `aiaikit --version` 与 `aiaikit --help`，可独立验证版本号、命令清单、退出码。

**Acceptance Scenarios**:

1. **Given** 已通过构建脚本生成本平台二进制，**When** 用户执行 `aiaikit --version`，**Then** 输出包含语义化版本号、构建时间（Git Commit Hash 可选）、退出码为 0。
2. **Given** 已通过构建脚本生成本平台二进制，**When** 用户执行 `aiaikit --help`，**Then** 输出包含本阶段已实现与预留的命令列表、每个命令的简短说明，并明确标识哪些命令为预留（尚未实现）。
3. **Given** 任意用户调用未实现的预留子命令（例如 `aiaikit init`），**Then** CLI 给出友好提示，指明该命令将在后续阶段实现并不退出非零码之外的副作用。

---

### User Story 2 — 跨平台一致体验（Priority: P1）

作为在 Windows / Linux / macOS 上工作的开发者或运维，我希望同一份源码构建出来的 `aiaikit` 在三类系统中输出一致的命令、选项与退出码，从而可以在不同环境中无缝复用脚本与文档。

**Why this priority**: 章程 XI 要求跨平台一致；没有跨平台验证则 CLI Core 的可移植性无法被信任，后续 Runtime / Plugin 分发将沿用同一通道。

**Independent Test**: 在三个目标平台各自执行同一组 CLI 命令，比对版本输出格式、帮助输出结构、退出码。

**Acceptance Scenarios**:

1. **Given** 已分别在 Windows / Linux / macOS 上完成构建，**When** 在三个平台执行相同的 `aiaikit --version` 与 `aiaikit --help`，**Then** 输出结构、字段顺序、退出码一致。
2. **Given** 任意目标平台，**When** 用户在仓库根目录执行标准构建命令，**Then** 构建成功并产出可在该平台直接运行的二进制。
3. **Given** CI 上同时跑三类平台的构建任务，**When** 任一平台构建失败，**Then** CI 报告明确的失败平台与失败阶段，便于定位。

---

### User Story 3 — 可扩展的命令与目录骨架（Priority: P2）

作为 aiaikit 后续阶段的开发者，我希望 CLI 已经为 `init`、`update`、`runtime`、`plugin`、`web` 等命令预留入口与目录结构，使后续业务实现只需在已划定的位置增量填充，而无需改动 Core 入口。

**Why this priority**: 章程 II / VII / X 强调 Core / Runtime / Plugin / Web 的职责分离；提前把目录与命令骨架定好能显著降低后续返工与分叉风险。

**Independent Test**: 在源码中确认预留命令的注册点存在、对应目录或包已创建、单元测试或 CLI 测试可以枚举到这些命令。

**Acceptance Scenarios**:

1. **Given** 开发者查看 CLI 源码，**When** 搜索 `init`、`update`、`runtime`、`plugin`、`web` 关键字，**Then** 可定位到对应的命令注册入口（命令可未实现业务）。
2. **Given** 开发者新增一个预留命令，**When** 不修改 Core 入口协议即可完成注册，**Then** 命令自动出现在 `aiaikit --help` 中并以“预留”状态展示。

---

### User Story 4 — 基础测试与质量门禁（Priority: P2）

作为 aiaikit 维护者，我希望从一开始就有可执行的单元测试、CLI 行为测试、格式化与静态检查工具链，使得后续 PR 都跑在统一的质量门禁之上。

**Why this priority**: 没有质量门禁就难以防止 Core 退化；章程 XIII / XIV 要求 Core 变更必须附带说明，自动化检查是低成本兑现途径。

**Independent Test**: 在仓库根目录执行标准命令即可触发格式化、静态检查、单元测试、CLI 行为测试，并全部通过。

**Acceptance Scenarios**:

1. **Given** 仓库根目录，**When** 执行标准格式化命令，**Then** 代码被自动规范化为统一风格。
2. **Given** 仓库根目录，**When** 执行静态检查命令，**Then** 未通过检查的代码在 CI 中被拦截。
3. **Given** 仓库根目录，**When** 执行单元测试与 CLI 测试，**Then** 所有测试通过且报告可读。

---

### Edge Cases

- 在不支持的操作系统或 CPU 架构上执行构建：构建脚本给出明确错误信息，不静默产生错误二进制。
- 多次重复执行构建：幂等，不残留临时文件污染源码树。
- 用户调用未实现的子命令并附带错误参数：CLI 输出 usage 提示与退出码，但不在磁盘上写入任何文件、不下载任何资源。
- 仓库缺少构建工具链（Go、格式化工具、静态检查工具）：构建脚本给出明确的依赖缺失提示。
- 版本号或构建信息缺失：构建产物仍可运行，但 `aiaikit --version` 显式标注 “dev / unknown” 而不是崩溃。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST 提供 `aiaikit --version` 子命令，输出当前语义化版本号、构建目标平台与架构、构建时间（可选 Git Commit Hash）。
- **FR-002**: CLI MUST 提供 `aiaikit --help`（与默认无参调用）子命令，输出所有已实现与已预留的命令清单，每个命令附带简短中文说明。
- **FR-003**: CLI MUST 预留 `init`、`update`、`runtime`、`plugin`、`web` 五个顶层子命令入口；本阶段 MUST 不实现其业务逻辑，但 MUST 在 `--help` 中明确标注“预留 / 未实现”。
- **FR-004**: 调用本阶段未实现的预留子命令 MUST 输出友好提示（中文文案），退出码符合 CLI 通用约定（参数错误类提示使用非零退出码），且 MUST NOT 下载、内置或执行任何 Spec-Kit Runtime 代码。
- **FR-005**: 源码目录 MUST 按章程 “Architecture Constraints” 划定的层级组织，包括 CLI 入口、版本信息、命令注册、跨平台适配占位、本阶段单元与 CLI 测试基目录；分层 MUST 保持单向依赖，不得反向依赖 Runtime 或 Plugin。
- **FR-006**: 项目 MUST 提供标准构建脚本，能在 Windows、Linux、macOS 各自的标准工具链下产出可执行二进制，并生成版本信息注入产物。
- **FR-007**: 项目 MUST 提供标准格式化、静态检查、单元测试、CLI 行为测试命令，开发者与 CI 可一键调用。
- **FR-008**: 本阶段 MUST NOT 实现任何 Spec-Kit SDD 工作流（`/constitution`、`/specify`、`/plan`、`/tasks`、`/implement` 等）；不得下载、内置、引用或执行 Spec-Kit Runtime 任何部分。
- **FR-009**: 本阶段 MUST NOT 修改或覆盖 Spec-Kit 已发布命令协议；后续阶段引入 aiaikit 专属命令时 MUST 使用与 Spec-Kit 不同的命名空间。
- **FR-010**: CLI 错误信息、日志、MUST 使用中文；版本号、命令名、路径、URL 等标识符保留原始形式。

### Key Entities

- **CLI 入口（Entry）**: 表示用户与 aiaikit CLI 交互的总入口；关键属性为二进制名（`aiaikit`）、支持的平台与架构、注入的版本信息。
- **命令节点（Command Node）**: 表示一个顶层子命令（`init` / `update` / `runtime` / `plugin` / `web` / ...）；关键属性为命令名、当前阶段状态（`已实现` / `预留未实现`）、中文说明、所属分层。
- **版本信息（Version Info）**: 表示一次构建注入的元数据；关键属性为语义化版本号、构建时间、Git Commit Hash、构建目标、构建者环境。
- **构建产物（Build Artifact）**: 表示一次构建产出的二进制；关键属性为目标平台与架构、可执行路径、是否经过质量门禁。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 开发者或 CI 在 Windows、Linux、macOS 任一平台执行标准构建命令，从仓库根目录到产出可运行二进制所需时间不超过 5 分钟（冷启动场景除外）。
- **SC-002**: `aiaikit --version` 在三类平台上的输出结构（字段集合与顺序）保持一致，且退出码为 0。
- **SC-003**: `aiaikit --help` 输出的预留命令数不少于 5 个（`init`、`update`、`runtime`、`plugin`、`web`），每个预留命令 MUST 标注“未实现 / 预留”状态。
- **SC-004**: 仓库根目录一键执行格式化、静态检查、单元测试、CLI 行为测试，CI 流水线中四项门禁全部通过率 100%。
- **SC-005**: 调用任意预留子命令不会触发网络请求、文件下载、Spec-Kit Runtime 加载；可通过网络抓包与文件系统快照验证。
- **SC-006**: 本阶段产出的二进制文件总大小不超过章程“Architecture Constraints”中同类型 Go CLI 通用基线（参考值 ≤ 30MB；具体数值以后续 Architecture Constraints 修订为准）。

## Assumptions

- 目标用户已安装 Go 工具链（具体最低版本以 Go 官方支持矩阵为准），并具备基本的命令行使用经验。
- 本阶段不引入 Spec-Kit Runtime 任何代码、模板、脚本，因此不依赖 Spec-Kit 仓库结构、`.specify/` 目录、`/constitution` 等斜杠命令协议。
- 本阶段不引入 Plugin 加载机制；Plugin 系统的实现推迟到后续阶段，但目录与命令占位保留以便后续增量接入。
- 本阶段不引入 Web 后端或前端代码；`web` 子命令仅为目录与命令占位。
- CI 环境提供 Windows、Linux、macOS 三类 runner 或等价容器镜像；构建脚本不假设具体 CI 平台，但提供标准本地复现命令。
- 跨平台适配层仅做占位（文件系统路径分隔、进程执行差异等），不引入重量级依赖。

## Out of Scope

- 任何 Spec-Kit Runtime 相关行为（安装、加载、模板解析、斜杠命令注册）。
- Plugin 加载、隔离、热更新机制；仅保留命令与目录占位。
- Web 管理中心的任何 UI、后端 API、鉴权。
- 自动更新（`aiaikit update`）的真实网络下载流程。
- 与 npm、GitHub、企业代理镜像的发布链路集成（属于后续阶段）。
- 任何与 AI 模型调用、Prompt 生成、SDD 工作流相关的功能。