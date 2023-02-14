package ZdnsBackup

import (
	"TestPlatform/Const"
	"TestPlatform/Fun/ZdnsBackup"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Recovery(context *gin.Context)  {
	logrus := Util.CreateLogger(Const.RecoveryLogPath,Const.RecoveryLogPath+ "/recovery.log")
	devices := Struct.SshStruct{}
	if err := context.ShouldBind(&devices); err != nil {
		logrus.Debug(err)
		context.SecureJSON(http.StatusOK, err)
	}else {
		fmt.Println()
		if taskExeInfo, taskExeError := ZdnsBackup.CheckRecovery(devices); taskExeError !=nil{
			context.JSON(http.StatusOK, map[string]interface{}{"err_info": taskExeError.Error()})
		}else {
			context.JSON(http.StatusOK, map[string]interface{}{"exe_info": taskExeInfo})
		}
	}
}