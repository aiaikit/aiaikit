//go:build windows

package platform

// isUnixLike 在 Windows 上返回 false。
func isUnixLike() bool {
	return false
}
