package Struct

import (
	"TestPlatform/Const"
	"TestPlatform/Util"
	"errors"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

//公共部分的校验
func PubCheck(TaskDB InstallZddiTaskStruct, DnsVersion int, DhcpVersion int, AddVersion int, ZddiPath string, BuildPath string, check_colony bool, ZddiDevices []ScpStruct, logrus *logrus.Logger) (string, error) {
	//数据库中是否已经存在相同的任务名
	if checkMsg, checkErr := CheckTaskName(TaskDB); checkErr != nil {
		return checkMsg, checkErr
	} else {
		logrus.Info("check task db succ ...")
	}
	//校验license版本
	if checkMsg, checkErr := CheckLicense(DnsVersion, DhcpVersion, AddVersion); checkErr != nil {
		return checkMsg, checkErr
	} else {
		logrus.Info(checkMsg)
	}
	//校验文件名合法
	if checkMsg, checkErr := CheckFileformat(ZddiPath, BuildPath,ZddiDevices); checkErr != nil {
		return checkMsg, checkErr
	} else {
		logrus.Info(checkMsg)
	}
	//校验架构是否正常
	if checkMsg, checkErr := CheckArchitecture(BuildPath, ZddiDevices); checkErr != nil {
		return checkMsg, checkErr
	} else {
		logrus.Info(checkMsg)
	}
	//校验角色信息
	if checkMsg, checkErr := CheckRole(ZddiDevices); checkErr != nil {
		return checkMsg, checkErr
	} else {
		logrus.Info(checkMsg)
	}
	//是否校验集群角色是否合法
	if check_colony {
		if checkMsg, checkErr := CheckColonyRole(ZddiDevices); checkErr != nil {
			return checkMsg, checkErr
		} else {
			logrus.Info("check colony model role check succ ...")
		}
		if checkMsg, checkErr := CheckColonyNetworkCard(ZddiDevices); checkErr != nil {
			return checkMsg, checkErr
		} else {
			logrus.Info("check colony ha network name succ ...")
		}
		if checkMsg, checkErr := ChechHAVIPConn(ZddiDevices); checkErr != nil {
			return checkMsg, checkErr
		} else {
			logrus.Info("Check HA VIP succ ...")
		}
	} else {
		logrus.Info("skip check colony model role !")
	}
	//校验环境是否干净
	if checkMsg, checkErr := CheckZddiDeviceEnvironmental(ZddiDevices); checkErr != nil {
		return checkMsg, checkErr
	} else {
		logrus.Info(checkMsg)
	}
	logrus.Info("check environmental succ ...")
	return "check succ start task ...", nil
}

//校验架构
func CheckArchitecture(BuildPath string, scp_struct_list []ScpStruct) (string, error) {
	for _, scp_dev := range scp_struct_list {
		client, _ := scp_dev.Conn()
		exe_recult, _ := scp_dev.Exec(client, "arch")
		if strings.Contains(BuildPath, "x86") {
			fmt.Println("456")
			if strings.Contains(exe_recult, "x86") {
				continue
			} else if strings.Contains(exe_recult, "aarch64") {
				return "check arch fail :" + scp_dev.Ipaddr + " arch is " + exe_recult + " pkg is " + BuildPath, errors.New("check arch fail :" + scp_dev.Ipaddr + " unknow arch :" + exe_recult)
			} else {
				return "check arch fail :" + scp_dev.Ipaddr + " unknow arch :" + exe_recult, errors.New("check arch fail :" + scp_dev.Ipaddr + " unknow arch :" + exe_recult)
			}
		} else if strings.Contains(BuildPath, "aarch64") {
			if strings.Contains(exe_recult, "aarch64") {
				continue
			} else if strings.Contains(exe_recult, "x86") {
				return "check arch fail :" + scp_dev.Ipaddr + " arch is " + exe_recult + " pkg is " + BuildPath, errors.New("check arch fail :" + scp_dev.Ipaddr + " is " + exe_recult + " pkg is " + BuildPath)
			} else {
				return "check arch fail :" + scp_dev.Ipaddr + " unknow arch :" + exe_recult, errors.New("check arch fail :" + scp_dev.Ipaddr + " unknow arch :" + exe_recult)
			}
		}
	}
	return "check license version succ ...", nil
}

//校验license
func CheckLicense(DnsVersion int, DhcpVersion int, AddVersion int) (string, error) {
	if DnsVersion >= 4 || DnsVersion < 0 {
		logrus.Error("dns version input limit 0-3")
		return "dns version input limit 0-3", errors.New("dns version input limit 0-3")
	}
	if AddVersion >= 3 || AddVersion < 0 {
		logrus.Error("add version input limit 0-2")
		return "add version input limit 0-2", errors.New("add version input limit 0-2")
	}
	if DhcpVersion >= 2 || DhcpVersion < 0 {
		logrus.Error("dhcp version input limit 0-1")
		return "dhcp version input limit 0-1", errors.New("dhcp version input limit 0-1")
	}
	return "check license version succ ...", nil
}

//选择对应的private
func ChoosePrivateVersion(version string) (string, string, error) {
	if strings.Contains(version, "3.10") || strings.Contains(version, "3.11") {
		return Const.LocalZddiPrivateTarGzOld, Const.RemoteZddiPrivateTarGzOld, nil
	} else if (strings.Contains(version, "3.12") || strings.Contains(version, "3.13") || strings.Contains(version, "3.14") || strings.Contains(version, "3.15")|| strings.Contains(version, "3.49")) && strings.Contains(version,"x86_64") {
		return Const.LocalZddiPrivateTarGzNew, Const.RemoteZddiPrivateTarGzNew, nil
	} else {
		return "", "", errors.New("Unknow Vserions " + version)
	}
}

//校验文件格式
func CheckFileformat(zddi_path string, build_path string, scp_struct_list []ScpStruct) (string, error) {
	if strings.Contains(build_path, "rpm") == false {
		return "build pkg file format don't  'rpm' ", errors.New("build pkg file format don't  'rpm'")
	}
	if strings.Contains(zddi_path, "rpm") == false && strings.Contains(zddi_path, "tar.gz") == false {
		return "work pkg file format don't  'rpm' or 'tar.gz' ", errors.New("work pkg file format don't  'rpm' or 'tar.gz'")
	}

	//假如包带ky 那么应该能读取到kylin-release
	if strings.Contains(build_path, "ky") {
		for _, scp_dev := range scp_struct_list {
			client, conn_err := scp_dev.Conn()
			if conn_err != nil {
				return "conn dev " + scp_dev.Ipaddr + " err", errors.New("conn dev " + scp_dev.Ipaddr + " err " + conn_err.Error())
			}
			exe_recult, _ := scp_dev.Exec(client, "ls /etc/kylin-release")
			fmt.Println(exe_recult)
			if strings.Contains(exe_recult, "No such file or directory") {
				error_msg := "get dev "+scp_dev.Ipaddr+" info : os not is Kylin10 !"
				return error_msg, errors.New(error_msg)
			} else {
				continue
			}
		}
		//假如包带oe 那么应该能读取到openEuler-release
	}else if strings.Contains(build_path, "oe"){
		for _, scp_dev := range scp_struct_list {
			client, conn_err := scp_dev.Conn()
			if conn_err != nil {
				return "conn dev " + scp_dev.Ipaddr + " err", errors.New("conn dev " + scp_dev.Ipaddr + " err " + conn_err.Error())
			}
			exe_recult, _ := scp_dev.Exec(client, "ls /etc/openEuler-release")
			fmt.Println(exe_recult)
			if strings.Contains(exe_recult, "No such file or directory") {
				error_msg := "get dev "+scp_dev.Ipaddr+" info : os not is openEuler !"
				return error_msg, errors.New(error_msg)
			} else {
				continue
			}
		}
		//假如包带el7 那么应该能读取到redhat-release的7.8
	}else if strings.Contains(build_path, "el7"){
		for _, scp_dev := range scp_struct_list {
			client, conn_err := scp_dev.Conn()
			if conn_err != nil {
				return "conn dev " + scp_dev.Ipaddr + " err", errors.New("conn dev " + scp_dev.Ipaddr + " err " + conn_err.Error())
			}
			exe_recult, _ := scp_dev.Exec(client, "cat /etc/redhat-release")
			fmt.Println(exe_recult)
			if !strings.Contains(exe_recult, "7.8") {
				error_msg := "get dev "+scp_dev.Ipaddr+" info : os not is Centos7.8 !"
				return error_msg, errors.New(error_msg)
			} else {
				continue
			}
		}
		//假如包以上都不带 那么应该能读取到redhat-release的6.4
	}else {
		for _, scp_dev := range scp_struct_list {
			client, conn_err := scp_dev.Conn()
			if conn_err != nil {
				return "conn dev " + scp_dev.Ipaddr + " err", errors.New("conn dev " + scp_dev.Ipaddr + " err " + conn_err.Error())
			}
			exe_recult, _ := scp_dev.Exec(client, "cat /etc/redhat-release")
			fmt.Println(exe_recult)
			if !strings.Contains(exe_recult, "6.4") && !strings.Contains(exe_recult, "6.10"){
				error_msg := "get dev "+scp_dev.Ipaddr+" info : os not is Centos6.4 or Centos6.10!"
				return error_msg, errors.New(error_msg)
			} else {
				continue
			}
		}
	}
	return "check file format succ ...", nil
}

//校验集群
func CheckColonyRole(scp_struct_list []ScpStruct) (string, error) {
	master_count := 0
	master_ha_flag := 0

	master_ha_temp := 0
	slave_ha_temp := 0
	baskcup_ha_temp := 0

	for _, scpStruct := range scp_struct_list {

		switch scpStruct.Role {
		case "master":
			master_ha_flag++
			master_count++

		case "slave":
			continue
		case "backup":
			continue
		case "master_m":
			master_ha_flag++
			master_ha_temp++
			if checkMsg, checkErr := CheckRemoteIPAndVip(scp_struct_list, scpStruct.Ipaddr, scpStruct.RemoteIp, scpStruct.VIp); checkErr != nil {
				return checkMsg, checkErr
			}
		case "master_s":
			master_ha_temp--
			if checkMsg, checkErr := CheckRemoteIPAndVip(scp_struct_list, scpStruct.Ipaddr, scpStruct.RemoteIp, scpStruct.VIp); checkErr != nil {
				return checkMsg, checkErr
			}
		case "slave_m":
			slave_ha_temp++
			if checkMsg, checkErr := CheckRemoteIPAndVip(scp_struct_list, scpStruct.Ipaddr, scpStruct.RemoteIp, scpStruct.VIp); checkErr != nil {
				return checkMsg, checkErr
			}
		case "slave_s":
			slave_ha_temp--
			if checkMsg, checkErr := CheckRemoteIPAndVip(scp_struct_list, scpStruct.Ipaddr, scpStruct.RemoteIp, scpStruct.VIp); checkErr != nil {
				return checkMsg, checkErr
			}
		case "backup_m":
			baskcup_ha_temp++
			if checkMsg, checkErr := CheckRemoteIPAndVip(scp_struct_list, scpStruct.Ipaddr, scpStruct.RemoteIp, scpStruct.VIp); checkErr != nil {
				return checkMsg, checkErr
			}
		case "backup_s":
			baskcup_ha_temp--
			if checkMsg, checkErr := CheckRemoteIPAndVip(scp_struct_list, scpStruct.Ipaddr, scpStruct.RemoteIp, scpStruct.VIp); checkErr != nil {
				return checkMsg, checkErr
			}
		default:
			return scpStruct.Ipaddr + " role: '" + scpStruct.Role + "' not in master/slave/backup/master_s/slave_m/backup_m/master_s/slave_s/backup_s", errors.New("not in not in master/slave/backup/master_s/slave_m/backup_m/master_s/slave_s/backup_s")
		}

	}
	if master_ha_temp != 0 {
		return "master ha num check fail", errors.New("master ha num check fail")
	}
	if slave_ha_temp != 0 {
		return "slave ha num check fail", errors.New("master ha num check fail")
	}
	if baskcup_ha_temp != 0 {
		return "backup ha num check fail", errors.New("master ha num check fail")
	}

	//判断是不是有两个master节点
	if master_count >= 2 {
		return "master node exist " + strconv.Itoa(master_count), errors.New("master node exist " + strconv.Itoa(master_count))
	}
	//判断是不是有master也有master_ha
	if master_ha_flag >= 2 {
		return "master and master_ha node exist ", errors.New("master node exist ")
	}
	return "", nil
}

//判断能否找到对端IP并核对
func CheckRemoteIPAndVip(scpStructList []ScpStruct, localIp string, rempteIp string, Vip string) (string, error) {
	for _, scpStruct := range scpStructList {
		if scpStruct.RemoteIp == localIp {
			if scpStruct.Ipaddr == rempteIp {
				if scpStruct.VIp == Vip {
					return "", nil
				} else {
					return localIp + " checke ha error :" + scpStruct.Ipaddr + " vip " + Vip + " diff", errors.New(localIp + "checke ha error :" + scpStruct.Ipaddr + " vip " + Vip + " diff")
				}

			} else {
				return localIp + "checke ha error :" + scpStruct.Ipaddr + " diff " + rempteIp, errors.New(localIp + "checke ha error :" + scpStruct.Ipaddr + "!=" + rempteIp)
			}
		}
	}
	return "not find ha remote dev " + rempteIp, errors.New("not find ha dev" + rempteIp)
}

//校验FTP服务器
func CheckFtpServer(ftp_server FtpStruct, zddi_path string, build_path string) (string, error) {
	url := ftp_server.Ipaddr + ":" + ftp_server.Port
	ftp_client, conn_err := ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))
	if conn_err != nil {
		logrus.Error("conn ftp server"+ftp_server.Ipaddr+" fail ...", conn_err)
		return "conn ftp server" + ftp_server.Ipaddr + " fail ...", conn_err
	} else {
		//假如连接成功
		login_err := ftp_client.Login(ftp_server.User, ftp_server.Password)
		if login_err != nil {
			logrus.Error("login ftp server"+ftp_server.Ipaddr+" fail ... "+ftp_server.User+"/"+ftp_server.Password, login_err)
			return "login ftp server" + ftp_server.Ipaddr + " fail ... " + ftp_server.User + "/" + ftp_server.Password, login_err
		} else {
			//查看ftp文件是否存在
			logrus.Info("login ftp server succ !")
			_, read_build_err := ftp_client.FileSize(build_path)
			if read_build_err != nil {
				logrus.Error("ftp server "+ftp_server.Ipaddr+" not find file "+build_path, read_build_err)
				return "ftp server " + ftp_server.Ipaddr + " not find file " + build_path, read_build_err
			}
			_, read_zddi_err := ftp_client.FileSize(zddi_path)
			if read_zddi_err != nil {
				logrus.Error("ftp server "+ftp_server.Ipaddr+" not find file "+zddi_path, read_zddi_err)
				return "ftp server " + ftp_server.Ipaddr + " not find file " + zddi_path, read_zddi_err
			}
		}
	}
	return "login ftp server succ !", nil
}

