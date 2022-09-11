package DialingTest

import (
	"TestPlatform/Fun/DialingTest"
	"TestPlatform/Struct"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateTask(context *gin.Context) {
	dialing_task := Struct.DialingTestTask{}
	if err := context.ShouldBind(&dialing_task); err != nil {
		log.Println(err)
		context.SecureJSON(http.StatusInternalServerError, err)
	}else {
		if dialing_task.Devices2 == ""{
			context.File(DialingTest.StartAlsoTask(dialing_task))
		}else {
			context.File(DialingTest.StartDoubleTask(dialing_task))
		}
	}
}