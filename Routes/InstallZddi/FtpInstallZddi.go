package InstallZddi

import (
	"TestPlatform/Fun/InstallZddi"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func FtpInstallZddi(context *gin.Context) {
	ftp_task := Struct.FtpTask{}
	if err := context.ShouldBind(&ftp_task); err != nil {
		logrus.Debug(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	} else {
		task_exe_info,_ := InstallZddi.InstallFtpZddiTask(&ftp_task)
		context.JSON(http.StatusOK, map[string]interface{}{"exe_info":task_exe_info})
	}
}