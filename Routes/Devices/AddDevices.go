package Devices

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func AddDevice(context *gin.Context) {
	logrus := Util.CreateLogger(Const.DevicesLogPath,Const.DevicesLogPath+ "/devices.log")
	devs := Struct.Devs{}
	succ_list := []map[string]string{}
	fail_list := []map[string]string{}
	if err := context.ShouldBind(&devs); err != nil {
		logrus.Debug(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	} else {

		logrus.Info("add devices :", devs)
		db, open_db_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
		if open_db_err != nil {
			logrus.Debug("open db file", open_db_err)
		}

		//开始循环
		for _, v := range devs.Devices {
			add_dev := db.First(&Struct.SshStruct{}, "ipaddr", v.Ipaddr)
			_, conn_err := v.Conn()
			//假如存在则添加到报错中
			if (add_dev.RowsAffected) == 1 {
				fail_list = append(fail_list, map[string]string{v.Ipaddr: "exist device skip add !"})
				logrus.Warning("device " + v.Ipaddr + " exist ! skip ")
			} else if conn_err != nil {
				fail_list = append(fail_list, map[string]string{v.Ipaddr: "conn device skip add !"})
				logrus.Warning("device " + v.Ipaddr + " conn fail ! skip ")
			} else {
				if _, conn_err := devs.Devices[0].Conn(); conn_err != nil {
					context.JSON(http.StatusOK, map[string]interface{}{"fail_list": fail_list, "succ_list": succ_list})
				}
				db.Create(&v)
				succ_list = append(succ_list, map[string]string{v.Ipaddr: "add succ !"})
				logrus.Info("add " + v.Ipaddr + " succ !")
			}
		}
		context.JSON(http.StatusOK, map[string]interface{}{"fail_list": fail_list, "succ_list": succ_list})
	}
}
