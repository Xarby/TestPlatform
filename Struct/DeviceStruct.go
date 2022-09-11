package Struct

//用来删除数据库中的设备IP集合
type DelDevList struct {
	DevList []string `json:"dev_list"`
}


//获取的设备信息
type DevInfoStruct struct {
	Ipaddr string `json:"ipaddr" gorm:"primaryKey"`
	CardName string `json:"card_name"`
	MacAddr string `json:"mac_addr"`
	OptVersion string `json:"opt_version"`
	RpmInfo string `json:"rpm_info"`
	MemSize string `json:"mem_size"`
	CpuNum string `json:"cpu_num"`
	CpuName string `json:"cpu_name"`
	SkuNum string `json:"sku_num"`
	SnNum string `json:"sn_num"`
	DiskSize string `json:"disk_size"`
	DiskUse string `json:"disk_use"`
	Status string `json:"status"`
}
//获取的设备信息的集合
type BatchDevInfoStruct struct {
	BatchDevInfo []DevInfoStruct `json:"batch_dev_info"`
}