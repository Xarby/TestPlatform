package InstallZddi

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"log"
	"os"
)

func InstallFtpZddiTask(task *Struct.FtpTask) (string, error) {
	//打印任务
	Util.PrintJson(task)
	check_info, check_err := task.CheckFtpTask()
	if check_err == nil {
		zddi_file_name := Util.GetFileName(task.ZddiPath)
		build_file_name := Util.GetFileName(task.BuildPath)
		local_zddi_file_name := Const.ZddiFileMenuName + zddi_file_name
		local_build_file_name := Const.ZddiFileMenuName + build_file_name
		//获取文件
		if _, open_err := os.Stat(Const.ZddiFileMenuName + zddi_file_name); open_err != nil {
			task.Ftp.GetFtpFile(task.ZddiPath)
		} else {
			log.Println("Is exist file " + local_zddi_file_name + ", skip get file")
		}
		if _, open_err := os.Stat(Const.ZddiFileMenuName + build_file_name); open_err != nil {
			task.Ftp.GetFtpFile(task.BuildPath)
		} else {
			log.Println("Is exist file " + local_build_file_name + ", skip get file")
		}
		//开始部署zddi
		for _, zddi_device := range task.ZddiDevices {
			//开始安装
			go func(zddi_device Struct.ScpStruct) {
				zddi_device.InstallZddi(zddi_file_name, build_file_name, task.DnsVersion, task.AddVersion, task.DhcpVersion, zddi_device.Role)
			}(zddi_device)
		}
		//开始添加节点
		if task.Colony == true{
			if  deployment_info,deployment_err :=task.FtpStartCreateColony();deployment_err!= nil {
				return deployment_info,deployment_err
			}
		}
	}
	return check_info, check_err
}

