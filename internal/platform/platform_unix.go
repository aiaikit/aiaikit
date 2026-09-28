//go:build unix

package platform

// isUnixLike 在 Unix 类系统上返回 true。
//
// Unix 类系统包括 linux / darwin / freebsd / openbsd / netbsd 等。
func isUnixLike() bool {
	return true
}
