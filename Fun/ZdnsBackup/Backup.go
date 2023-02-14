package ZdnsBackup

import (
	"TestPlatform/Struct"
)

func BackgroundBackupZddi(device Struct.SshStruct) (string, error) {
	if check_info,check_err:=device.CheckBackupTask();check_err!=nil{
		return check_info, check_err
	}else {
		go device.StartBackup()
		return check_info, nil
	}
}