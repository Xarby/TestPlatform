package Util

import (
	"TestPlatform/Const"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func InitLog()  {

	CreateDir(Const.LogPath)
	CreateDir(Const.DevicesInfoMenuName)
	CreateDir(Const.DialingTestFilePath)
	CreateDir(Const.ZddiFileMenuName)
	CreateDir(Const.ZddiPrivateMenuName)
	gin_log, _ := os.Create(Const.GinLogPath)
	gin.DefaultWriter = io.MultiWriter(gin_log, os.Stdout)
	work_path,_ := os.Create(Const.WorkLogPath)
	log.SetOutput(work_path)
}

func CreateDir(path string)  {
	_,file_err := os.Stat(path)
	if file_err!= nil{
		log.Println("start create dir "+path)
		os.MkdirAll(path,os.ModePerm)
	}
}