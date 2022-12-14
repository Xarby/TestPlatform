package Util

import (
	"TestPlatform/Const"
	"strings"
)

func ChooseVersion(version string) (string, string) {
	if strings.Contains(version, "3.10") || strings.Contains(version, "3.11") {
		return Const.ZddiPrivateMenuName + "zddi-private1.tar.gz", Const.ZddiRemoteFilePath + "zddi-private1.tar.gz"

	} else if strings.Contains(version, "3.13") || strings.Contains(version, "3.14") || strings.Contains(version, "3.15") {
		return Const.ZddiPrivateMenuName + "zddi-private2.tar.gz", Const.ZddiRemoteFilePath + "zddi-private2.tar.gz"
	}else if strings.Contains(version, "3.16") || strings.Contains(version, "3.17"){
		return Const.ZddiPrivateMenuName + "zddi-private3.tar.gz", Const.ZddiRemoteFilePath + "zddi-private3.tar.gz"
	} else {
		return "", ""
	}
}
