package Util

import (
	"TestPlatform/Const"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func InitProgram() {
	initDir()
	logInit()
}

func CreateDir(path string) {
	_, file_err := os.Stat(path)
	if file_err != nil {
		os.MkdirAll(path, os.ModePerm)
	}
}

//初始化文件夹
func initDir() {
	CreateDir(Const.LogPath)
	CreateDir(Const.DevicesInfoMenuName)
	CreateDir(Const.DialingTestFilePath)
	CreateDir(Const.ZddiFileMenuName)
	CreateDir(Const.ZddiPrivateMenuName)
}

//日志初始化
func logInit() {
	//gin框架日志init
	gin_log, _ := os.Create(Const.GinLogPath)
	gin.DefaultWriter = io.MultiWriter(gin_log, os.Stdout)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:true,    //键值对加引号
		TimestampFormat:"2006-01-02 15:04:05",  //时间格式
		FullTimestamp:true,
		ForceColors: true,
	})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	file, err := os.OpenFile(Const.WorkLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{
		file,
		os.Stdout}
	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		logrus.SetOutput(fileAndStdoutWriter)
	} else {
		logrus.Info("failed to log to file.")
	}
	//设置最低loglevel
	logrus.SetLevel(logrus.InfoLevel)
}
