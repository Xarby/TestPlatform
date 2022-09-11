package Util

import (
	"TestPlatform/Const"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func InitLog()  {
	if err:=os.MkdirAll(Const.LogPath,os.ModePerm);err != nil{
		fmt.Println(err)
	}
	if err:=os.MkdirAll(Const.DevicesInfoMenuName,os.ModePerm);err != nil{
		fmt.Println(err)
	}
	gin_log, _ := os.Create(Const.GinLogPath)
	gin.DefaultWriter = io.MultiWriter(gin_log, os.Stdout)
	work_path,_ := os.Create(Const.WorkLogPath)
	log.SetOutput(work_path)
}
