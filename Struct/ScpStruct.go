package Struct

import (
	"TestPlatform/Const"
	"TestPlatform/Util"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
)

type Devs struct {
	Devices []SshStruct `json:"devices"`
}

// SCP获取文件结构体
type ScpStruct struct {
	SshStruct
	Path      string `json:"path"`
	Role      string `json:"role"`
	sshclient *ssh.Client
}

func (scp_dev *ScpStruct) PutFile(local string, remote string) error {
	ssh_client, conn_err := scp_dev.Conn()
	if conn_err != nil {
		return conn_err
	}
	scpclient, new_client_err := scp.NewClientBySSH(ssh_client)

	defer scpclient.Close()
	if new_client_err != nil {
		log.Println("ssh change scp fail")
	}

	if _, file_err := os.Stat(local); file_err != nil {
		log.Println("local not exist file " + local)
	}
	f, file_err := os.Open(local)

	if file_err != nil {
		log.Println("not exist file " + local)
	}
	log.Println("start push local file： " + local + " to " + scp_dev.Ipaddr + " " + remote + " ...")
	if push_err :=scpclient.CopyFromFile(context.Background(), *f, remote, "0655");push_err!= nil{
		log.Println(push_err)
	}
	log.Println("push ok ...")
	return nil
}

func (scp_dev *SshStruct) GetFile(local string, remote string) error {

	ssh_client, conn_err := scp_dev.Conn()
	if conn_err != nil {
		return conn_err
	}
	scpclient, new_client_err := scp.NewClientBySSH(ssh_client)
	if new_client_err != nil {
		log.Println("ssh change scp fail")
	}
	_, err := os.Stat(local)
	//判断是本地是否存在文件
	if err != nil {
		defer scpclient.Close()
		//判断远程是否存在文件
		client , connErr :=scp_dev.Conn()
		if connErr != nil{
			return connErr
		}
		defer client.Close()
		_ ,file_err:=scp_dev.Exec(client,"ls "+remote)
		//远端有文件
		if file_err == nil{
			srcFile, open_file_err := os.OpenFile(local, os.O_RDWR|os.O_CREATE, 0755)
			if open_file_err != nil {
				log.Println("open file " + local + " fail")
			}
			log.Println("start pull remote file： " + remote + " to " + local + " ...")
			if get_file_err := scpclient.CopyFromRemote(context.Background(), srcFile, remote); get_file_err != nil {
				log.Println("push err ...", get_file_err)
				return get_file_err
			}
			log.Println("push ok ...")
		}else {
			log.Println("remote not exist file:"+remote+" task stop ")
			return file_err
		}
	}
	return nil
}



