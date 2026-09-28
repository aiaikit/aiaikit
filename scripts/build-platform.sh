#!/usr/bin/env bash
#
# build-platform.sh —— 单平台构建脚本，供 Makefile 调用。
#
# 用法：
#   ./scripts/build-platform.sh <GOOS> <GOARCH> [<VERSION> <COMMIT> <BUILD_TIME>]
#
# 输出：
#   bin/aiaikit-<GOOS>-<GOARCH>[.exe]
#
# 环境变量（可覆盖默认值）：
#   VERSION, COMMIT, BUILD_TIME, GOVERSION

set -eu

GOOS="${1:?缺少 GOOS}"
GOARCH="${2:?缺少 GOARCH}"
VERSION="${3:-dev}"
COMMIT="${4:-unknown}"
BUILD_TIME="${5:-unknown}"
GOVERSION="${GOVERSION:-$(go version | awk '{print $3}')}"

VERSION_PKG="github.com/aiaikit/aiaikit/internal/version"
SUFFIX=""
if [ "$GOOS" = "windows" ]; then
  SUFFIX=".exe"
fi

mkdir -p bin

OUTPUT="bin/aiaikit-${GOOS}-${GOARCH}${SUFFIX}"

# 使用单一引号包裹每个 -X，再以双引号包裹整个 ldflags 字符串。
# 平台字段需要运行时拼接到单引号块中（Platform=<GOOS>/<GOARCH>）。
LDFLAGS="-X '${VERSION_PKG}.Version=${VERSION}' \
         -X '${VERSION_PKG}.Commit=${COMMIT}' \
         -X '${VERSION_PKG}.BuildTime=${BUILD_TIME}' \
         -X '${VERSION_PKG}.Platform=${GOOS}/${GOARCH}' \
         -X '${VERSION_PKG}.GoVersion=${GOVERSION}'"

echo "→ 构建 ${GOOS}/${GOARCH} -> ${OUTPUT}"
GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags "$LDFLAGS" -o "$OUTPUT" ./cmd/aiaikit
ls -lh "$OUTPUT"