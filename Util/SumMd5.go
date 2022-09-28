package Util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func GetFileMd5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		logrus.Error("os file"+filename+" error")
		return "", err
	}
	md5 := md5.New()
	_, err = io.Copy(md5, file)
	if err != nil {
		logrus.Error("io copy error")
		return "", err
	}
	md5Str := hex.EncodeToString(md5.Sum(nil))
	return md5Str, nil
}
