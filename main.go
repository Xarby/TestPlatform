package main

import (
	"TestPlatform/Const"
	"TestPlatform/Routes/Devices"
	"TestPlatform/Routes/DialingTest"
	"TestPlatform/Routes/InstallZddi"
	"TestPlatform/Util"
	"github.com/gin-gonic/gin"
)
func main()  {


	Util.InitLog()

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
	dialing_test:= r.Group("/dialing_test")
	{
		dialing_test.POST("/create_task",DialingTest.CreateTask)
	}

	r.Run(":"+Const.RunPort)
}
