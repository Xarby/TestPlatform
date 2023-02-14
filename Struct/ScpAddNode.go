package Struct

import (
	"TestPlatform/Const"
	"TestPlatform/Util"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func (scp_dev *ScpStruct) newAddNode(master_ip string, logrus *logrus.Logger) (string, error) {
	name := scp_dev.Role + scp_dev.Ipaddr[strings.LastIndex(scp_dev.Ipaddr, ".")+1:]
	role := ""
	if strings.Contains(scp_dev.Role, "master") {
		role = "master"
	} else if strings.Contains(scp_dev.Role, "slave") {
		role = "slave"
	} else if strings.Contains(scp_dev.Role, "backup") {
		role = "backup"
	}

	body := []byte(fmt.Sprintf(`{
					"name": "%s",
					"ip": "%s",
					"username": "%s",
					"password": "%s",
					"role": "%s",
					"group": "local",
					"is_extend":"no"}`, name, scp_dev.Ipaddr, Const.ZddiAddNodeUser, Const.ZddiAddNodePasswd, role))
	logrus.Info("start post:"+string(body))
	if requestsErr, _ := Util.PostRequests("POST", fmt.Sprintf("https://%s:20120/groups/local/members", master_ip), body); requestsErr != nil {
		logrus.Error("add node"+scp_dev.Ipaddr+"role:"+scp_dev.Role+"fail", requestsErr)
		//超时则重试三次
		flag := 0
		if (strings.Contains(requestsErr.Error(), "Timeout") || strings.Contains(requestsErr.Error(), "timeout")) && flag < 3 {
			flag = flag + 1
			time.Sleep(time.Second * 3)
			logrus.Info("try " + strconv.Itoa(flag))
			if requestsErr, _ := Util.PostRequests("POST", fmt.Sprintf("https://%s:20120/groups/local/members", master_ip), body); requestsErr != nil {
				logrus.Error("add node"+scp_dev.Ipaddr+"role:"+scp_dev.Role+"fail", requestsErr)
			} else {
				logrus.Info("add node " + scp_dev.Ipaddr + "role: " + scp_dev.Role + " succ")
				return "add node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ", nil
			}
		}
		return "add node" + scp_dev.Ipaddr + " role: " + scp_dev.Role + " fail", requestsErr
	}
	logrus.Info("add node " + scp_dev.Ipaddr + "role: " + scp_dev.Role + " succ")
	return "add node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ", nil
}

func (scp_dev ScpStruct) newAddHaNode(master_ip string, logrus *logrus.Logger) (string, error) {
	client, _ := scp_dev.Conn()
	network_card_name, _ := scp_dev.Exec(client, "ip addr | grep -B 2 "+scp_dev.Ipaddr+" | head -n 1 | awk -F: '{ print $2 }' | tr -d [:blank:]")
	name := scp_dev.Role[:len(scp_dev.Role)-1] + "m" + scp_dev.RemoteIp[strings.LastIndex(scp_dev.VIp, ".")+1:]
	owners := "local." + name
	//判断节点IP
	if scp_dev.Role == "master_s" {
		owners = "local.master"
	} else {
		owners = "local." + name
	}
	iface := network_card_name
	vips1 := scp_dev.VIp
	vips2 := network_card_name
	vr_id := scp_dev.VIp[strings.LastIndex(scp_dev.Ipaddr, ".")+1:]
	node_owners := "node_ha_" + scp_dev.VIp[strings.LastIndex(scp_dev.Ipaddr, ".")+1:]
	body := []byte(fmt.Sprintf(`{"role": "ha","group": "local","name": "ha_%s","ip": "%s","username": "%s","password": "%s",
					"owners": ["%s"],"member_ha_info": {"iface": "%s","master_ip": "%s/24","backup_ip": "%s/24",
					"vips": ["%s/24 %s"],"vr_id": "%s","arp_interval": "10","preempt": "yes",
					"owners": ["%s"]},"DNS_LICENSE": false,"DHCP_LICENSE": false,"ADD_LICENSE": false}`, name, scp_dev.Ipaddr, Const.ZddiAddNodeUser, Const.ZddiAddNodePasswd, owners, iface, scp_dev.RemoteIp, scp_dev.Ipaddr, vips1, vips2, vr_id, node_owners))
	logrus.Info("start post:"+string(body))
	if requestsErr, _ := Util.PostRequests("POST", fmt.Sprintf("https://%s:20120/ha/members", master_ip), body); requestsErr != nil {
		logrus.Error("add node "+scp_dev.Ipaddr+" role: "+scp_dev.Role+"fail", requestsErr)
		flag := 0
		if (strings.Contains(requestsErr.Error(), "Timeout") || strings.Contains(requestsErr.Error(), "timeout")) && flag < 3 {
			flag = flag + 1
			time.Sleep(time.Second * 3)
			logrus.Info("try " + strconv.Itoa(flag))
			if requestsErr, _ := Util.PostRequests("POST", fmt.Sprintf("https://%s:20120/ha/members", master_ip), body); requestsErr != nil {
				logrus.Error("add node"+scp_dev.Ipaddr+"role:"+scp_dev.Role+"fail", requestsErr)
			} else {
				logrus.Info("add node " + scp_dev.Ipaddr + "role: " + scp_dev.Role + " succ")
				return "add node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ", nil
			}
		}
		return "add HA node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + "fail", requestsErr
	}
	time.Sleep(time.Second * 5)
	logrus.Info("add HA node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ")
	return "add HA node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ", nil
}

func (scp_dev ScpStruct) oldAddNode(master_ip string, logrus *logrus.Logger) (string, error) {
	name := scp_dev.Role + scp_dev.Ipaddr[strings.LastIndex(scp_dev.Ipaddr, ".")+1:]
	body := []byte(fmt.Sprintf(`{
					"name": "%s",
					"ip": "%s",
					"username": "%s",
					"password": "%s",
					"role": "%s",
					"group": "local"}`, name, scp_dev.Ipaddr, Const.ZddiAddNodeUser, Const.ZddiAddNodePasswd, scp_dev.Role))
	logrus.Info("start post:"+string(body))
	if requestsErr, _ := Util.PostRequests("POST", fmt.Sprintf("https://%s:20120/groups/local/members", master_ip), body); requestsErr != nil {
		logrus.Error("add node "+scp_dev.Ipaddr+" role: "+scp_dev.Role+" fail", requestsErr)
		flag := 0
		if (strings.Contains(requestsErr.Error(), "Timeout") || strings.Contains(requestsErr.Error(), "timeout")) && flag < 3 {
			flag = flag + 1
			time.Sleep(time.Second * 3)
			logrus.Info("try " + strconv.Itoa(flag))
			logrus.Info("start post:"+string(body))
			if requestsErr, _ := Util.PostRequests("POST", fmt.Sprintf("https://%s:20120/groups/local/members", master_ip), body); requestsErr != nil {
				logrus.Error("add node"+scp_dev.Ipaddr+"role:"+scp_dev.Role+"fail", requestsErr)
			} else {
				logrus.Info("add node " + scp_dev.Ipaddr + "role: " + scp_dev.Role + " succ")
				return "add node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ", nil
			}
		}
		return "add node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " fail", requestsErr
	}
	logrus.Info("add  node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ")
	return "add  node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ", nil
}

func (scp_dev ScpStruct) oldAddHaNode(logrus *logrus.Logger) (string, error) {
	client, _ := scp_dev.Conn()
	//获取网卡名列表
	requestsErr, getResult := Util.PostRequests("GET", fmt.Sprintf("https://%s:20120/groups/local/members/master/config/ha", scp_dev.Ipaddr), nil)
	if requestsErr != nil {
		logrus.Debug(requestsErr)
	} else {
		logrus.Info("get network info" + getResult)
	}
	//获取当前IP使用的网卡名
	network_card_name, _ := scp_dev.Exec(client, "ip addr | grep -B 2 "+scp_dev.Ipaddr+" | head -n 1 | awk -F: '{ print $2 }' | tr -d [:blank:]")
	//设备状态 /单机/主机/备机
	var dev_type string
	if strings.Contains(scp_dev.Role, "_m") {
		dev_type = "master"
	} else if strings.Contains(scp_dev.Role, "_s") {
		dev_type = "backup"
	} else {
		return "role  err " + scp_dev.Ipaddr + " role:" + scp_dev.Role, errors.New("role  err " + scp_dev.Ipaddr + " role:" + scp_dev.Role)
	}
	//判断角色
	vip := scp_dev.VIp
	//vid
	vr_id := scp_dev.VIp[strings.LastIndex(scp_dev.Ipaddr, ".")+1:]
	//当前设备的网卡名
	iface := network_card_name
	var ifaces string
	if strings.Index(getResult, "ifaces\":") != -1 && strings.Index(getResult, "ifaces\":") != -1 {
		ifaces = getResult[strings.Index(getResult, "ifaces\":")+8 : strings.Index(getResult, ",\"sync_data_time")]
	} else {
		logrus.Error(getResult)
		logrus.Error("not get ifaces "+scp_dev.Ipaddr, network_card_name)
		return "not get ifaces " + scp_dev.Ipaddr, errors.New("not get ifaces " + scp_dev.Ipaddr)
	}
	body := []byte(fmt.Sprintf(`{"type": "%s","vips": ["%s/24 %s"],"pair_ip": "%s",
													"vr_id": "%s","preempt": "yes","arp_interval": "10","iface": "%s",
													"ifaces":%s,
													"keepalive_status": "yes","use_vmac": "no","use_vmac_vrrp": "no",
													"ha_run_role": "%s"}`, dev_type, vip, network_card_name, scp_dev.RemoteIp, vr_id, iface, ifaces, scp_dev.Role[0:len(scp_dev.Role)-2]))
	logrus.Info("start post:"+string(body))
	if requestsErr, _ := Util.PostRequests("PUT", fmt.Sprintf("https://%s:20120/groups/local/members/master/config/ha", scp_dev.Ipaddr), body); requestsErr != nil {

		flag := 0
		if (strings.Contains(requestsErr.Error(), "Timeout") || strings.Contains(requestsErr.Error(), "timeout")) && flag < 3 {
			flag = flag + 1
			time.Sleep(time.Second * 3)
			logrus.Info("try " + strconv.Itoa(flag))
			if requestsErr, _ := Util.PostRequests("PUT", fmt.Sprintf("https://%s:20120/groups/local/members/master/config/ha", scp_dev.Ipaddr), body); requestsErr != nil {
				logrus.Error("add node"+scp_dev.Ipaddr+"role:"+scp_dev.Role+"fail", requestsErr)
			} else {
				logrus.Info("add node " + scp_dev.Ipaddr + "role: " + scp_dev.Role + " succ")
				return "add node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ", nil
			}
		}

		logrus.Error("add node "+scp_dev.Ipaddr+" role: "+scp_dev.Role+" fail", requestsErr)
		return "add node " + scp_dev.Ipaddr + " role:" + scp_dev.Role + " fail", requestsErr
	}
	time.Sleep(time.Second * 5)
	logrus.Info("add  node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ")
	return "add  node " + scp_dev.Ipaddr + " role: " + scp_dev.Role + " succ", nil
}

func PutMasterIP(master_ip string) (string, error) {
	//编辑master的本机IP
	logrus.Info("start put master cloud center ip!")
	time.Sleep(time.Second * 5)
	if requestsErr, _ := Util.PostRequests("PUT", fmt.Sprintf("https://%s:20120/groups/local/members/master", master_ip), []byte(fmt.Sprintf(`{
											"group": "local",
											"name": "master",
											"id": "master",
											"ip": "%s",
											"positionX": "0.423",
											"positionY": "0.406"
										}`, master_ip))); requestsErr != nil {
		logrus.Error("put master ip faile " + master_ip + requestsErr.Error())
		return "put master ip faile " + master_ip, requestsErr
	}
	logrus.Info("put master ip succ " + master_ip)
	return "put master ip succ " + master_ip, nil
}
