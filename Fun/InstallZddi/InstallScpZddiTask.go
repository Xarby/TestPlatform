package InstallZddi

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"log"
)

func InstallScpZddiTask(task *Struct.ScpTask) (string,error) {
	Util.PrintJson(task)

	check_info,check_err := task.CheckScpTask()
	if check_err == nil {
		log.Println(check_info)

		//获取文件名字
		zddi_file_name := Util.GetFileName(task.GetScpZddi.Path)
		build_file_name := Util.GetFileName(task.GetScpBuild.Path)
		task.GetScpZddi.Conn()
		task.GetScpBuild.Conn()

		if err := task.GetScpZddi.GetFile(Const.ZddiFileMenuName+zddi_file_name, task.GetScpZddi.Path); err != nil {
			return "get file fail "+zddi_file_name,err
			//return err
		}
		if err := task.GetScpBuild.GetFile(Const.ZddiFileMenuName+build_file_name, task.GetScpBuild.Path); err != nil {
			return "get file fail "+zddi_file_name,err
		}
		for _, zddi_device := range task.ZddiDevices {
			go func(zddi_device Struct.ScpStruct) {
				zddi_device.InstallZddi(zddi_file_name, build_file_name, task.DnsVersion, task.AddVersion, task.DhcpVersion, zddi_device.Role)
			}(zddi_device)
		}
	}
	return check_info,nil
}
