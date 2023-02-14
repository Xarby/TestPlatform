package Util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func GetTxtFileMsg(path string) (string, error) {
	fmt.Println(path)
	if _, open_file_err := os.Stat(path); open_file_err != nil {
		return "", errors.New("open file " + path + " fail ")
	}
	read_context, read_err := ioutil.ReadFile(path)
	if read_err != nil {
		return "", errors.New("read file " + path + " fail ,err msg:" + read_err.Error())
	}
	// 将字节流转换为字符串
	return string(read_context), nil
}
