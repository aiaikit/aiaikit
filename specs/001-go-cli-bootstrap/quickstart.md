# Quickstart：aiaikit 001 — 本地构建与验证

本指南给出从仓库根目录开始的可复现步骤，用于验证本阶段交付物（CLI Core 骨架、`--version`、`--help`、预留命令、质量门禁）。

> 适用平台：Windows（含 Git Bash / WSL）、Linux、macOS。
> 前置依赖：Go 工具链（具体最低版本以 `go.mod` 中 `go` 指令为准）、GNU Make 3.81+（Windows 上 Git for Windows 自带）。
> 推荐：CI runner 与本地开发者使用同一 Go 主版本，避免“本地过 / CI 红”。

---

## 1. 仓库根目录布局（验证项）

执行：

```bash
ls -1
```

期望可见（按章程“Architecture Constraints”分层）：

```
.golangci.yml
LICENSE
Makefile
README.md
cmd/
go.mod
go.sum
internal/
testdata/
```

不期望出现：`vendor/`、`runtime/`、`plugins/`、`web/`、`node_modules/`（本阶段均未引入）。

---

## 2. 开发模式快速验证（无须构建）

```bash
go run ./cmd/aiaikit --version
go run ./cmd/aiaikit --help
go run ./cmd/aiaikit init
```

期望：

- `--version` 输出 `aiaikit 版本 dev` 等五行字段（参见 `contracts/cli.md §1`）；
- `--help` 输出 `contracts/cli.md §2` 中定义的 Available Commands 列表；
- `init` 提示 “命令 "init" 尚未实现，将在后续阶段提供”，退出码 `2`。

---

## 3. 标准本地构建

```bash
make build
./bin/aiaikit --version
./bin/aiaikit --help
```

期望：二进制生成于 `./bin/aiaikit`，版本输出同上。

---

## 4. 注入真实版本信息的构建

```bash
VERSION="0.1.0"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

make build VERSION="$VERSION" COMMIT="$COMMIT" BUILD_TIME="$BUILD_TIME"
./bin/aiaikit --version
```

期望：版本输出对应字段被替换为真实值；`Commit` 与 `构建时间` 不再是 `unknown`。

> Windows PowerShell 等价：

```powershell
$VERSION="0.1.0"
$COMMIT=(git rev-parse --short HEAD 2>$null); if (-not $COMMIT) { $COMMIT="unknown" }
$BUILD_TIME=(Get-Date -AsUTC -Format "yyyy-MM-ddTHH:mm:ssZ")

make build VERSION=$VERSION COMMIT=$COMMIT BUILD_TIME=$BUILD_TIME
.\bin\aiaikit.exe --version
```

---

## 5. 跨平台矩阵构建（本地可选，CI 必跑）

```bash
make build-all
ls -1 bin/
```

期望（具体命名以 Makefile 实际产出的清单为准）：

```
aiaikit-darwin-amd64
aiaikit-darwin-arm64
aiaikit-linux-amd64
aiaikit-linux-arm64
aiaikit-windows-amd64.exe
aiaikit-windows-arm64.exe
```

在每个产物上执行 `--version` 应均可运行（Windows 上需在目标 Windows 机器 / WSL 中执行；本机仅能验证当前平台产物）。

---

## 6. 质量门禁

### 6.1 格式化

```bash
make fmt
```

期望：无输出（所有文件已符合格式）；非零退出码意味着存在需要修复的文件。

### 6.2 静态检查

```bash
make lint
```

期望：所有 linter 通过；首次运行可能在 Windows 上因 `golangci-lint` 安装位置差异需要预装到 PATH。

### 6.3 单元测试与 CLI 行为测试

```bash
make test
```

期望：所有测试通过；输出包含 `ok` 字样。

更新 golden 文件（仅在契约变更后）：

```bash
UPDATE_GOLDEN=1 make test
```

### 6.4 覆盖率

```bash
make cover
```

期望：生成 `coverage.out` 与覆盖率摘要；覆盖率目标值由后续阶段决定，本阶段仅生成报告。

### 6.5 一键 CI

```bash
make ci
```

等价于 `fmt + lint + test + build`，CI 与本地均使用此入口。

---

## 7. 行为契约验证（按 contracts/cli.md）

### 7.1 `--version`

```bash
./bin/aiaikit --version
```

期望字段集合与顺序与 `contracts/cli.md §1` 一致；快照对照：

```bash
diff <(./bin/aiaikit --version) <(sed 's/<[^>]*>/.../g' testdata/cli/version.golden) && echo OK
```

注：snapshot 比对仅用于人工冒烟；正式断言由 `make test` 中 golden 测试承担。

### 7.2 `--help`

```bash
./bin/aiaikit --help
```

期望：八个 Available Commands 全部出现，预留命令描述以“（预留 / 未实现）”结尾。

### 7.3 预留命令

```bash
./bin/aiaikit init; echo "exit=$?"
./bin/aiaikit update; echo "exit=$?"
./bin/aiaikit runtime; echo "exit=$?"
./bin/aiaikit plugin; echo "exit=$?"
./bin/aiaikit web; echo "exit=$?"
```

期望：

- stdout 出现友好提示文案；
- 退出码 `exit=2`；
- 没有网络请求、文件下载、Spec-Kit Runtime 加载。

---

## 8. 章程一致性自检

```bash
# 章程 I / V：不应引用 Spec-Kit 任何模块
grep -r "spec-kit\|speckit\|github.com/specify" go.mod go.sum || echo "OK: 无 Spec-Kit 依赖"

# 章程 VI / XII：扩展点存在
test -f internal/cli/reserved.go && echo "OK: 预留命令扩展点存在"

# 章程 XV：CLI 文案为中文
grep -l "aiaikit" internal/cli/messages.go >/dev/null && echo "OK: 中文文案集中"
```

期望三项均通过。

---

## 9. 故障排查

| 现象 | 可能原因 | 处置 |
|---|---|---|
| `--version` 输出全 `unknown` | Make 注入失败 / 变量未导出 | 检查 Makefile 中的 ldflags；PowerShell 下确认 `$VAR` 传递 |
| `--help` 顺序漂移 | Cobra 版本升级 | 锁定 `github.com/spf13/cobra` 版本；更新 golden 文件并修订本契约 |
| 静态检查失败：`errcheck` 报错 | 未处理 `cmd.Execute()` 错误 | 在 `main.go` 中通过 `exitcode.FromError` 处理 |
| Windows 上 `make` 不可用 | 未安装 Git for Windows / WSL | 安装 Git for Windows；或在 PowerShell 中使用等价命令 |
| 跨平台构建失败：`CGO_ENABLED` 相关 | 默认开启 CGO | Makefile 中显式 `CGO_ENABLED=0` |
| `--version` 输出乱码 | 终端编码非 UTF-8 | 设置终端编码为 UTF-8（Windows: `chcp 65001`） |

---

## 10. 下一步

- 本阶段交付后，进入 `/speckit-tasks` 生成实现任务清单；
- 后续阶段：`/speckit-specify` 启动 Runtime / Plugin / Web 子模块的独立 spec。