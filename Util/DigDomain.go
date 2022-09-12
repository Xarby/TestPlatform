package Util

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

func DigDoamin(ipaddr string, domain string, domain_type string) ([]string) {
	var out []byte
	var err error
	if runtime.GOOS == "linux" {
		fmt.Println("/usr/bin/dig"+"@"+ipaddr+ domain+ domain_type+ "+short")
		out, err = exec.Command("dig", "@"+ipaddr, domain, domain_type, "+short").Output()
	} else if runtime.GOOS == "windows" {
		out, err = exec.Command("cmd.exe","/c", "dig", "@"+ipaddr, domain, domain_type, "+short").Output()
	}
	exec_result := bat_result(string(out))
	sort.Strings(exec_result)
	if err != nil {
		log.Println(err)
		return []string{err.Error()}
	} else {
		return exec_result
	}
}

func bat_result(domain string) []string {
	var result []string
	if domain != "" && len(domain) > 2 {
		for _, v := range strings.Split(domain,"\n") {
			result = append(result,strings.Replace(strings.Replace(v,"\n","",-1),"\t","",-1))
		}
		return result
	} else {
		return []string{"none"}
	}
}
