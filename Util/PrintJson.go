package Util

import (
	"encoding/json"
	"log"
)
func PrintJson(Point any)  {
	if task_info,err := json.MarshalIndent(Point,"","\t"); err != nil{
		log.Println(err)
	}else {
		log.Println(string(task_info))
	}
}
