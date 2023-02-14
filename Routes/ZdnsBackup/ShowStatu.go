package ZdnsBackup

import (
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func ShowStatu(context *gin.Context)  {
	devices := Struct.SshStruct{}
	if err := context.ShouldBind(&devices); err != nil {
		context.SecureJSON(http.StatusOK, err)
	}else {
		var backup_file_info string
		var process_info string
		var rpm_info string
		if file,process,rpm,task_exe_error := devices.ShowStatu();task_exe_error!=nil{
			context.JSON(http.StatusOK, map[string]interface{}{"err_info":task_exe_error.Error()})
		}else {
			logrus.Debug(file,process,rpm,task_exe_error)
			if rpm == 0 {
				rpm_info = "Not Install Zddi Build Pkg ..."
			}else if rpm == 1{
				rpm_info = "Exist Install Zddi Build Pkg ..."
			}

			if process == 0{
				process_info = "Not Execute Rsync Task ..."
			}else if process == 1{
				process_info = "Exist Execute Rsync Task ,Please Wait ..."
			}

			if file == 0{
				backup_file_info = "Not Find Backup File , Please Backup ..."
			}else if file == 1{
				backup_file_info = "Exist Execute Backup Task , Please Wait ..."
			}else if file == 2{
				backup_file_info = "Execute Backup Task Done ..."
			}
			context.JSON(http.StatusOK, map[string]interface{}{"backup_info":backup_file_info,"rsync_process_info":process_info,"zddi_rpm_info":rpm_info})
		}
	}
}
