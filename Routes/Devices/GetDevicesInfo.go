package Devices

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
)

func GetDeviceInfo(context *gin.Context) {
	//打开连接数据库
	db, open_db_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
	db.AutoMigrate(Struct.SshStruct{})
	db.AutoMigrate(Struct.DevInfoStruct{})
	if open_db_err!=nil {
		log.Println("open db file", open_db_err)
	}
	//存放查询到的主机
	ssh_list := []Struct.SshStruct{}
	//存放查询到的主机的结构体
	ssh_bat := Struct.BatchSshStruct{}
	//存放主机或许信息的批量信息结构体
	dev_info_bat := Struct.BatchDevInfoStruct{}

	//查询并赋值给ssh_list
	db.Find(&ssh_list)
	ssh_bat.BatchSsh = ssh_list

	//开始执行任务
	wg := sync.WaitGroup{}
	var mutex sync.Mutex
	db.Delete(&Struct.DevInfoStruct{})

	//读取一个机器的信息并开始获取信息并存储
	for _, dev := range ssh_bat.BatchSsh {
		wg.Add(1)
		go func(sshStruct Struct.SshStruct) {
			getInfo, get_err := sshStruct.GetDevInfo()
			log.Println(getInfo)
			if get_err!=nil {
				log.Println(get_err)
			}
			mutex.Lock()
			dev_info_bat.BatchDevInfo= append(dev_info_bat.BatchDevInfo, getInfo)
			db.Create(getInfo)
			mutex.Unlock()
			wg.Done()
		}(dev)
	}
	wg.Wait()
	context.JSON(http.StatusOK,dev_info_bat)
}