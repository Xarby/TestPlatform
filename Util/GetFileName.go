package Util

import (
	"strings"
)

func GetFileName(Path string) (string) {
	if strings.Contains(Path, "/") {
		return Path[strings.LastIndex(Path, "/")+1:]
	} else if strings.Contains(Path, "\\") {
		return Path[strings.LastIndex(Path, "\\")+1:]
	} else {
		return ""
	}
}