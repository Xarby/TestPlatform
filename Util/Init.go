package Util

import (
	"TestPlatform/Const"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func InitProgram() error {
	//初始化文件夹
	initDir()
	if runtime.GOOS == "linux" {
		InitLicense()
	}
	GinLogInit()
	return JudgeFile([]string{
		Const.RsyncFilePath+Const.RsyncFilePathCentos7X86,
		Const.RsyncFilePath+Const.RsyncFilePathKylin10ARM,
		Const.RsyncFilePath+Const.RsyncFilePathopenEulerX86,
		Const.RsyncFilePath+Const.RsyncFilePathCentos6X86,
		Const.LocalBackupShellFile})
}

// 创建目录
func CreateDir(path string) {
	_, file_err := os.Stat(path)
	if file_err != nil {
		os.MkdirAll(path, os.ModePerm)
	}
}

// 初始化文件夹
func initDir() {
	CreateDir(Const.LogPath)
	CreateDir(Const.ZddiTaskLogPath)
	CreateDir(Const.DevicesLogPath)
	CreateDir(Const.DBMenuName)
	CreateDir(Const.DialingTestFilePath)
	CreateDir(Const.ZddiFileMenuName)
	CreateDir(Const.ZddiPrivateMenuName)
	CreateDir(Const.ZddiPrivateTarGzPath)
	CreateDir(Const.TempLicensePath)
	CreateDir(Const.MachineLicensePath)
	CreateDir(Const.MachineLicensePath)
}

// 日志初始化
func GinLogInit() {
	//gin框架日志init
	gin_log, _ := os.Create(Const.GinLogPath)
	gin.DefaultWriter = io.MultiWriter(gin_log, os.Stdout)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
		ForceColors:     true,
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
	logrus.SetLevel(Const.DebugLevel)
}

func CreateLogger(path string, log_path string) *logrus.Logger {
	CreateDir(path)
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
		ForceColors:     true,
	})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	file, err := os.OpenFile(log_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{
		file,
		os.Stdout}
	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		log.SetOutput(fileAndStdoutWriter)
	} else {
		log.Info("failed to log to file.")
	}
	//设置最低loglevel
	log.SetLevel(logrus.InfoLevel)
	return log
}

func JudgeFile(paths []string) error {
	for _, path := range paths {
		_, file_err := os.Stat(path)
		if file_err != nil {
			logrus.Error(path + ": not exist file")
			return errors.New(path + ": not exist file")
		}
	}
	return nil
}

func InitLicense() (string,error) {
	parmPath,_ := GetPronPath()
	parmPathNewLicense := parmPath+Const.ZddiPrivateNewPath
	parmPathOldLicense := parmPath+Const.ZddiPrivateOldPath
	OsExecCmd("gcc -g -o "+parmPathNewLicense+"license -lcrypto "+parmPathNewLicense+"license.c  "+parmPathNewLicense+"dig_license.c "+parmPathNewLicense+"gen_clientId.c")
	OsExecCmd("/bin/mv "+parmPathNewLicense+"license /usr/local/bin/license_new")
	OsExecCmd("gcc -g -o "+parmPathOldLicense+"license -lcrypto "+parmPathOldLicense+"license.c  "+parmPathOldLicense+"dig_license.c "+parmPathOldLicense+"gen_clientId.c")
	OsExecCmd("/bin/mv "+parmPathOldLicense+"license /usr/local/bin/license_old")
	return "",nil
}


func OsExecCmd(cmd string)(string,error){
	logrus.Info("Local Exec :"+cmd)
	out, err := exec.Command("/bin/bash","-c",cmd).CombinedOutput()
	if err!=nil{
		logrus.Warning(string(out), err.Error() )
		return string(out), err
	}else {
		return string(out), nil
	}
}


