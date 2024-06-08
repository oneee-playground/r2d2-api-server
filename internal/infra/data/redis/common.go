package redis

import "strings"

func buildKey(strs ...string) string {
	return strings.Join(strs, ":")
}
