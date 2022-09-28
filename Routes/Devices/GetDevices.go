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

func GetDevice(context *gin.Context) {
	db, open_db_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
	if open_db_err!=nil {
		logrus.Debug("open db file", open_db_err)
	}
	ssh_bat := []Struct.SshStruct{}
	ssh_list := Struct.BatchSshStruct{}
	db.Find(&ssh_bat)
	ssh_list.BatchSsh = ssh_bat
		context.JSON(http.StatusOK,ssh_list)
}
