package InstallZddi

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

func InstallFtpZddiTask(task *Struct.FtpTask) (string, error) {
	//打印任务

	if task_info,err := json.MarshalIndent(task,"","\t"); err != nil{
		logrus.Debug(err)
	}else {
		logrus.Info(string(task_info))
	}
	//校验任务
	check_info, check_err := task.CheckFtpTask()
	//如果校验错误为空,则开始任务
	if check_err == nil {
		go func() {
			wg := sync.WaitGroup{}
			zddi_file_name := Util.GetFileName(task.ZddiPath)
			build_file_name := Util.GetFileName(task.BuildPath)
			local_zddi_file_name := Const.ZddiFileMenuName + zddi_file_name
			local_build_file_name := Const.ZddiFileMenuName + build_file_name
			//获取文件
			if _, open_err := os.Stat(Const.ZddiFileMenuName + zddi_file_name); open_err != nil {
				logrus.Warning("local exist file "+Const.ZddiFileMenuName + zddi_file_name, " start get!")
				task.Ftp.GetFtpFile(task.ZddiPath)
			} else {
				logrus.Info("Is exist file " + local_zddi_file_name + ", skip get file")
			}
			if _, open_err := os.Stat(Const.ZddiFileMenuName + build_file_name); open_err != nil {
				logrus.Warning("local exist file "+Const.ZddiFileMenuName + build_file_name, " start get!")
				task.Ftp.GetFtpFile(task.BuildPath)
			} else {
				logrus.Info("Is exist file " + local_build_file_name + ", skip get file")
			}


			//开始部署zddi
			for _, zddi_device := range task.ZddiDevices {
				//开始安装
				wg.Add(1)
				go func(zddi_device Struct.ScpStruct) {
					zddi_device.InstallZddi(zddi_file_name, build_file_name, task.DnsVersion, task.AddVersion, task.DhcpVersion, zddi_device.Role)
					wg.Done()
				}(zddi_device)
			}
			wg.Wait()
			//开始添加节点
			if task.Colony == true {
				logrus.Info("start add zddi group")
				task.FtpStartCreateColony()
			}
		} ()
	}
	return check_info, check_err
}
