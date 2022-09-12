package InstallZddi

import (
	"TestPlatform/Fun/InstallZddi"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ScpInstallZddi(context *gin.Context) {
	scp_task := Struct.ScpTask{}
	if err := context.ShouldBind(&scp_task); err != nil {
		log.Println(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	}else {
		context.JSON(http.StatusOK, InstallZddi.InstallScpZddiTask(&scp_task))
	}
}