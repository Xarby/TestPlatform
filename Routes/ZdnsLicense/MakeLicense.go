package ZdnsLicense

import (
	"TestPlatform/Fun/ZdnsLicense"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MakeLicnese(context *gin.Context)  {
	devices := Struct.MakeLicenseStruct{}
	if err := context.ShouldBind(&devices); err != nil {
		context.SecureJSON(http.StatusOK, err)
	}else {
		if task_exe_info,task_exe_error := ZdnsLicense.ZdnsMakeLicense(devices);task_exe_error!=nil{
			context.JSON(http.StatusOK, map[string]interface{}{"err_info":task_exe_error.Error()})
		}else {
			context.JSON(http.StatusOK, map[string]interface{}{"exe_info":task_exe_info})
		}
	}
}