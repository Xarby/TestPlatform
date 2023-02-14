package ZdnsLicense

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func ZdnsUploadMachineLicense(Machine Struct.UpLoadMachineStruct) (string, error) {
	logrus := logrus.New()
	programPath, _ := Util.GetPronPath()
	tmpPath := programPath + Const.TempUploadMachineLicnesePath
	err := os.RemoveAll(tmpPath)
	if err != nil {
		errmsg := errors.New(tmpPath+" remove dir fail : "+err.Error())
		logrus.Error(errmsg)
		return "", errmsg
	}
	Util.CreateDir(tmpPath)
	fmt.Println(Machine)
	logrus.Debug(Machine)
	//加密
	//readBin,_:= os.ReadFile(tmpFilePath)
	//fmt.Println(base64.URLEncoding.EncodeToString(readBin))
	//切割字符串

	//文件目录
	local_machine_info := tmpPath + "machine.info"
	local_pub_key := tmpPath + "pub.key"
	local_pri_key := tmpPath + "pri.key"
	local_licnese_file := tmpPath + "license.file"

	////处理js和go base加密/解密的 "/"和"_"问题
	WaitDecodeStr := Machine.BinData[strings.Index(Machine.BinData, "base64,")+len("base64,"):]
	tempDecode := strings.Replace(strings.Replace(WaitDecodeStr, "/", "_", -1),"+","-",-1)
	logrus.Info("base64 info : "+tempDecode)
	DecodeStr, _ := base64.URLEncoding.DecodeString(tempDecode)
	//根据base64串生成文件
	os.WriteFile(local_machine_info, DecodeStr, 666)
	t := time.Now()
	ddi_version := Machine.Version.DdiVersion
	time_now := fmt.Sprintf("%d-%d-%d %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	//保留信息
	Util.OsExecCmd("echo '"+fmt.Sprintf("%+v", Machine)+"' > "+tmpPath+"license.info")
	if ddi_version == "3.10" || ddi_version == "3.11" || ddi_version == "3.12" {
		cmds := fmt.Sprintf("license_old -p '%s' -q '%s' -i '%s' -l '%s' -e 'license_start_time:#:%s;license_valid_days:#:-1;license_remote_auth_flag:#:0;license_flag:#:1;root_change_password_flag:#:%s;license_control_node_num:#:1;license_normal_flag:#:1;license_virtual_validate:#:1;user_name:#:test;license_type:#:0;device_type:#:0;elk:#:1;sgcloud_st_time:#:%s;sgcloud_vaild_days:#:-1;sgcloud_role:#:1' -r 'DNSv%s,DHCPv%s,ADDv%s,REGv0 0 0'", local_pub_key, local_pri_key, local_machine_info, local_licnese_file, time_now,Machine.Version.ChangePasswd, time_now, Machine.Version.DnsVersion, Machine.Version.DhcpVersion, Machine.Version.AddVersion)
		logrus.Debug(cmds)
		Util.OsExecCmd(cmds)
		cmds = "tar -zcvf "+tmpPath+"license.tar.gz "+tmpPath+"*"
		logrus.Debug(cmds)
		Util.OsExecCmd(cmds)
	} else if strings.Contains(ddi_version, "3.13") || strings.Contains(ddi_version, "3.14") || strings.Contains(ddi_version, "3.15") || strings.Contains(ddi_version, "3.49") {
		cmds := fmt.Sprintf("license_new -p '%s' -q '%s' -i '%s' -l '%s' -e 'license_start_time:#:%s;license_valid_days:#:-1;license_remote_auth_flag:#:0;root_change_password_flag:#:%s;device_type:#:0' -r 'DNSv%s,DHCPv%s,ADDv%s 0 0'", local_pub_key, local_pri_key, local_machine_info, local_licnese_file, time_now,Machine.Version.ChangePasswd, Machine.Version.DnsVersion, Machine.Version.DhcpVersion, Machine.Version.AddVersion)
		logrus.Debug(cmds)
		Util.OsExecCmd(cmds)
		cmds = "cd "+tmpPath+" && "+"tar -zcvf "+tmpPath+"license.tar.gz -C "+tmpPath+" *"
		logrus.Debug(cmds)
		Util.OsExecCmd(cmds)
	}

	//读流
	filePath := tmpPath+"license.tar.gz"
	fileByte,_ := os.ReadFile(filePath)
	EncodeStr := "data:application/x-gzip;base64,"+strings.Replace(strings.Replace(base64.URLEncoding.EncodeToString(fileByte), "_", "/", -1),"-","+",-1)
	logrus.Debug("result base info :"+EncodeStr)
	return EncodeStr, nil
}
