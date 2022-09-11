package Util

import "net"

func CheckIP(ipaddr string) bool {
	if address := net.ParseIP(ipaddr); address ==nil{
		return false
	}else {
		return true
	}
}
