# TestPlatform

###测试组工具平台
仅供学习使用，本项目采用go-gin框架
实现功能:
    1.一键集群化部署环境(支持SCP/FTP方式获取源rpm包文件,部署+激活+组建集群)
    2.使用rsync对完整的(centos6/7,kylin10,openEuler22)进行备份/还原
    3.支持对多设备进行获取CPU/磁盘/内存/SN码等等设备信息获取并生成表格


部署DDI接口
[GIN-debug] POST   /install_zddi/ftp         --> TestPlatform/Routes/InstallZddi.FtpInstallZddi (3 handlers)
[GIN-debug] POST   /install_zddi/scp         --> TestPlatform/Routes/InstallZddi.ScpInstallZddi (3 handlers)
[GIN-debug] GET    /install_zddi/get_all_task_info --> TestPlatform/Routes/InstallZddi.GetAllTaskInfo (3 handlers)
[GIN-debug] GET    /install_zddi/get_task_info --> TestPlatform/Routes/InstallZddi.GetTaskLog (3 handlers)
[GIN-debug] DELETE /install_zddi/del_task_info --> TestPlatform/Routes/InstallZddi.DelTaskInfo (3 handlers)
[GIN-debug] DELETE /install_zddi/del_all_task_info --> TestPlatform/Routes/InstallZddi.DelAllTaskInfo (3 handlers)
[GIN-debug] GET    /device/get_dev_info      --> TestPlatform/Routes/Devices.GetDeviceInfo (3 handlers)
[GIN-debug] GET    /device/get_dev           --> TestPlatform/Routes/Devices.GetDevice (3 handlers)
[GIN-debug] POST   /device/add_dev           --> TestPlatform/Routes/Devices.AddDevice (3 handlers)
[GIN-debug] PUT    /device/put_dev           --> TestPlatform/Routes/Devices.PutDevice (3 handlers)
[GIN-debug] DELETE /device/del_dev           --> TestPlatform/Routes/Devices.DelDevice (3 handlers)
[GIN-debug] POST   /dialing_test/create_task --> TestPlatform/Routes/DialingTest.CreateTask (3 handlers)
[GIN-debug] POST   /zdns_recovery/backup     --> TestPlatform/Routes/ZdnsBackup.Backup (3 handlers)
[GIN-debug] POST   /zdns_recovery/recovery   --> TestPlatform/Routes/ZdnsBackup.Recovery (3 handlers)
[GIN-debug] POST   /zdns_recovery/show_statu --> TestPlatform/Routes/ZdnsBackup.ShowStatu (3 handlers)
[GIN-debug] POST   /license/make_license     --> TestPlatform/Routes/ZdnsLicense.MakeLicnese (3 handlers)
[GIN-debug] POST   /license/upload_machine_license --> TestPlatform/Routes/ZdnsLicense.UploadMachineLicense (3 handlers)