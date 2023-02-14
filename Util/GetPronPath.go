package Util

import (
	"log"
	"os"
	"strings"
)

//获取当前程序目录
func GetPronPath() (string,error) {

	out, getPathErr :=  os.Getwd()
	if getPathErr != nil{
		log.Panic("Don't get program fail")
	}
	return strings.Replace(out, "\n", "", -1) +"/",nil
}
