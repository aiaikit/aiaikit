<!-- Sync Impact Report
- Version change: 无 → 1.0.0（首次正式生效）
- Modified principles: 全部新增（项目初始章程）
- Added sections: 核心原则（15 条）、架构约束、治理
- Removed sections: 无
- Follow-up TODOs:
  - TODO(RATIFICATION_DATE): 若需追溯原始采用日期，补充到 Last Amended 之外的历史记录；当前按首次生效处理。
  - TODO(EXTENSION_HOOKS): 当前 .specify/extensions.yml 不存在，后续若引入扩展钩子须重新评估章程同步策略。
-->

# aiaikit Constitution
<!-- 项目：aiaikit —— 基于 Go 的 AI/SDD 工具生态平台 -->

本文档是 aiaikit 项目的最高治理准则，对 Spec、Plan、Tasks 及后续所有交付物具备强制约束力。任何与本文档冲突的项目实践 MUST 以本章程为准。

## Core Principles

### I. 不重复实现 Spec-Kit（Do-Not-Reinvent Spec-Kit）
aiaikit MUST 不重复实现 Spec-Kit 已经具备的能力。Spec-Kit 已实现的 SDD 工作流、模板、脚本、Agent Integration、斜杠命令等 MUST 保持原有行为；aiaikit 不得因复刻而引入行为偏移、行为修改或语义变更。理由：避免与上游 Spec-Kit 形成事实分叉，确保生态一致性与可审计性。

### II. Spec-Kit Runtime 是 aiaikit 的核心分发对象
aiaikit 主要职责 MUST 限定为 Spec-Kit Runtime 的分发、安装、版本管理、更新、插件管理与生态扩展。Runtime 相关的所有能力（包括但不限于模板、脚本、Integration、Agent 配置）一律通过 Runtime 包分发，禁止在 aiaikit CLI 中重新实现 Runtime 的业务逻辑。理由：单一职责，明确 aiaikit 的生态定位。

### III. `aiaikit init` 第一阶段兼容 `specify init`
`aiaikit init` 的第一阶段 MUST 保持与 `specify init` 的兼容行为（产出 `.specify/` 目录与对应斜杠命令）。后续阶段在此基础上叠加 aiaikit 扩展，禁止破坏 `specify init` 已建立的目录结构与文件契约。理由：保证现有 Spec-Kit 用户的迁移成本为零。

### IV. aiaikit 通过 npm 分发与升级
`aiaikit update` MUST 用于升级 aiaikit 自身，aiaikit 自身的分发与升级通道 MUST 走 npm。aiaikit CLI 的版本号、Release Notes、变更说明 MUST 与 npm 包保持一致；不得绕过 npm 直接发布二进制或源码。理由：统一分发入口，便于版本追溯与企业代理镜像。

### V. Spec-Kit Runtime 与 aiaikit CLI 解耦
Spec-Kit Runtime 与 aiaikit CLI MUST 解耦，独立维护各自的版本号、变更日志与兼容性矩阵。aiaikit MUST 支持 Runtime 的版本锁定、显式升级、跨大版本兼容性检查；不得因 aiaikit 升级而强制用户升级 Runtime，也不得因 Runtime 升级而强制 aiaikit 升级。理由：解耦带来独立的演进路径与回滚能力。

### VII. 插件系统是 aiaikit 的核心扩展机制
插件系统 MUST 作为 aiaikit 的核心扩展机制；后续 MUST 支持内置插件、第三方插件、Plugin Registry 与 Web 管理。Plugin 不得直接修改 aiaikit Core 的源码或运行时行为；扩展点 MUST 通过明确定义的 Plugin API 暴露。理由：保护 Core 的稳定性与可升级性，使生态可在 Core 之外演化。

### VIII. Spec-Kit Prompt Generator 作为首个内置插件
Spec-Kit Prompt Generator MUST 作为 aiaikit 的首个内置插件落地，其核心逻辑 MUST 不耦合到 aiaikit CLI。Prompt Generator 的输入、输出、版本与依赖 MUST 由插件自身声明并随插件升级，独立于 CLI 主版本。理由：把可变能力外置到插件层，验证 Plugin 体系的可行性。

### IX. 保留 Spec-Kit 原生斜杠命令与工作方式
aiaikit MUST 保留 Spec-Kit 原生的 `/constitution`、`/specify`、`/plan`、`/tasks`、`/implement` 等斜杠命令及其工作方式，不得改名、改协议、改入口。aiaikit 新增的命令 MUST 使用新的命名空间（例如 `/speckit-` 前缀），避免与 Spec-Kit 命令冲突。理由：保持用户已有的工作流可迁移、可复用。

### X. Web 管理中心复用 Core、Runtime 与 Plugin
Web 管理中心作为 aiaikit 的后续能力 MUST 与 CLI 共用 Go Core、Runtime 与 Plugin 系统，禁止重复实现业务逻辑。Web 的视图层、后端 API、权限模型 MUST 视作 Plugin API 的消费者，而非新的实现分支。理由：避免双栈双实现带来的维护成本与行为漂移。

### XI. 跨平台与安全基线
aiaikit MUST 支持 Windows、Linux、macOS 三类主流操作系统，且在三类系统上保持一致的 CLI、Runtime、Plugin 行为。跨平台优先级排序为：CLI 可维护性 > Runtime 可扩展性 > Plugin 安全性 > Web 可管理性。理由：保障多平台下用户与企业的部署一致性。