//校验环境是否干净
func CheckZddiDeviceEnvironmental(scpStruct []ScpStruct) (string, error) {
	for _, device := range scpStruct {
		client, conn_scp_err := device.Conn()
		if conn_scp_err != nil {
			logrus.Error("conn fail zddi "+device.Ipaddr, conn_scp_err)
			return "conn fail zddi " + device.Ipaddr, conn_scp_err
		} else {
			exe_rpm_result, _ := device.Exec(client, "rpm -qa | grep zddi")
			exe_recovery_result, _ := device.Exec(client, "zdns-recovery-tool")
			exe_environmental_result, _ := device.Exec(client, "systemd-detect-virt")
			//假如安装了rpm包
			if strings.Contains(exe_rpm_result, "zddi") {
				//安装了工具未安装备份工具
				if strings.Contains(exe_recovery_result, "command not found") || strings.Contains(exe_recovery_result, "incomplete") {
					// 并且还是硬件
					if strings.Contains(exe_environmental_result, "none") {
						return device.Ipaddr + " environmental pollution", errors.New(device.Ipaddr + " environmental pollution")
					}
					//安装了工具并且备份成功了
				}
			}
		}
	}
	return "check environmental succ ...", nil
}

//校验scp文件是否存在
func CheckScpFileExist(scp_dev ScpStruct) (string, error) {
	scp_conn, conn_scp_err := scp_dev.Conn()
	defer func() {
		if scp_conn != nil {
			scp_conn.Close()
		}
	}()
	if conn_scp_err != nil {
		return "conn fail zddi " + scp_dev.Ipaddr, conn_scp_err
	} else {
		exe_result, err := scp_dev.Exec(scp_conn, "ls "+scp_dev.Path)
		if err != nil {
			return "exe get remote file cmd fail :" + scp_dev.Ipaddr + ":" + scp_dev.Path, err
		} else if strings.Contains(exe_result, "No such file or directory") {
			return "get remote build fail :" + scp_dev.Ipaddr + ":" + scp_dev.Path, err
		}
	}
	return "file" + scp_dev.Path + " exist", nil
}

