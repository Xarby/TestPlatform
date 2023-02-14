package Util

import (
	"strings"
)

func ChangeRole(role string) string {
	if strings.Contains(role,"master"){
		return "master"
	}else if strings.Contains(role,"slave") {
		return "slave"
	}else {
		return "backup"
	}
}
