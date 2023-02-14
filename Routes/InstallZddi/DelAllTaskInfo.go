package InstallZddi

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func DelAllTaskInfo(context *gin.Context) {
	installZddiTask := Struct.InstallZddiTaskStruct{}
	del_sql_result_msg, del_sql_result_err := installZddiTask.DelAllTask()
	_, log_err := os.Stat(Const.ZddiTaskLogPath)
	if del_sql_result_err != nil {
		context.SecureJSON(http.StatusOK, del_sql_result_err.Error())
	} else {
		if log_err != nil {
			context.SecureJSON(http.StatusOK, del_sql_result_msg+"and not find log menu")
		} else {
			os.RemoveAll(Const.ZddiTaskLogPath)
			context.SecureJSON(http.StatusOK, del_sql_result_msg+"and del task log succ ...")
		}
	}
}
