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

func PutDevice(context *gin.Context) {
	devs := Struct.Devs{}
	succ_list := []map[string]string{}
	fail_list := []map[string]string{}
	if err := context.ShouldBind(&devs); err != nil {
		log.Println(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	} else {
		db, open_db_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
		if open_db_err!=nil {
			log.Println("open db file", err)
		}
		db.AutoMigrate(Struct.SshStruct{})
		for _, v := range devs.Devices {
			put_dev := db.First(&Struct.SshStruct{},"ipaddr",v.Ipaddr)
			if (put_dev.RowsAffected) == 1{
				put_dev.Updates(&v)
				succ_list = append(succ_list, map[string]string{v.Ipaddr:"put succ !"})
			}else {

				fail_list = append(fail_list, map[string]string{v.Ipaddr:"exist devices skip add !"})
			}
		}
		context.JSON(http.StatusOK, map[string]interface{}{"fail_list":fail_list,"succ_list":succ_list})
	}
}
