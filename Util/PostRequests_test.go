package Util

import (
	"fmt"
	"testing"
)

func TestTryConn(t *testing.T) {
	//fmt.Println("-----------")
	//fmt.Println(TryConn("10.1.104.55",443))
	body := []byte(`{
	"name": "slave",
	"ip": "10.1.107.51",
	"username": "admin",
	"password": "admincns",
	"role": "slave",
	"group": "local",
	"is_extend":"no"}`)
	fmt.Println(string(body))
	fmt.Println(PostRequests("https://10.1.107.54:20120/groups/local/members", body))
	fmt.Println("-----------")
}
