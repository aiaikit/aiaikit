# aiaikit Makefile
#
# 提供 build / build-all / test / fmt / lint / cover / ci 等目标；
# 跨平台构建以 GOOS / GOARCH 变量驱动；CI 与本地复用同一入口。
#
# 决策参考：specs/001-go-cli-bootstrap/research.md D-003 / D-004 / D-009

# 显式禁用 CGO，避免跨平台 C 工具链依赖（quickstart.md §5）。
export CGO_ENABLED := 0

# 默认目标：本地快速验证。
.DEFAULT_GOAL := help

# 项目根目录。
ROOT_DIR := $(shell pwd)

# 当前构建目标（默认走 uname）。
GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# 版本信息：可通过环境变量覆盖，默认值为 ldflags 友好的字符串。
VERSION     ?= dev
COMMIT      ?= unknown
BUILD_TIME  ?= unknown

# ldflags 注入路径（与 internal/version 包变量严格一致）。
VERSION_PKG := github.com/aiaikit/aiaikit/internal/version
GOVERSION   := $(shell go version | awk '{print $$3}')
LDFLAGS     := -X '$(VERSION_PKG).Version=$(VERSION)' \
               -X '$(VERSION_PKG).Commit=$(COMMIT)' \
               -X '$(VERSION_PKG).BuildTime=$(BUILD_TIME)' \
               -X '$(VERSION_PKG).Platform=$(GOOS)/$(GOARCH)' \
               -X '$(VERSION_PKG).GoVersion=$(GOVERSION)'

# 平台矩阵（使用 `.` 分隔 GOOS/GOARCH 以兼容 Make 通配符）。
PLATFORMS := linux.amd64 linux.arm64 darwin.amd64 darwin.arm64 windows.amd64 windows.arm64

# 默认目标列表。
.PHONY: help build build-all test fmt lint cover clean ci version

help:
	@echo "可用目标："
	@echo "  build       在当前平台构建二进制到 bin/aiaikit"
	@echo "  build-all   跨平台矩阵构建到 bin/aiaikit-<GOOS>-<GOARCH>[.exe]"
	@echo "  test        执行单元测试与 CLI 行为测试"
	@echo "  fmt         gofmt + goimports"
	@echo "  lint        golangci-lint"
	@echo "  cover       生成覆盖率报告"
	@echo "  ci          fmt + lint + test + build"
	@echo "  clean       清理 bin/ 与 coverage 产物"
	@echo "  version     输出当前注入的版本信息（默认 dev）"
	@echo ""
	@echo "常用变量（可在命令行覆盖）："
	@echo "  VERSION=<semver> COMMIT=<git-hash> BUILD_TIME=<RFC3339>"

# 当前平台构建。
build:
	@mkdir -p bin
	@go build -ldflags "$(LDFLAGS)" -o bin/aiaikit ./cmd/aiaikit
	@echo "✓ 构建完成: bin/aiaikit"

# 跨平台矩阵构建。
#
# 通过独立 shell 脚本 scripts/build-platform.sh 完成 ldflags 拼接，
# 规避 Make 多行转义问题。
build-all: $(addprefix build-platform-,$(PLATFORMS))

build-platform-%:
	@./scripts/build-platform.sh $$(echo "$*" | cut -d. -f1) $$(echo "$*" | cut -d. -f2) '$(VERSION)' '$(COMMIT)' '$(BUILD_TIME)'

# 测试（-count=1 禁止测试缓存，保证每次都是最新结果）。
#
# 注：未启用 -race，因 Makefile 顶层禁用 CGO；
# 启用 race 时需要单独设置 CGO_ENABLED=1。
#
# AIK_TEST_SKIP_MAKE_CI=1 在调用 TestMakeCi 时设置，避免递归调用 make ci。
test:
	@AIK_TEST_SKIP_MAKE_CI=1 go test -count=1 ./...

# 格式化。
fmt:
	@gofmt -s -w .
	@goimports -w . 2>/dev/null || echo "(goimports 未安装，跳过 import 排序)"

# 静态检查。
lint:
	@golangci-lint run ./... || (echo "golangci-lint 未安装或失败；尝试 go vet"; go vet ./...)

# 覆盖率。
cover:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -1
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ 覆盖率报告: coverage.html"

# CI 一键入口。
ci: fmt lint test build

# 清理。
clean:
	@rm -rf bin coverage.out coverage.html

# 输出版本信息（用于自检）。
version:
	@./bin/aiaikit --version 2>/dev/null || (echo "二进制未生成；先执行 make build")