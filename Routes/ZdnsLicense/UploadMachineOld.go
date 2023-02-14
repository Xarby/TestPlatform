package ZdnsLicense

import (
	"TestPlatform/Const"
	"TestPlatform/Fun/ZdnsLicense"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)



func UploadMachineLicense(context *gin.Context)  {
	logrus := Util.CreateLogger(Const.RecoveryLogPath,Const.RecoveryLogPath+ "/machine_license.log")
	machineStruct := Struct.UpLoadMachineStruct{}
	context.Header("contentType","application/json")
	context.Header("Access-Control-Allow-Origin","*")
	context.Header("Access-Control-Allow-Methods","POST,GET,OPTION")
	context.Header("Access-Control-Max-Age","3600")
	context.Header("Access-Control-Allow-Headers","x-requested-with,Authorization,token, content-type")
	fmt.Println(context.Request.Header.Get("contentType"))
	fmt.Println(context.Request.Header.Get("Access-Control-Allow-Origin"))
	if err := context.BindJSON(&machineStruct); err != nil {
		fmt.Println(err.Error())
		logrus.Debug(err)
		context.SecureJSON(http.StatusOK, err)
	}else {
		if EncodeStr, taskExeError := ZdnsLicense.ZdnsUploadMachineLicense(machineStruct); taskExeError !=nil{
			context.JSON(http.StatusOK, map[string]interface{}{"Error_info": taskExeError})
		}else {
			context.JSON(http.StatusOK,map[string]interface{}{"EncodeStr_info": EncodeStr})
		}
	}
}