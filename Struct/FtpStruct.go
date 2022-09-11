package Struct

import (
	"TestPlatform/Const"
	"TestPlatform/Util"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"os"
	"time"
)

type FtpStruct struct {
	Ipaddr string `json:"ipaddr"`
	Port string `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
}

func (ftpStruct *FtpStruct)GetFtpFile( ftpFilePath string) map[string]any {

	url := ftpStruct.Ipaddr + ":" + ftpStruct.Port
	ftp_client, conn_err := ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))
	if conn_err != nil {
		return map[string]any{"result": false, "mes": "conn fail", "file_code": conn_err}
	}
	if log_err := ftp_client.Login(ftpStruct.User, ftpStruct.Password); log_err != nil {
		return map[string]any{"result": false, "mes": "login fail", "file_code": log_err}
	}
	remote_buf, buf_err := ftp_client.Retr(ftpFilePath)
	defer remote_buf.Close()
	if buf_err != nil {
		log.Println(ftpFilePath)
		return map[string]any{"result": false, "mes": "get file fail", "file_code": buf_err}
	}
	file_name := Util.GetFileName(ftpFilePath)

	Util.CheckZddiFileMenu()
	outFile, out_file_err := os.Create(Const.ZddiFileMenuName + file_name)
	defer outFile.Close()
	if out_file_err != nil {
		return map[string]any{"result": false, "mes": "new file" + Const.ZddiFileMenuName + file_name + "file", "file_code": out_file_err}
	}
	log.Println("start get file")
	if _,copy_err := io.Copy(outFile, remote_buf);copy_err != nil {
		return map[string]any{"result": false, "mes": "write file fail", "file_code": copy_err}
	}
	result := map[string]any{"result": true, "mes": "get file"+ftpFilePath+" in local:"+Const.ZddiFileMenuName + file_name+" succ"}
	log.Println(result)
	return result
}
