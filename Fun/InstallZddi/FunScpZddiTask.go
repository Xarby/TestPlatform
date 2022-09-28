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

func InstallScpZddiTask(task *Struct.ScpTask) (string, error) {
	//打印任务
	if task_info,err := json.MarshalIndent(task,"","\t"); err != nil{
		logrus.Debug(err)
	}else {
		logrus.Info(string(task_info))
	}
	//校验任务
	check_info, check_err := task.CheckScpTask()
	//如果校验错误为空,则开始任务
	if check_err == nil {
		//文件名称
		logrus.Info("check task info succ..")
		zddi_file_name := Util.GetFileName(task.GetScpZddi.Path)
		build_file_name := Util.GetFileName(task.GetScpBuild.Path)

		//本地MD5
		local_zddi_md5, _ := Util.GetFileMd5(Const.ZddiFileMenuName + zddi_file_name)
		local_build_md5, _ := Util.GetFileMd5(Const.ZddiFileMenuName + build_file_name)
		//连接远程
		zddi_client, connErr := task.GetScpZddi.Conn()
		if connErr != nil {
			return "conn err" + task.GetScpZddi.Ipaddr, connErr
		}
		build_client, connErr := task.GetScpBuild.Conn()
		if connErr != nil {
			return "conn err" + task.GetScpBuild.Ipaddr, connErr
		}
		//远端文件MD5
		remote_zddi_md5, _ := task.GetScpZddi.Exec(zddi_client, "md5sum "+task.GetScpZddi.Path+" | awk {'print $1'}")
		remote_build_md5, _ := task.GetScpBuild.Exec(build_client, "md5sum "+task.GetScpBuild.Path+" | awk {'print $1'}")


		logrus.Info("-------------GET SCP FILE MD5 INFO-------------")
		logrus.Info("local zddi md5: "+local_zddi_md5)
		logrus.Info("remote zddi md5: "+remote_zddi_md5)
		logrus.Info("local build md5: "+local_build_md5)
		logrus.Info("remote build md5: "+remote_build_md5)
		logrus.Info("-------------------------------------------")

		//查看本地是否有文件
		if _, get_file_err := os.Stat(Const.ZddiFileMenuName + zddi_file_name); get_file_err != nil {
			//没有直接去远端获取
			logrus.Info("local not file" + Const.ZddiFileMenuName + zddi_file_name + " start get file" + task.GetScpZddi.Path)
			if msg, getFileErr := task.GetScpZddi.GetFile(Const.ZddiFileMenuName+zddi_file_name, task.GetScpZddi.Path); getFileErr != nil {
				return msg, getFileErr
			}
		} else {
			//有的话查看下MD5是否一致
			if remote_zddi_md5 == local_zddi_md5 {
				logrus.Info("remote file " + zddi_file_name + " and local same md5")
				logrus.Info("local :" +" /" + local_zddi_md5 +"/" )
				logrus.Info("remote" + " /" + remote_zddi_md5 + "/")
				logrus.Info("local " + zddi_file_name + " exist! skip get")
			} else {
				//不一致就去远端获取
				logrus.Warning("remote file " + zddi_file_name + " diff md5")
				logrus.Warning("local :" +" /" + local_zddi_md5 +"/" )
				logrus.Warning("remote" + " /" + remote_zddi_md5 + "/")
				if msg, getFileErr := task.GetScpZddi.GetFile(Const.ZddiFileMenuName+zddi_file_name, task.GetScpZddi.Path); getFileErr != nil {
					return msg, getFileErr
				}
			}
		}
		//查看本地是否有文件
		if _, get_file_err := os.Stat(Const.ZddiFileMenuName + build_file_name); get_file_err != nil {
			logrus.Info("local not file" + Const.ZddiFileMenuName + build_file_name + " start get file" + task.GetScpZddi.Path)
			if msg, getFileErr := task.GetScpZddi.GetFile(Const.ZddiFileMenuName+build_file_name, task.GetScpBuild.Path); getFileErr != nil {
				return msg, getFileErr
			}
		} else {
			//有的话查看下MD5是否一致
			if local_build_md5 == remote_build_md5 {
				logrus.Info("remote file " + build_file_name + " and local same md5")
				logrus.Info("local :" +" /" + local_build_md5 +"/" )
				logrus.Info("remote" + " /" + remote_build_md5 + "/")
				logrus.Info("local " + zddi_file_name + " exist! skip get")
			} else {
				//不一致就去远端获取
				logrus.Warning("remote file " + build_file_name + " diff md5")
				logrus.Warning("local :" +" /" + local_build_md5 +"/" )
				logrus.Warning("remote" + " /" + remote_build_md5 + "/")
				if msg, getFileErr := task.GetScpBuild.GetFile(Const.ZddiFileMenuName+build_file_name, task.GetScpBuild.Path); getFileErr != nil {
					return msg, getFileErr
				}
			}
		}
		logrus.Info(" all file get succ !!")
		//部署zddi
		go func() {
			{
				//开始部署ddi
				wg := sync.WaitGroup{}
				for _, zddi_device := range task.ZddiDevices {
					wg.Add(1)
					go func(zddi_device Struct.ScpStruct) {
						zddi_device.InstallZddi(zddi_file_name, build_file_name, task.DnsVersion, task.AddVersion, task.DhcpVersion, zddi_device.Role)
						wg.Done()
					}(zddi_device)

				}
				wg.Wait()
				//是否搭建集群
				if task.Colony == true {
					logrus.Info("start add zddi group")
					if _, deployment_err := task.ScpStartCreateColony(); deployment_err != nil {
						return
					}
				}
			}
		}()
	}
	//返回校验结果
	return check_info, check_err
}
