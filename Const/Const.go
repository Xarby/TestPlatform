package Const

const RunPort = "30120"

//业务包和tar包存放目录
const ZddiFileMenuName = "File/ZddiFile/"

//private包存放目录
const ZddiPrivateMenuName = "File/Private/"

//本地恢复工具路径
const LocalRecoverToolPath = ZddiPrivateMenuName+"sys_recovery_test.tar.gz"
//远端
const RemoteRecoverToolPath = ZddiRemoteFilePath+"sys_recovery_test.tar.gz"

//恢复工具安装路径
const RecoverToolInstallPath = ZddiRemoteFilePath+"sys_recovery/install.sh"

//设备信息目录存放目录
const DevicesInfoMenuName = "File/Devices/"

//设备设备信息文件存放路径
const DevicesInfoRootSqlPath = DevicesInfoMenuName+"Devices.db"

//put远程目录
const ZddiRemoteFilePath = "/root/"

//日志根目录
const LogPath = "Log/"

//接口日志目录
const GinLogPath = LogPath+"Gin.log"

//业务日志目录
const WorkLogPath = LogPath+"Work.log"

//拨测生成表格列表
const DialingTestFilePath  = "File/DialingTest/"

const ZddiApiDefaultUser = "admin"

const ZddiApiDefaultPasswd = "admin"
