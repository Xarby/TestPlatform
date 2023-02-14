package Util

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

func BitChangeGB(BitNum string)  (string,error){
	atoi, err := strconv.Atoi(BitNum)
	if err == nil {
		return strconv.Itoa(atoi/1024/1024)+" GB",nil
	}else {
		logrus.Error(err)
		return "",err
	}
}


func KbChangeGB(BitNum string)  (string,error){
	atoi, err := strconv.Atoi(BitNum)
	if err == nil {
		return strconv.Itoa(atoi/1024)+" GB",nil
	}else {
		logrus.Error(err)
		return "",err
	}
}