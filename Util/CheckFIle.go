package Util

import (
	"TestPlatform/Const"
	"log"
	"os"
)


func CheckZddiFile(zddi_file_path string)  bool {
	path := GetFileName(zddi_file_path)
	if _,err :=os.Stat(Const.ZddiFileMenuName+path);err != nil{
		log.Println("local "+Const.ZddiFileMenuName+"not exist file"+path)
		return false
	}else {
		return true
	}
}

func CheckZddiFileMenu()  {
	if _,err :=os.Stat(Const.ZddiFileMenuName);err != nil{
		log.Println("create Path:"+Const.ZddiFileMenuName+"file ")
		os.MkdirAll(Const.ZddiFileMenuName,0777)
	}
}