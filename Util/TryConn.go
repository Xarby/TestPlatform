package Util

import (
	"TestPlatform/Const"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

func TryConn(ipaddr string, port int,timeout int) (bool, error) {
	dial := net.Dialer{Timeout: time.Second*time.Duration(timeout)}
	conn, connerr := dial.Dial("tcp", ipaddr+":"+strconv.Itoa(port))
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if connerr != nil {
		return false, connerr
	} else {
		return true, nil
	}
}

func PostRequests(method string,url string, body []byte) (error,string) {

	client := &http.Client{
		Timeout: time.Second * 20,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))

	if err != nil {
		return fmt.Errorf("Got error %s", err.Error()),""
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(Const.ZddiApiDefaultUser, Const.ZddiApiDefaultPasswd)
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error()),""
	}
	defer response.Body.Close()
	result_body, _ := ioutil.ReadAll(response.Body)
	logrus.Info(string(result_body))
	return nil,string(result_body)
}