//校验任务名称是否已存在
func CheckTaskName(TaskDB InstallZddiTaskStruct) (string, error) {
	if len(TaskDB.TaskName) == 0 {
		logrus.Error("task name is nil,please input task name")
		return "task name is nil,please input task name", errors.New("task name is nil,please input task name")
	} else if len(TaskDB.TaskName) > 30 {
		logrus.Error("task name " + TaskDB.TaskName + " exceed 30 character,please del character")
	}
	_, flag, find_err := TaskDB.DBFindTask()
	if find_err != nil {
		logrus.Error(find_err)
	}
	if flag == false {
		return "check db succ !!!", nil
	} else {
		return "exist task " + TaskDB.TaskName, errors.New("exist task " + TaskDB.TaskName)
	}
}

//校验HA网卡名字是否一样
func CheckColonyNetworkCard(scp_struct_list []ScpStruct) (string, error) {
	Dev := make(map[string]string)
	for _, scpStruct := range scp_struct_list {

		//获取所有的HA节点的网卡名
		if scpStruct.RemoteIp != "" {
			if local_client, conn_err := scpStruct.Conn(); conn_err != nil {
				return "conn " + scpStruct.Ipaddr + " faile", conn_err
			} else {
				exec_result, _ := scpStruct.Exec(local_client, "ifconfig | awk '/"+scpStruct.Ipaddr+"/{print a}{a=$1}'")
				Dev[scpStruct.Ipaddr] = exec_result
			}
		}
	}
	for _, scpStruct := range scp_struct_list {
		//对比自己和远端的网卡名是不是一样
		if scpStruct.RemoteIp != "" {
			if Dev[scpStruct.Ipaddr] == Dev[scpStruct.RemoteIp] {
				continue
			} else {
				return "HA dev " + scpStruct.Ipaddr + " network is" + Dev[scpStruct.Ipaddr] + " " + scpStruct.RemoteIp + " network name is " + Dev[scpStruct.RemoteIp] + "diff", errors.New("HA dev " + scpStruct.Ipaddr + " network is" + Dev[scpStruct.Ipaddr] + " " + scpStruct.RemoteIp + " network name is " + Dev[scpStruct.RemoteIp] + "diff")
			}
		}
	}
	return "HA network card name check succ ...", nil
}