//zddi_path string, build_path string, dns_version string, add_version string, dhcp_version string, role string
func (dev *ScpStruct) InstallZddi(zddi_path string, build_path string, dns_version string, add_version string, dhcp_version string, role string) error {

	local_zddi_file_name := Const.ZddiFileMenuName + zddi_path
	local_build_file_name := Const.ZddiFileMenuName + build_path
	remote_zddi_file_name := Const.ZddiRemoteFilePath + zddi_path
	remote_build_file_name := Const.ZddiRemoteFilePath + build_path
	local_private, remote_private := Util.ChooseVersion(build_path)
	ssh_client,conn_err := dev.Conn()
	if conn_err != nil{
		return conn_err
	}
	defer ssh_client.Close()
	//传输文件

	remote_zddi_md5, get_zddi_md5_err := dev.Exec(ssh_client,"md5sum " + remote_zddi_file_name)
	local_zddi_md5, _ := Util.GetFileMd5(local_zddi_file_name)


	if get_zddi_md5_err == nil {
		if (strings.Contains(remote_zddi_md5,local_zddi_md5) == false){
			dev.PutFile(local_zddi_file_name, remote_zddi_file_name)
		}else {
			log.Println("remote exist file "+remote_zddi_file_name+" md5:" + remote_zddi_md5)
		}

	}else {
		dev.PutFile(local_zddi_file_name, remote_zddi_file_name)
	}

	remote_build_md5, get_build_md5_err := dev.Exec(ssh_client,"md5sum " + remote_build_file_name)
	local_build_md5, _ := Util.GetFileMd5(local_build_file_name)
	if get_build_md5_err== nil {
		if (strings.Contains(remote_build_md5,local_build_md5) == false){
			dev.PutFile(local_build_file_name, remote_build_file_name)
		}else {
			log.Println("remote exist file "+remote_build_file_name+" md5:" + remote_build_md5)
		}
	}else {
		dev.PutFile(local_build_file_name, remote_build_file_name)
	}
	remote_private_md5, get_private_err := dev.Exec(ssh_client,"md5sum " + remote_private)
	local_private_md5, _ := Util.GetFileMd5(local_private)

	if get_private_err == nil {
		if  (strings.Contains(remote_private_md5,local_private_md5) == false){
			dev.PutFile(local_private, remote_private)
		}else {
			log.Println("remote exist file"+remote_private+" md5:" + remote_private_md5)
		}
	}else {
		dev.PutFile(local_private, remote_private)
	}

	log.Println("start install zddi " + dev.Ipaddr)
	dev.Exec(ssh_client,"rpm -ivh " + build_path + " --force")
	if strings.Contains(zddi_path, "rpm") {
		dev.Exec(ssh_client,"rpm -ivh " + zddi_path + " --force")
	} else {
		dev.Exec(ssh_client,"tar -xzvf " + zddi_path)
	}
	dev.Exec(ssh_client,"tar -xzvf " + remote_private)
	dev.Exec(ssh_client,"cd /root/zddi-private/publisher/ && sh install")
	dev.Exec(ssh_client,"lurker -c fetch -o /etc/machine.info")

	//3.15||3.16版本处理方式
	if strings.Contains(build_path, "3.13") || strings.Contains(build_path, "3.14") || strings.Contains(build_path, "3.15") || strings.Contains(build_path, "3.16") {
		dev.Exec(ssh_client,fmt.Sprintf("publisher -c license -p /etc/pub.key -q /etc/pri.key -i /etc/machine.info -l /etc/license.file -N %s -H %s -A %s -m 0", dns_version, dhcp_version, add_version))
		dev.Exec(ssh_client,"sed -i '/release/d' /etc/rc.d/rc.local")
		dev.Exec(ssh_client,fmt.Sprintf("sed -i '/CLISH_PATH/a /usr/local/appsys/normal/package/zdns_startmgr/zdns_startmgr_ctl release /root/zddi %s' /etc/rc.d/rc.local", role))
		dev.Exec(ssh_client,"nohup /etc/rc.local &")
		dev.Exec(ssh_client,"sed -i '/PermitRootLogin/d' /etc/ssh/sshd_config")
		dev.Exec(ssh_client,"sed -i '/StrictModes/a PermitRootLogin yes' /etc/ssh/sshd_config")
		dev.Exec(ssh_client,"service  sshd restart")
		//3.10||3.11||3.13||3.14版本处理方式
	} else if strings.Contains(build_path, "3.10") || strings.Contains(build_path, "3.11") {
		dev.Exec(ssh_client,fmt.Sprintf("publisher -c license -p /etc/pub.key -q /etc/pri.key -i /etc/machine.info -l /etc/license.file -N %s -H %s -A %s -R 1 -m 0", dns_version, dhcp_version, add_version))
		dev.Exec(ssh_client,fmt.Sprintf("sed -i 's#^lurker.*#lurker -c run -d /root/zddi -r %s -p /etc/pub.key -l /etc/license.file#g'  /etc/rc.d/rc.local", role))
		dev.Exec(ssh_client,"nohup /etc/rc.local &")
	}
	return nil
}

