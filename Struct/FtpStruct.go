package Struct

import (
	"TestPlatform/Const"
	"TestPlatform/Util"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type FtpStruct struct {
	Ipaddr string `json:"ipaddr"`
	Port string `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
}

func (ftpStruct *FtpStruct)GetFtpFile( ftpFilePath string) (string,error) {

	//ftp连接
	url := ftpStruct.Ipaddr + ":" + ftpStruct.Port
	ftp_client, conn_err := ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))
	if conn_err != nil {
		logrus.Error("conn ftp server succ !!"+ftpStruct.Ipaddr+" error:"+conn_err.Error())
		return "conn ftp "+ftpStruct.Ipaddr,conn_err
	}else {
		logrus.Info("conn ftp server succ !!")
	}
	//ftp登录
	if log_err := ftp_client.Login(ftpStruct.User, ftpStruct.Password); log_err != nil {
		return "login ftp "+ftpStruct.Ipaddr,conn_err
	}else {
		logrus.Info("login ftp server succ !!"+ftpStruct.Ipaddr)
	}
	//获取远端文件的流
	remote_buf, buf_err := ftp_client.Retr(ftpFilePath)
	defer remote_buf.Close()
	if buf_err != nil {
		logrus.Error("get remote file "+ftpFilePath+" fail", " error:"+buf_err.Error())
		return "get remote file "+ftpFilePath+" fail", buf_err
	}else {
		logrus.Info("get ftp file succ:"+ftpFilePath)
	}

	//本地创建一个文件并开流
	file_name := Util.GetFileName(ftpFilePath)
	outFile, out_file_err := os.Create(Const.ZddiFileMenuName + file_name)
	defer outFile.Close()
	if out_file_err != nil {
		logrus.Error("open local file fail :"+Const.ZddiFileMenuName + file_name,"error :"+out_file_err.Error())
		return "open local file fail :"+Const.ZddiFileMenuName + file_name, out_file_err
	}else {
		logrus.Info("open local file succ"+Const.ZddiFileMenuName + file_name)
	}

	//开始下载
	logrus.Info("start get file")
	if _,copy_err := io.Copy(outFile, remote_buf);copy_err != nil {
		logrus.Error("download file "+ftpFilePath+" fail","error :"+copy_err.Error())
		return "download file "+ftpFilePath+" fail",copy_err
	}else {
		logrus.Info("downloca succ"+ftpFilePath)
	}
	return "downloca succ"+ftpFilePath,nil
}
