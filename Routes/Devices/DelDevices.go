package Devices

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func DelDevice(context *gin.Context) {
	devs := Struct.DelDevList{}
	succ_list := []map[string]string{}
	fail_list := []map[string]string{}
	if err := context.ShouldBind(&devs); err != nil {
		logrus.Debug(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	} else {
		logrus.Info("del devices :",devs)
		db, open_db_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
		if open_db_err!=nil {
			logrus.Debug("open db file", err)
		}
		db.AutoMigrate(Struct.SshStruct{})
		for _, v := range devs.DevList {
			del_dev := db.First(&Struct.SshStruct{},"ipaddr",v)
			if (del_dev.RowsAffected) == 1{
				del_dev.Delete(1)
				succ_list = append(succ_list, map[string]string{v:"del succ !"})
				logrus.Info("del "+v+" succ!")
			}else {
				fail_list = append(fail_list, map[string]string{v:"exist devices skip del !"})
				logrus.Warning( "exist devices"+ v+" skip del !")
			}
		}
		context.JSON(http.StatusOK, map[string]interface{}{"fail_list":fail_list,"succ_list":succ_list})
	}
}

