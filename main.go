package main

import (
	"TestPlatform/Const"
	"TestPlatform/Routes/Devices"
	"TestPlatform/Routes/DialingTest"
	"TestPlatform/Routes/InstallZddi"
	"TestPlatform/Routes/ZdnsBackup"
	"TestPlatform/Routes/ZdnsLicense"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func main() {
	//初始化
	r := gin.Default()
	//监听接口
	install_zddi := r.Group("/install_zddi") //安装ZDDI
	{
		install_zddi.POST("/ftp", InstallZddi.FtpInstallZddi)
		install_zddi.POST("/scp", InstallZddi.ScpInstallZddi)
		install_zddi.GET("/get_all_task_info", InstallZddi.GetAllTaskInfo)
		install_zddi.GET("/get_task_info", InstallZddi.GetTaskLog)
		install_zddi.DELETE("/del_task_info", InstallZddi.DelTaskInfo)
		install_zddi.DELETE("/del_all_task_info", InstallZddi.DelAllTaskInfo)
	}
	device := r.Group("/device") //扫描设备
	{
		device.GET("/get_dev_info", Devices.GetDeviceInfo)
		device.GET("/get_dev", Devices.GetDevice)
		device.POST("/add_dev", Devices.AddDevice)
		device.PUT("/put_dev", Devices.PutDevice)
		device.DELETE("/del_dev", Devices.DelDevice)
	}
	dialing_test := r.Group("/dialing_test") //拨测
	{
		dialing_test.POST("/create_task", DialingTest.CreateTask)
	}
	zdns_recovery := r.Group("/zdns_recovery") //备份/恢复license
	{
		zdns_recovery.POST("/backup", ZdnsBackup.Backup)
		zdns_recovery.POST("/recovery", ZdnsBackup.Recovery)
		zdns_recovery.POST("/show_statu", ZdnsBackup.ShowStatu)
	}
	ZdnsLicnese:= r.Group("/license") //制作license
	{
		ZdnsLicnese.POST("/make_license",ZdnsLicense.MakeLicnese)
		ZdnsLicnese.POST("/upload_machine_license",ZdnsLicense.UploadMachineLicense)
	}
	r.Run("0.0.0.0" + ":" + Const.RunPort)
}

func init() {
	//新建需要的文件夹以及公共日志初始化
	if Util.InitProgram() != nil {
		logrus.Error("Inif Fail ...")
		os.Exit(0)
	}
	device_db, open_db_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
	if open_db_err != nil {
		logrus.Debug("open db file", open_db_err)
		return
	}
	//建数据库
	device_db.AutoMigrate(Struct.SshStruct{})
	device_db.AutoMigrate(Struct.DevInfoStruct{})

	Zddi_install_info_db, open_db_err := gorm.Open(sqlite.Open(Const.InstallZddiTaskSqlPath), &gorm.Config{})
	if open_db_err != nil {
		logrus.Debug("open db file", open_db_err)
		return
	}
	Zddi_install_info_db.AutoMigrate(Struct.InstallZddiTaskStruct{})
	logrus.Info("Server init Succ ...")
}
