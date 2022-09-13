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

func AddDevice(context *gin.Context) {
	devs := Struct.Devs{}
	succ_list := []map[string]string{}
	fail_list := []map[string]string{}
	if err := context.ShouldBind(&devs); err != nil {
		log.Println(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	} else {
		log.Println("add devices :",devs)
		db, open_db_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
		if open_db_err!=nil {
			log.Println("open db file", err)
		}
		for _, v := range devs.Devices {
			add_dev := db.First(&Struct.SshStruct{},"ipaddr",v.Ipaddr)
			if (add_dev.RowsAffected) == 1{
				fail_list = append(fail_list, map[string]string{v.Ipaddr:"exist devices skip add !"})
				log.Println("device "+v.Ipaddr+" exist ! skip task")
			}else {
				db.Create(&v)
				succ_list = append(succ_list, map[string]string{v.Ipaddr:"add succ !"})
				log.Println( "add "+v.Ipaddr+" succ !")
			}
		}
		context.JSON(http.StatusOK, map[string]interface{}{"fail_list":fail_list,"succ_list":succ_list})
	}
}
