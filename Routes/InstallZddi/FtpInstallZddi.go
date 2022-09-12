package InstallZddi

import (
	"TestPlatform/Fun/InstallZddi"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func FtpInstallZddi(context *gin.Context) {
	ftp_task := Struct.FtpTask{}
	if err := context.ShouldBind(&ftp_task); err != nil {
		log.Println(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	} else {
		context.JSON(http.StatusOK, InstallZddi.InstallFtpZddiTask(&ftp_task))
	}
}
