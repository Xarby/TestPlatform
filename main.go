package main

import (
	"TestPlatform/Const"
	"TestPlatform/Routes/Devices"
	"TestPlatform/Routes/DialingTest"
	"TestPlatform/Routes/InstallZddi"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func main() {

	Util.InitLog()
	db, open_db_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
	if open_db_err != nil {
		log.Println("open db file", open_db_err)
		return
	}
	db.AutoMigrate(Struct.SshStruct{})
	db.AutoMigrate(Struct.DevInfoStruct{})

	r := gin.Default()
	install_zddi := r.Group("/install_zddi")
	{
		install_zddi.POST("/ftp", InstallZddi.FtpInstallZddi)
		install_zddi.POST("/scp", InstallZddi.ScpInstallZddi)
	}
	device := r.Group("/device")
	{
		device.GET("/get_dev_info", Devices.GetDeviceInfo)
		device.GET("/get_dev", Devices.GetDevice)
		device.POST("/add_dev", Devices.AddDevice)
		device.PUT("/put_dev", Devices.PutDevice)
		device.DELETE("/del_dev", Devices.DelDevice)
	}
	dialing_test := r.Group("/dialing_test")
	{
		dialing_test.POST("/create_task", DialingTest.CreateTask)
	}
	r.Run("0.0.0.0" + ":" + Const.RunPort)
}