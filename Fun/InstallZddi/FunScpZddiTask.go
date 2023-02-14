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

func InstallScpZddiTask(task *Struct.ScpTask) (string, error) {
	//生成日志
	logrus := Util.CreateLogger(Const.ZddiTaskLogPath+task.TaskName, Const.ZddiTaskLogPath+task.TaskName+"/"+task.TaskName+".log")
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
	license_version := "DNSv" + strconv.Itoa(task.DnsVersion) + ", DHCPv" + strconv.Itoa(task.DhcpVersion) + ", ADDv" + strconv.Itoa(task.AddVersion)
	task_db := Struct.InstallZddiTaskStruct{
		TaskName:       task.TaskName,
		Nodes:          nodes,
		LicenseVersion: license_version,
		BuildPkg:       task.GetScpBuild.Path,
		ZddiPkg:        task.GetScpZddi.Path,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		StatusMsg:      "start task ....",
	}

	//如果校验错误为空,则开始任务
	if check_info, check_err := task.CheckScpTask(task_db, task.Colony, logrus); check_err == nil {
		go func() (string, error) {
			{
				logrus.Info("start task " + task.TaskName + " ...")
				//插入数据库
				task_db.DBCreateTask(*logrus)
				//日志文件路径
				zddi_file_name := Util.GetFileName(task.GetScpZddi.Path)
				build_file_name := Util.GetFileName(task.GetScpBuild.Path)

				//本地MD5
				local_zddi_md5, _ := Util.GetFileMd5(Const.ZddiFileMenuName + zddi_file_name)
				local_build_md5, _ := Util.GetFileMd5(Const.ZddiFileMenuName + build_file_name)

				//获取远端文件的MD5
				zddi_client, connErr := task.GetScpZddi.Conn()
				if connErr != nil {
					return "conn err" + task.GetScpZddi.Ipaddr, connErr
				}
				remote_zddi_md5, _ := task.GetScpZddi.Exec(zddi_client, "md5sum "+task.GetScpZddi.Path+" | awk {'print $1'}")
				build_client, connErr := task.GetScpBuild.Conn()
				if connErr != nil {
					return "conn err" + task.GetScpBuild.Ipaddr, connErr
				}
				remote_build_md5, _ := task.GetScpBuild.Exec(build_client, "md5sum "+task.GetScpBuild.Path+" | awk {'print $1'}")

				logrus.Debug("local zddi md5: " + local_zddi_md5)
				logrus.Debug("remote zddi md5: " + remote_zddi_md5)
				logrus.Debug("local build md5: " + local_build_md5)
				logrus.Debug("remote build md5: " + remote_build_md5)
				//查看本地是否有文件
				if _, get_file_err := os.Stat(Const.ZddiFileMenuName + zddi_file_name); get_file_err != nil {
					//没有直接去远端获取
					logrus.Warning("local not file" + Const.ZddiFileMenuName + zddi_file_name + " start get file" + task.GetScpZddi.Path)
					task_db.UpdateTaskMsg("Get Zddi Pkg ...")
					if msg, getFileErr := task.GetScpZddi.GetFile(Const.ZddiFileMenuName+zddi_file_name, task.GetScpZddi.Path, logrus); getFileErr != nil {
						return msg, getFileErr
					}
				} else {
					//有的话查看下MD5是否一致
					if remote_zddi_md5 == local_zddi_md5 {

						logrus.Info("local " + zddi_file_name + " exist! skip get")
					} else {
						//不一致就去远端获取
						logrus.Warning("remote file " + zddi_file_name + " diff md5")
						logrus.Warning("local :" + " /" + local_zddi_md5 + "/")
						logrus.Warning("remote" + " /" + remote_zddi_md5 + "/")
					}
				}
				//查看本地是否有文件
				if _, get_file_err := os.Stat(Const.ZddiFileMenuName + build_file_name); get_file_err != nil {
					logrus.Info("local not file" + Const.ZddiFileMenuName + build_file_name + " start get file" + task.GetScpZddi.Path)
					task_db.UpdateTaskMsg("Get Build Pkg ...")
					if msg, getFileErr := task.GetScpZddi.GetFile(Const.ZddiFileMenuName+build_file_name, task.GetScpBuild.Path, logrus); getFileErr != nil {
						return msg, getFileErr
					}
				} else {
					//有的话查看下MD5是否一致
					if local_build_md5 == remote_build_md5 {
						logrus.Debug("local " + zddi_file_name + " exist! skip get")
					} else {
						//不一致就去远端获取
						logrus.Warning("remote file " + build_file_name + " diff md5")
						logrus.Warning("local :" + " /" + local_build_md5 + "/")
						logrus.Warning("remote" + " /" + remote_build_md5 + "/")
						task_db.UpdateTaskMsg("Get Build Pkg ...")
						if msg, getFileErr := task.GetScpBuild.GetFile(Const.ZddiFileMenuName+build_file_name, task.GetScpBuild.Path, logrus); getFileErr != nil {
							return msg, getFileErr
						}
					}
				}
				logrus.Info(" all need file get succ !!!")
				//部署zddi
				go func() {
					{
						task_db.UpdateTaskMsg("Install Zddi in Devices ...")
						//开始部署ddi
						wg := sync.WaitGroup{}
						for _, zddi_device := range task.ZddiDevices {
							wg.Add(1)
							go func(zddi_device Struct.ScpStruct) {
								zddi_device.InstallZddi(zddi_file_name, build_file_name, task.DnsVersion, task.AddVersion, task.DhcpVersion, zddi_device.Role, task.TaskName)
								wg.Done()
							}(zddi_device)
						}
						wg.Wait()
						logrus.Info("all node install succ !!")
						//是否搭建集群
						if task.Colony == true {
							task_db.UpdateTaskMsg("start add zddi group ...")
							logrus.Info("start add zddi group")
							Struct.StartCreateColony(task.ZddiDevices, task.GetScpBuild.Path, logrus)
						} else {
							logrus.Info("task exec succ ...")
						}
					}
					task_db.UpdateTaskMsg("task exec succ ...")
				}()
			}
			return "", nil
		}()
		return "check task " + task.TaskName + " succ , please wait ...", nil
	} else {
		return check_info, check_err
	}
}
