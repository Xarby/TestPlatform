package ZdnsBackup

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"errors"
	"golang.org/x/crypto/ssh"
	"strings"
)

func CheckRecovery(device Struct.SshStruct) (string,error) {
	logrus := Util.CreateLogger(Const.RecoveryLogPath, Const.RecoveryLogPath+"/recovery.log")

	//开始连接
	client, conn_err := device.Conn()
	if conn_err != nil {
		return "", errors.New("conn device " + device.Ipaddr + " fail")
	} else {
		logrus.Info("conn device " + device.Ipaddr + " succ !")
	}

	//校验rsync
	if exe_result, _ := device.Exec(client, "rpm -qa | grep Rsync"); strings.Contains(exe_result, "Rsync") == false {
		warn_msg := "Environment not installed Rsync , please install Rsync ..."
		logrus.Warning(warn_msg)
		arch ,_:= device.Execute(client, "arch")
		system_info ,_:= device.Execute(client, "cat /etc/system-release")
		if strings.Contains(arch,"x86_64")&& strings.Contains(system_info,"7.8"){
			device.PutFile(Const.RsyncFilePath+Const.RsyncFilePathCentos7X86, Const.ZddiRemoteFilePath+Const.RsyncFilePathCentos7X86, logrus)
		}else if strings.Contains(arch,"x86_64")&& strings.Contains(system_info,"openEuler"){
			device.PutFile(Const.RsyncFilePath+Const.RsyncFilePathopenEulerX86, Const.ZddiRemoteFilePath+Const.RsyncFilePathopenEulerX86, logrus)
		}else if strings.Contains(arch,"aarch64")&& strings.Contains(system_info,"V10"){
			device.PutFile(Const.RsyncFilePath+Const.RsyncFilePathopenEulerX86, Const.ZddiRemoteFilePath+Const.RsyncFilePathopenEulerX86, logrus)
		}else {
			err_msg := "Not super arch " + arch+" system info "+system_info
			logrus.Error(err_msg)
			return "", errors.New(err_msg)
		}

	} else {
		logrus.Info("device " + device.Ipaddr + " exist Rsync !")
	}


	file, process, rpm, _ := device.ShowStatu()
	//纯干净环境
	if file == 0 && process == 0 && rpm == 0 {
		result_msg := device.Ipaddr + " not find backup file , please start backup task ,..."
		logrus.Info(result_msg)
		return result_msg, nil
	}else if process == 1 {  //正在执行任务
		result_msg := device.Ipaddr + " exist execute backup or recovery task , please wait ..."
		logrus.Warning(result_msg)
		return result_msg, nil
	}

	////开始还原
	go startRecovery(device,client)
	return "recovery succ", nil
}


func startRecovery(device Struct.SshStruct,client *ssh.Client)  (string, error){
	device.Execute(client, "zdns-recovery-tool recovery")
	return "", nil
}