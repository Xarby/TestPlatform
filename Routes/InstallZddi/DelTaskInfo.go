package InstallZddi

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func DelTaskInfo(context *gin.Context) {
	installZddiTask := Struct.InstallZddiTaskStruct{}
	if err := context.ShouldBind(&installZddiTask); err == nil {
		del_sql_result_msg, del_sql_result_err := installZddiTask.DelFindTask()
		_, log_err := os.Stat(Const.ZddiTaskLogPath + "/" + installZddiTask.TaskName)
		if del_sql_result_err !=nil{
			context.SecureJSON(http.StatusOK, del_sql_result_err.Error())
		}else {
			if log_err!= nil{

				context.SecureJSON(http.StatusOK, del_sql_result_msg+"and not find log menu")
			}else {
				os.RemoveAll(Const.ZddiTaskLogPath + "/" + installZddiTask.TaskName)
				context.SecureJSON(http.StatusOK, del_sql_result_msg+"and del task log succ ...")
			}
		}

	}else {
		fmt.Println(err)
		logrus.Debug(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	}
}
