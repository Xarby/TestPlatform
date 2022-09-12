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
		task_exe_info,_ := InstallZddi.InstallScpZddiTask(&scp_task)
		context.JSON(http.StatusOK, map[string]interface{}{"exe_info":task_exe_info})
	}
}