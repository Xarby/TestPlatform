package Fun

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

//func GetDevicesInfo() Struct.DevsInfo {
func GetDevicesInfo() string {
	//wg := sync.WaitGroup{}
	//var mutex sync.Mutex
	db, open_sqlite_err := gorm.Open(sqlite.Open(Const.DevicesInfoRootSqlPath), &gorm.Config{})
	if open_sqlite_err != nil {
		log.Println(open_sqlite_err)
	}


	db.Find(&Struct.SshStruct{}, " ipaddr = '10.1.111.64'")
	fmt.Println(db.Find(&Struct.SshStruct{}, "ipaddr", "10.1.111.63").RowsAffected)
	fmt.Println(db.Find(&Struct.SshStruct{}, "ipaddr", "10.1.111.64").RowsAffected)

	//log.Println("start  get devices info")
	//jsonFile, read_json_err := os.Open(Const.DevicesInfoRootFileName)
	//if read_json_err != nil {
	//	log.Println("read json file fail", read_json_err)
	//}
	//defer jsonFile.Close()
	//
	//jsonData, err := ioutil.ReadAll(jsonFile)
	//if err != nil {
	//	log.Println("json change byte fail", err)
	//}

	//var devs Struct.Devs

	//var devinfos Struct.DevsInfo
	//json.Unmarshal(jsonData, &devs)
	//for _, dev := range devs.Devices {
	//	wg.Add(1)
	//	go func(sshStruct Struct.SshStruct) {
	//		getInfo, get_err := sshStruct.GetDevInfo()
	//		log.Println(getInfo)
	//		if get_err!=nil {
	//			log.Println(get_err)
	//		}
	//		mutex.Lock()
	//		devinfos.DevInfos = append(devinfos.DevInfos, getInfo)
	//		mutex.Unlock()
	//		wg.Done()
	//	}(dev)
	//}
	//wg.Wait()

	//return
	return ""
}