### XII. 资源版本、依赖与兼容性管理
Runtime、Plugin、Template、Integration 等所有可分发资源 MUST 具备明确的语义化版本、依赖声明与兼容性矩阵；aiaikit MUST 在安装、更新、加载阶段执行兼容性校验，不满足依赖或兼容性约束的资源 MUST 被拒绝安装或加载。理由：避免运行时静默失败与不可重现的崩溃。

### XIII. 新增能力优先通过扩展机制实现
所有新增能力 MUST 优先通过扩展机制（Plugin、Template、Integration、Preset）实现，禁止直接污染 aiaikit Core。Core 的变更仅在以下场景允许：(a) 新增明确的扩展点；(b) 修复阻塞性缺陷；(c) 重构以提升可扩展性，且必须附带 ADR 说明。理由：保证 Core 长期可演进。

### XIV. 全栈架构纳入范围
前端、后端、CLI、Runtime、Plugin、数据库、API、测试、安全以及 PC/移动端 UX MUST 纳入整体架构设计；不得仅关注后端或 CLI 而忽视其他层面。任何跨层面变更（例如 Plugin 增加 Web UI）MUST 同步更新对应的架构图与 ADR。理由：避免局部最优导致的系统级失衡。

### XV. 与 Spec-Kit 官方生态兼容优先
aiaikit MUST 优先保持与 Spec-Kit 官方生态的兼容，避免不必要的 fork 分叉。仅当官方缺乏变更且 aiaikit 有强需求时，方可创建本地扩展，并通过 Plugin / Preset / Integration 等可卸载机制交付，禁止直接修改上游 Runtime 的相关文件。理由：减少分叉负担与升级摩擦。

### XV. 中文优先（Chinese-First）
所有由本项目生成的文档、规范、计划、任务、ADR、README、代码注释、提交消息、Release Notes、错误消息、CLI 输出文案，以及 Spec-Kit 相关产物（`spec.md`、`plan.md`、`tasks.md`、`checklist.md`、`analysis.md`、本 `constitution.md`、插件文档、Integration 文档、Slash Command 文档、Preset 说明等），默认 MUST 使用简体中文。技术术语、API/CLI/代码标识符、错误码、版本号、URL、命令路径、库名（如 npm、Spec-Kit、Go、Plugin）等 MUST 保持原始形式；仅在英文为通用惯例且与上游生态对齐时允许保留英文（例如 Spec-Kit 官方命令 `/constitution` 等）。禁止机翻痕迹明显的文案；中英文混排时 MUST 保留必要空格与标点规范。理由：保证团队内审、跨区域协作与中文用户的可读性、可审计性，降低语义损耗。

## Architecture Constraints

aiaikit MUST 沿以下架构分层组织代码与发布物，任何新增模块必须明确其所在层与跨层契约：

1. **Go Core**：CLI 入口、版本管理、安装/升级流程、扩展点注册、Plugin Loader；不得包含 Runtime 业务实现。
2. **Spec-Kit Runtime**：以版本化包形式由 aiaikit 分发；Runtime 内部变更由 Spec-Kit 上游治理，aiaikit 不直接修改。
3. **Plugin Runtime**：Plugin 的加载、隔离、生命周期管理；Plugin API 必须在 Go Core 中以稳定签名定义。
4. **Plugin 实现**：内置插件（如 Prompt Generator）与第三方插件必须遵循同一份 Plugin Manifest 与依赖声明。
5. **Template / Integration / Preset**：作为可分发资源随 Runtime 或独立包发布；具备版本与兼容性声明。
6. **Web 管理中心**：复用 Go Core 暴露的 Plugin API；Web 视图层可独立选型，但禁止重新实现 Runtime 业务。
7. **跨平台适配层**：封装 Windows/Linux/macOS 的文件系统、进程、可执行位、权限差异；上层一律以 Json 标准接口访问。

新增能力 MUST 落点定位清晰：若能在 Plugin 层实现，禁止下沉到 Core；若能在 Preset/Template 层实现，禁止新增 Plugin。

## Governance

本章程是 aiaikit 项目治理的最高准则。任何 Spec、Plan、Tasks、ADR、PR 描述、Release Notes MUST 显式或隐式遵循本章程；与本章程冲突的项目实践、模板、自动化脚本一律以本章程为准。

- **修订程序**：本章程的修订必须通过 `speckit-constitution` 命令（或等价工作流）执行；修订 MUST 在 PR 中提供 Sync Impact Report（版本变化、原则增删、推迟项）。
- **版本策略**：本章程遵循 `MAJOR.MINOR.PATCH` 语义化版本：
  - MAJOR：删除或重新定义任一现有原则，或引入向后不兼容的治理规则；
  - MINOR：新增原则、章节，或对现有原则做实质性扩展；
  - PATCH：措辞澄清、笔误修复、非语义性润色。
- **合规审查**：所有 PR MUST 在描述中确认对章程的合规性；若涉及 Core 变更，必须附 ADR 说明动机与扩展点影响。
- **推迟项**：章程中标注 `TODO(<FIELD_NAME>)` 的条目 MUST 在后续 PR 中补齐，否则每三个月复查一次直至清空。

**Version**: 1.0.0 | **Ratified**: 2026-09-28 | **Last Amended**: 2026-09-28