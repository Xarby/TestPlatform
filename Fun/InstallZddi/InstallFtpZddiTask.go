package InstallZddi

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"log"
)


func InstallFtpZddiTask(task *Struct.FtpTask) {
	//打印任务
	Util.PrintJson(task)
	zddi_file_name := Util.GetFileName(task.ZddiPath)
	build_file_name := Util.GetFileName(task.BuildPath)
	local_zddi_file_name := Const.ZddiFileMenuName + zddi_file_name
	local_build_file_name := Const.ZddiFileMenuName + build_file_name
	//获取文件
	if Util.CheckZddiFile(task.ZddiPath) != true {
		task.Ftp.GetFtpFile(task.ZddiPath)
	} else {
		log.Println("Is exist file " + local_zddi_file_name + ", skip get file")
	}
	if Util.CheckZddiFile(task.BuildPath) != true {
		task.Ftp.GetFtpFile(task.BuildPath)
	} else {
		log.Println("Is exist file " + local_build_file_name + ", skip get file")
	}

	for _, zddi_device := range task.ZddiDevices {
		//开始安装

		go func(zddi_device Struct.ScpStruct) {
			zddi_device.InstallZddi(zddi_file_name, build_file_name, task.DnsVersion, task.AddVersion, task.DhcpVersion, zddi_device.Role)
		}(zddi_device)
	}
}
