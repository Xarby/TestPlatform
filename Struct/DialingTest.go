package Struct

import (
	"TestPlatform/Const"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"time"
)

type DomainStruct struct {
	Domain string `json:"name"`
	Type   string `json:"type"`
}

//总任务
type DialingTestTask struct {
	Devices1    string         `json:"devices_1"`
	Devices2    string         `json:"devices_2"`
	BatchDomain []DomainStruct `json:"batch_domain"`
}

//单机拨测的结构体
type AlsoDialingTestStruct struct {
	Doamin string   `json:"doamin"`
	Type   string   `json:"type"`
	Result []string `json:"result"`
}

//单机拨测的结果
type AlsoDialingTestResult struct {
	AlsoDialingTest []AlsoDialingTestStruct `json:"also_dialing_test"`
	Path            string                  `json:"path"`
}

func (also_dialing_test *AlsoDialingTestResult) AlsoToExcel() string {
	file_name := Const.DialingTestFilePath + time.Now().Format("2006_01_02_15_04_05") + "_also.xlsx"
	xlsx := excelize.NewFile()
	defer func() {
		err := xlsx.SaveAs(file_name)
		if err != nil {
			fmt.Println(err)
		}
	}()
	Sheet_name := "Sheet1"
	index := xlsx.NewSheet("Sheet1")
	xlsx.SetActiveSheet(index)
	xlsx.SetSheetRow(Sheet_name, "A1", &[]interface{}{"域名", "类型", "解析结果"})
	start_flag := 2
	for _, testStruct := range also_dialing_test.AlsoDialingTest {
		xlsx.SetSheetRow(Sheet_name, "A"+strconv.Itoa(start_flag), &[]interface{}{testStruct.Doamin, testStruct.Type, testStruct.Result})
		start_flag++
	}
	xlsx.SetColWidth("Sheet1", "A", "A", 40)
	xlsx.SetColWidth("Sheet1", "C", "C", 40)
	also_dialing_test.Path = file_name
	return file_name
}

//双设备拨测的结构体
type DoubleDialingTestStruct struct {
	Doamin         string   `json:"doamin"`
	Type           string   `json:"type"`
	Devices1Result []string `json:"devices_1_result"`
	Devices2Result []string `json:"devices_2_result"`
	ContrastResult string   `json:"contrast_result"`
	Path           string   `json:"path"`
}

type DoubleDialingTestResult struct {
	AlsoDialingTest []DoubleDialingTestStruct `json:"double_dialing_test"`
	Path            string                    `json:"path"`
}

func (double_dialing_test *DoubleDialingTestResult) DoubleToExcel(dev1_ipaddr string, dev2_ipaddr string) string {
	file_name := Const.DialingTestFilePath + time.Now().Format("2006_01_02_15_04_05") + "_double.xlsx"
	xlsx := excelize.NewFile()
	defer func() {
		err := xlsx.SaveAs(file_name)
		if err != nil {
			fmt.Println(err)
		}
	}()
	Sheet_name := "Sheet1"
	index := xlsx.NewSheet("Sheet1")
	xlsx.SetActiveSheet(index)
	xlsx.SetSheetRow(Sheet_name, "A1", &[]interface{}{"域名", "类型", dev1_ipaddr + "解析结果", dev2_ipaddr + "解析结果","对比结果"})
	start_flag := 2
	for _, testStruct := range double_dialing_test.AlsoDialingTest {
		xlsx.SetSheetRow(Sheet_name, "A"+strconv.Itoa(start_flag), &[]interface{}{testStruct.Doamin, testStruct.Type, testStruct.Devices1Result, testStruct.Devices2Result, testStruct.ContrastResult})
		start_flag++
	}
	xlsx.SetColWidth("Sheet1", "A", "A", 40)
	xlsx.SetColWidth("Sheet1", "C", "D", 40)
	double_dialing_test.Path = file_name
	return file_name
}
