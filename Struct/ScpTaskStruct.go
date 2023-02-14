package Struct

import (
	"github.com/sirupsen/logrus"
)

//  SCP部署的结构
type ScpTask struct {
	TaskName    string      `json:"task_name"`
	GetScpZddi  ScpStruct   `json:"get_scp_zddi"`
	GetScpBuild ScpStruct   `json:"get_scp_build"`
	DnsVersion  int         `json:"dns_version"`
	AddVersion  int         `json:"add_version"`
	DhcpVersion int         `json:"dhcp_version"`
	ZddiDevices []ScpStruct `json:"zddi_devices"`
	Colony      bool        `json:"colony"`
}

func (scp_task ScpTask) CheckScpTask(TaskDB InstallZddiTaskStruct, check_colony bool, logrus *logrus.Logger) (string, error) {
	//和ftp公用的校验规则
	logrus.Info("start check task " + TaskDB.TaskName)
	if checkMsg, checkErr := PubCheck(TaskDB,scp_task.DnsVersion, scp_task.DhcpVersion, scp_task.AddVersion, scp_task.GetScpZddi.Path, scp_task.GetScpBuild.Path, check_colony, scp_task.ZddiDevices, logrus); checkErr != nil {
		return checkMsg, checkErr
	}
	if checkMsg, checkErr := CheckScpFileExist(scp_task.GetScpBuild); checkErr != nil {
		return checkMsg, checkErr
	}
	if checkMsg, checkErr := CheckScpFileExist(scp_task.GetScpZddi); checkErr != nil {
		return checkMsg, checkErr
	}
	logrus.Info("start check task " + TaskDB.TaskName + " succ ...")
	return "task " + TaskDB.TaskName + "all check succ ...", nil
}
