package Util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func TryConn(ipaddr string, port int) (bool, error) {
	conn, connerr := net.Dial("tcp", ipaddr+":"+strconv.Itoa(port))
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

func PostRequests(method string,url string, body []byte) error {

	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))

	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()
	result_body, _ := ioutil.ReadAll(response.Body)
	log.Println(string(result_body))
	return nil
}
