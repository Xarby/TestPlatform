package DialingTest

import (
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"github.com/sirupsen/logrus"
	"log"
	"reflect"
	"sync"
)

func StartAlsoTask(task Struct.DialingTestTask) string {
	wg := sync.WaitGroup{}
	var mutex sync.Mutex
	dig_result := Struct.AlsoDialingTestResult{}
	log.Println("start sola dialing mode ...")

	for _, domainStruct := range task.BatchDomain {
		wg.Add(1)
		go func(temp_task Struct.DialingTestTask, temp_domain Struct.DomainStruct) {
			if temp_domain.Type == "" {
				temp_domain.Type = "A"
			}
			mutex.Lock()
			dig_result.AlsoDialingTest = append(dig_result.AlsoDialingTest, Struct.AlsoDialingTestStruct{
				Doamin: temp_domain.Domain,
				Type:   temp_domain.Type,
				Result: Util.DigDoamin(temp_task.Devices1, temp_domain.Domain, temp_domain.Type)})
			mutex.Unlock()
			wg.Done()
		}(task, domainStruct)
	}
	wg.Wait()

	return dig_result.AlsoToExcel()
	//返回json
	//return dig_result

}

func StartDoubleTask(task Struct.DialingTestTask) string {
	//校验是否存在相同元素
	wg := sync.WaitGroup{}
	var mutex sync.Mutex
	dig_result := Struct.DoubleDialingTestResult{}
	logrus.Info("start double dialing mode ...")

	for _, domainStruct := range task.BatchDomain {
		wg.Add(1)
		go func(temp_task Struct.DialingTestTask, temp_domain Struct.DomainStruct) {
			if temp_domain.Type == "" {
				temp_domain.Type = "A"
			}
			dev1_result := Util.DigDoamin(temp_task.Devices1, temp_domain.Domain, temp_domain.Type)
			dev2_result := Util.DigDoamin(temp_task.Devices2, temp_domain.Domain, temp_domain.Type)
			var contrast_result string
			if reflect.DeepEqual(dev1_result, dev2_result){
				contrast_result = "相同"
			}else {
				contrast_result = "不相同"
			}
				mutex.Lock()
			dig_result.AlsoDialingTest = append(dig_result.AlsoDialingTest, Struct.DoubleDialingTestStruct{
				Doamin:         temp_domain.Domain,
				Type:           temp_domain.Type,
				Devices1Result: Util.DigDoamin(temp_task.Devices1, temp_domain.Domain, temp_domain.Type),
				Devices2Result: Util.DigDoamin(temp_task.Devices2, temp_domain.Domain, temp_domain.Type),
				ContrastResult:contrast_result,
			})
			mutex.Unlock()
			wg.Done()
		}(task, domainStruct)
	}
	wg.Wait()
	logrus.Info("double dialing task execute end ...")
	return dig_result.DoubleToExcel(task.Devices1,task.Devices2)
	//返回json
	//return dig_result
}