func ChechHAVIPConn(scp_struct_list []ScpStruct) (string, error) {
	for _, scpStruct := range scp_struct_list {
		if scpStruct.VIp != "" {
			flag, _ := Util.TryConn(scpStruct.VIp, 22, 1)
			if flag == false {
				continue
			} else {
				return scpStruct.Ipaddr + " VIP " + scpStruct.VIp + " occupied , please change !!!", errors.New(scpStruct.Ipaddr + " VIP " + scpStruct.VIp + " , please change !!!")
			}
		}
	}
	return "Check HA VIP succ ...", nil
}

//校验角色
func CheckRole(ZddiDevices []ScpStruct) (string, error) {
	for _, scpStruct := range ZddiDevices {
		if scpStruct.Role == "" {
			return scpStruct.Ipaddr + " role is nil", errors.New(scpStruct.Ipaddr + " role is nil")
		} else if scpStruct.Role == "master_m" || scpStruct.Role == "master_s" || scpStruct.Role == "slave_m" || scpStruct.Role == "slave_s" || scpStruct.Role == "backup_s" || scpStruct.Role == "backup_m" || scpStruct.Role == "master" || scpStruct.Role == "slave" || scpStruct.Role == "backup" {
			continue
		} else {
			return "unknow role " + scpStruct.Role, errors.New("unknow role " + scpStruct.Role)
		}
	}
	return "check role succ", nil
}
