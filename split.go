//go:build !windows
// +build !windows

package hapesay

import "strings"

func splitPath(s string) []string {
	return strings.Split(s, ":")
}
