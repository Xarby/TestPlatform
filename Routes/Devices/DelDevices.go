package Devices

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func DelDevice(context *gin.Context) {
	devs := Struct.DelDevList{}
	succ_list := []map[string]string{}
	fail_list := []map[string]string{}
	if err := context.ShouldBind(&devs); err != nil {
		log.Println(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	} else {
		log.Println("del devices :",devs)
		db, open_db_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
		if open_db_err!=nil {
			log.Println("open db file", err)
		}
		db.AutoMigrate(Struct.SshStruct{})
		for _, v := range devs.DevList {
			del_dev := db.First(&Struct.SshStruct{},"ipaddr",v)
			if (del_dev.RowsAffected) == 1{
				del_dev.Delete(1)
				succ_list = append(succ_list, map[string]string{v:"del succ !"})
				log.Println("del "+v+" succ!")
			}else {
				fail_list = append(fail_list, map[string]string{v:"exist devices skip del !"})
				log.Println( "exist devices"+ v+" skip del !")
			}
		}
		context.JSON(http.StatusOK, map[string]interface{}{"fail_list":fail_list,"succ_list":succ_list})
	}
}

