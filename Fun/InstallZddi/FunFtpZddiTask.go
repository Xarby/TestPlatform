package InstallZddi

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"
)

func InstallFtpZddiTask(task *Struct.FtpTask) (string, error) {
	//生成日志
	logrus := Util.CreateLogger(Const.LogPath+task.TaskName, Const.LogPath+task.TaskName+"/"+task.TaskName+".log")
	//打印任务
	if task_info, err := json.MarshalIndent(task, "", "\t"); err != nil {
		logrus.Debug(err)
	} else {
		logrus.Info(string(task_info))
	}

	//生成数据库信息
	var nodes string
	for _, zddi_device := range task.ZddiDevices {
		nodes = nodes + zddi_device.Ipaddr + " "
	}
	license_version := "DNSv" + strconv.Itoa(task.DnsVersion) + ", DHCPv" + strconv.Itoa(task.DhcpVersion) + ", ADDv2" + strconv.Itoa(task.AddVersion)
	task_db := Struct.InstallZddiTaskStruct{
		TaskName:       task.TaskName,
		Nodes:          nodes,
		LicenseVersion: license_version,
		BuildPkg:       task.BuildPath,
		ZddiPkg:        task.ZddiPath,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		StatusMsg:      "start task ....",
	}

	//校验任务

	//如果校验错误为空,则开始任务
	if check_info, check_err := task.CheckFtpTask(task_db, task.Colony, logrus); check_err == nil {
		go func() {
			task_db.UpdateTaskMsg("Install Zddi in Devices ...")
			wg := sync.WaitGroup{}
			zddi_file_name := Util.GetFileName(task.ZddiPath)
			build_file_name := Util.GetFileName(task.BuildPath)
			local_zddi_file_name := Const.ZddiFileMenuName + zddi_file_name
			local_build_file_name := Const.ZddiFileMenuName + build_file_name
			//获取文件
			if _, open_err := os.Stat(Const.ZddiFileMenuName + zddi_file_name); open_err != nil {
				logrus.Warning("local exist file "+Const.ZddiFileMenuName+zddi_file_name, " start get!")
				task_db.UpdateTaskMsg("Get Zddi Pkg ...")
				task.Ftp.GetFtpFile(task.ZddiPath)
			} else {
				logrus.Info("Is exist file " + local_zddi_file_name + ", skip get file")
			}
			if _, open_err := os.Stat(Const.ZddiFileMenuName + build_file_name); open_err != nil {
				logrus.Warning("local exist file "+Const.ZddiFileMenuName+build_file_name, " start get!")
				task_db.UpdateTaskMsg("Get Build Pkg ...")
				task.Ftp.GetFtpFile(task.BuildPath)
			} else {
				logrus.Info("Is exist file " + local_build_file_name + ", skip get file")
			}
			//开始部署zddi
			for _, zddi_device := range task.ZddiDevices {
				//开始安装
				wg.Add(1)
				go func(zddi_device Struct.ScpStruct) {
					zddi_device.InstallZddi(zddi_file_name, build_file_name, task.DnsVersion, task.AddVersion, task.DhcpVersion, zddi_device.Role, task.TaskName)
					wg.Done()
				}(zddi_device)
			}
			wg.Wait()
			//开始添加节点
			if task.Colony == true {
				task_db.UpdateTaskMsg("start add zddi group ...")
				logrus.Info("start add zddi group")
				Struct.StartCreateColony(task.ZddiDevices, task.BuildPath, logrus)
			}

		}()
		task_db.UpdateTaskMsg("task exec succ ...")
		return "check task " + task.TaskName + " succ , please wait ...", nil
	} else {
		return check_info, check_err
	}

}
