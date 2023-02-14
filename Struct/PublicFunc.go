package Struct

import (
	"TestPlatform/Util"
	"errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func StartCreateColony(zddi_devices []ScpStruct, build_path string,logrus *logrus.Logger) (string, error) {
	logrus.Info("start Create Colony!")
	var master_ip string
	//获取master的IP  并循环等待master的443端口开放 slave的4583端口开放
	for _, scp_dev := range zddi_devices {
		var port int
		if scp_dev.Role == "master" || scp_dev.Role == "master_m" {
			master_ip = scp_dev.Ipaddr
			port = 443
		} else if scp_dev.Role == "master_s"{
			port = 443
		} else {
			port = 4583
		}
		var count = 0
		for {
			if conn_flag, conn_err := Util.TryConn(scp_dev.Ipaddr, port,3); conn_flag == true {
				logrus.Info("conn ip:" + scp_dev.Ipaddr + " port " + strconv.Itoa(port) + " succ")
				break
			} else {
				//时间超过 60则放弃
				if count > 500 {
					logrus.Error("start "+scp_dev.Ipaddr+"server time exceed 180s", errors.New("start "+scp_dev.Ipaddr+"server time exceed 180s"))
					return "start " + scp_dev.Ipaddr + "server time exceed 180s", errors.New("start " + scp_dev.Ipaddr + "server time exceed 180s")
				} else {
					//时间+3秒并打印
					count = count + 3
					logrus.Warning(scp_dev.Ipaddr+" sleep 3s ", conn_err)
					time.Sleep(time.Second * 3)
				}
			}
		}
	}
	logrus.Info("check all dev port open!")
	logrus.Info("master_ip:" + master_ip)
	//编辑master的本机IP
	if _, put_err := PutMasterIP(master_ip); put_err != nil {
		return "put master ip faile " + master_ip, put_err
	}

	if strings.Contains(build_path, "3.15") || strings.Contains(build_path, "3.16") || strings.Contains(build_path, "3.17") ||strings.Contains(build_path, "3.49") {
		//添加非HA节点和HA主节点
		for _, scp_dev := range zddi_devices {
			if scp_dev.Role == "slave" || scp_dev.Role == "backup"  || scp_dev.Role == "slave_m" || scp_dev.Role == "backup_m" {
				logrus.Info("new add node ip:"+scp_dev.Ipaddr+" role:"+scp_dev.Role)
				for i := 0; i < 4; i++ {
					time.Sleep(time.Second*2)
					if _, err := scp_dev.newAddNode(master_ip,logrus);err == nil{
						break
					}
				}
			}
		}
		//添加HA非主节点
		for _, scp_dev := range zddi_devices {
			if scp_dev.Role == "master_s" || scp_dev.Role == "slave_s" || scp_dev.Role == "backup_s" {
				logrus.Info("new add ha node")
				for i := 0; i < 4; i++ {
					time.Sleep(time.Second * 2)
					if _, err := scp_dev.newAddHaNode(master_ip, logrus);err == nil{
						break
					}
				}
			}
		}
	} else {
		//添加所有非master_ha节点和非master节点
		for _, scp_dev := range zddi_devices {
			if scp_dev.Role == "master_s" || scp_dev.Role == "master_m"||scp_dev.Ipaddr == master_ip{
				continue
			} else {
				logrus.Info("old add node ip:"+scp_dev.Ipaddr+" role:"+scp_dev.Role)
				scp_dev.oldAddNode(master_ip,logrus)
			}
		}
		//用来决策master还是slave
		for _, scp_dev := range zddi_devices {
			if scp_dev.Role == "master_m" || scp_dev.Role == "slave_m" || scp_dev.Role == "backup_m" {
				logrus.Info("old add ha node "+scp_dev.Ipaddr)
				scp_dev.oldAddHaNode(logrus)
			}
		}
		for _, scp_dev := range zddi_devices {
			if scp_dev.Role == "master_s" || scp_dev.Role == "slave_s" || scp_dev.Role == "backup_s" {
				scp_dev.oldAddHaNode(logrus)
			}
		}
	}

	logrus.Info("add all node succ task shutdown !!!")
	return "add all node succ", nil
}
