package Struct

import (
	"TestPlatform/Util"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)


type BatchSshStruct struct {
	BatchSsh []SshStruct `json:"batch_ssh"`
}

type SshStruct struct {
	Ipaddr    string `json:"ipaddr" gorm:"primaryKey"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
}

func (dev *SshStruct) Conn() (*ssh.Client, error){
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 3,
		User:            dev.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(dev.Password)}
	client, err := ssh.Dial("tcp", dev.Ipaddr+":"+dev.Port, config)
	if err != nil {
		logrus.Error("conn devices " + dev.Ipaddr + " err " + err.Error())
		return nil,err
	} else {
		logrus.Debug("conn devices " + dev.Ipaddr + " succ")
	}
	return client,nil
}
func (ssh_dev *SshStruct) Exec(client *ssh.Client,cmd string) (string, error) {
	logrus.Info(ssh_dev.Ipaddr + " start exec cmd : " + cmd)
	session, net_session_err := client.NewSession()
	if net_session_err != nil {
		return "",net_session_err
	}
	result_exec, exe_err := session.CombinedOutput(cmd)
	if exe_err != nil {
		logrus.Warning(ssh_dev.Ipaddr + " start exec cmd : " + cmd+" /error info:"+exe_err.Error())
		return string(result_exec), exe_err
	}
	result := Util.SplitFormat(string(result_exec))
	return result, exe_err
}

func (ssh_dev *SshStruct) Execute(client *ssh.Client,cmd string) (string, error) {
	logrus.Info(ssh_dev.Ipaddr + " start exec cmd : " + cmd)
	session, net_session_err := client.NewSession()
	if net_session_err != nil {
		return "",net_session_err
	}
	result_exec, exe_err := session.CombinedOutput(cmd)
	if exe_err != nil {
		logrus.Warning(ssh_dev.Ipaddr + " start exec cmd : " + cmd+" /error info:"+exe_err.Error())
		return string(result_exec), exe_err
	}
	return string(result_exec), exe_err
}

func (dev *SshStruct) GetDevInfo() (DevInfoStruct, error) {
	ssh_client,conn_err := dev.Conn()
	defer func(){
	 if ssh_client != nil{
		 ssh_client.Close()
	 }
	}()
	if conn_err != nil {
		return DevInfoStruct{
			Ipaddr:     dev.Ipaddr,
			CardName:   "none",
			MacAddr:    "none",
			OptVersion: "none",
			RpmInfo:    "none",
			MemSize:    "none",
			CpuNum:     "none",
			CpuName:    "none",
			SkuNum:     "none",
			SnNum:      "none",
			DiskSize:   "none",
			DiskUse:    "none",
			Status:     conn_err.Error(),
		}, conn_err
	}
	dev_info := DevInfoStruct{}
	dev_info.CardName, _ = dev.Exec(ssh_client,"ip route | grep '" + string(dev.Ipaddr) + "' | awk '{print $3}'")
	dev_info.MacAddr, _ = dev.Exec(ssh_client,"cat /sys/class/net/" + strings.Replace(dev_info.CardName, "\n", "", -1) + "/address")
	dev_info.OptVersion, _ = dev.Exec(ssh_client,"cat /etc/redhat-release")
	rpm_temp, _ := dev.Exec(ssh_client,"rpm -qa | grep zddi_build")
	if rpm_temp == "" {
		dev_info.RpmInfo = "Not Install Zddi Build Pkg"
	} else {
		dev_info.RpmInfo = rpm_temp
	}
	dev_info.Ipaddr = dev.Ipaddr
	dev_info.MemSize, _ = dev.Exec(ssh_client,"free -m | grep \"Mem\" |awk {'print $2'}")
	dev_info.CpuNum, _ = dev.Exec(ssh_client,"lscpu | grep \"CPU(s):\" | grep -v \"NUMA\" | awk {'print $2'}")
	dev_info.CpuName, _ = dev.Exec(ssh_client,"lscpu | grep name |  cut -f2 -d: | sed 's/^[            ]\\+//'")
	dev_info.SkuNum, _ = dev.Exec(ssh_client,"dmidecode -t 1 | grep \"SKU\" | awk '{print $3}'\n")
	dev_info.SnNum, _ = dev.Exec(ssh_client,"dmidecode -t 1 | grep \"Serial Number\" | awk '{print $3}'\n")
	dev_info.DiskSize, _ = dev.Exec(ssh_client,"df -h | grep \"sda3\" | awk {'print $2'}")
	dev_info.DiskUse, _ = dev.Exec(ssh_client,"df -h | grep \"sda3\" | awk {'print $5'}")
	dev_info.Status = "conn succ"
	return dev_info, conn_err
}