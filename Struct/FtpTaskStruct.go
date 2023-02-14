package Struct

import (
	"github.com/sirupsen/logrus"
)

// 部署节点的结构

// FTP部署的结构
type FtpTask struct {
	TaskName    string      `json:"task_name"`
	Ftp         FtpStruct   `json:"ftp"`
	DnsVersion  int         `json:"dns_version"`
	AddVersion  int         `json:"add_version"`
	DhcpVersion int         `json:"dhcp_version"`
	ZddiPath    string      `json:"zddi_path"`
	BuildPath   string      `json:"build_path"`
	ZddiDevices []ScpStruct `json:"zddi_devices"`
	Colony      bool        `json:"colony"`
}

func (ftp_task FtpTask) CheckFtpTask(TaskDB InstallZddiTaskStruct, check_colony bool, logrus *logrus.Logger) (string, error) {
	if checkMsg, checkErr := PubCheck(TaskDB, ftp_task.DnsVersion, ftp_task.DhcpVersion, ftp_task.AddVersion, ftp_task.ZddiPath, ftp_task.BuildPath, check_colony, ftp_task.ZddiDevices, logrus); checkErr != nil {
		return checkMsg, checkErr
	}
	if checkMsg, checkErr := CheckFtpServer(ftp_task.Ftp, ftp_task.ZddiPath, ftp_task.BuildPath); checkErr != nil {
		return checkMsg, checkErr
	}

	logrus.Info("check ftp task succ start task !")
	return "check ftp task succ start task !", nil
}
