package InstallZddi

import (
	"TestPlatform/Const"
	"TestPlatform/Util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTaskLog(context *gin.Context) {
	//判断传来的JSON是否异常
	task_name := context.Query("task_name")
	node := context.Query("node")
	if task_name == "" {
		context.JSON(http.StatusOK, map[string]string{"msg": "no task name !!!"})
		return
	}
	if node == "" {
		if read_context, read_err := Util.GetTxtFileMsg(Const.ZddiTaskLogPath  + task_name + "/" + task_name + ".log"); read_err == nil {
			context.JSON(http.StatusOK, map[string]string{"msg": read_context})
		} else {
			context.JSON(http.StatusOK, map[string]string{"msg": read_err.Error()})
		}
	} else {
		if read_context, read_err := Util.GetTxtFileMsg(Const.ZddiTaskLogPath  + task_name + "/" + node + ".log"); read_err == nil {
			context.JSON(http.StatusOK, map[string]string{"msg": read_context})
		} else {
			context.JSON(http.StatusOK, map[string]string{"msg": read_err.Error()})
		}
	}
}
