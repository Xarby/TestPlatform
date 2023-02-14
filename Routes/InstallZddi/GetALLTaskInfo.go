package InstallZddi

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func GetAllTaskInfo(context *gin.Context) {
	db, open_db_err := gorm.Open(sqlite.Open(Const.InstallZddiTaskSqlPath), &gorm.Config{})
	if open_db_err!=nil {
		logrus.Error("open db file", open_db_err)
	}
	task_bat := []Struct.InstallZddiTaskStruct{}
	task_list:= Struct.BatchInstallZddiTaskStruct{}
	db.Find(&task_bat)
	task_list.BatchSsh=task_bat
	context.JSON(http.StatusOK, task_list)
}