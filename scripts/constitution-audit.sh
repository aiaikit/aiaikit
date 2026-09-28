#!/usr/bin/env bash
#
# constitution-audit.sh —— 一键执行章程一致性自检。
#
# 章程参考：.specify/memory/constitution.md
# 本脚本实现 quickstart.md §8 中的三条核心自检命令，
# 并扩展覆盖章程 X（跨平台）与 XV（中文优先）的常规检查。

set -eu

# 颜色（仅在 TTY 时启用）。
if [ -t 1 ]; then
  RED=$'\e[31m'; GREEN=$'\e[32m'; YELLOW=$'\e[33m'; RESET=$'\e[0m'
else
  RED=""; GREEN=""; YELLOW=""; RESET=""
fi

ok()   { printf "%s✓%s %s\n" "$GREEN" "$RESET" "$1"; }
warn() { printf "%s!%s %s\n" "$YELLOW" "$RESET" "$1"; }
fail() { printf "%s✗%s %s\n" "$RED" "$RESET" "$1"; }

ROOT="$(git rev-parse --show-toplevel)"
cd "$ROOT"

echo "=== 章程一致性自检 ==="
echo

# --- 章程 I / V：不应引用 Spec-Kit 任何模块 ---
echo "[章程 I / V] 不重复 / 不耦合 Spec-Kit"
if grep -rE 'spec-kit|speckit|github.com/specify' go.mod go.sum 2>/dev/null; then
  fail "go.mod / go.sum 包含 Spec-Kit 依赖（违反章程 I / V）"
else
  ok "go.mod / go.sum 无 Spec-Kit 依赖"
fi
echo

# --- 章程 VI / XII：扩展点存在 ---
echo "[章程 VI / XII] 扩展点与预留命令"
if [ -f internal/cli/reserved.go ]; then
  ok "internal/cli/reserved.go 存在"
else
  fail "缺少 internal/cli/reserved.go（违反章程 VI）"
fi
if grep -q 'ReservedAnnotationKey' internal/cli/*.go; then
  ok "预留命令 Annotation 扩展点已标注"
else
  warn "未发现 ReservedAnnotationKey 标注"
fi
echo

# --- 章程 X：跨平台构建矩阵 ---
echo "[章程 X] 跨平台矩阵"
PLATFORMS=("linux-amd64" "linux-arm64" "darwin-amd64" "darwin-arm64" "windows-amd64.exe" "windows-arm64.exe")
for p in "${PLATFORMS[@]}"; do
  f="bin/aiaikit-${p}"
  if [ -f "$f" ]; then
    ok "$f 已构建"
  else
    warn "$f 未构建（执行 make build-all）"
  fi
done
echo

# --- 章程 XV：中文文案集中 ---
echo "[章程 XV] 中文优先"
if [ -f internal/cli/messages.go ]; then
  ok "internal/cli/messages.go 集中中文文案"
else
  fail "缺少 internal/cli/messages.go（违反章程 XV）"
fi
# 抽查二进制输出是否包含中文
if [ -x ./bin/aiaikit ]; then
  if ./bin/aiaikit --help 2>&1 | grep -q "初始化 aiaikit 工作区"; then
    ok "--help 输出包含中文章节"
  else
    fail "--help 输出缺少中文（违反章程 XV）"
  fi
  if ./bin/aiaikit init 2>&1 | grep -q "尚未实现"; then
    ok "init 提示包含中文"
  else
    fail "init 提示缺少中文"
  fi
  true
else
  warn "./bin/aiaikit 不存在，跳过二进制抽查（执行 make build）"
fi

echo
echo "=== 自检完成 ==="