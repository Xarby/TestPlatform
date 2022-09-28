package Struct

import (
	"fmt"
	"testing"
)

func TestScp(t *testing.T) {
	scp := SshStruct{
		Ipaddr:   "10.1.104.201",
		Port:     "22",
		User:     "root",
		Password: "zdns@knet.cn",
	}
	fmt.Println(scp)
	scp.CheckBackup()
}

