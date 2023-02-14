package Const

import "github.com/sirupsen/logrus"

const RunPort = "30120"

//业务包和tar包存放目录
const ZddiFileMenuName = "File/ZddiFile/"

//private包存放目录
const ZddiPrivateMenuName = "File/Private/"

//源码路径
const ZddiPrivateNewPath = ZddiPrivateMenuName+"new/"
const ZddiPrivateOldPath = ZddiPrivateMenuName+"old/"

//制作licnese
	//压缩包文件目录(本地指定ip方式)
const ZddiPrivateTarGzPath = ZddiPrivateMenuName+"License/"
const TempLicensePath = ZddiPrivateMenuName+"tempLicense/"
const MachineLicensePath = TempLicensePath+"machine/"
const RemoteLicnesePath = "/etc/"
	//上传machine.info方式(UploadMachine)
const TempUploadMachineLicnesePath = TempLicensePath+"Upload_Machine_license/"
const TempUploadMachineLicneseFileName = "machine.info"


//本地和原创目录
const LocalZddiPrivateTarGzNew = ZddiPrivateTarGzPath+"PrivateNew.tar.gz"
const RemoteZddiPrivateTarGzNew = ZddiRemoteFilePath+"PrivateNew.tar.gz"
const LocalZddiPrivateTarGzOld = ZddiPrivateTarGzPath+"PrivateOld.tar.gz"
const RemoteZddiPrivateTarGzOld = ZddiRemoteFilePath+"PrivateOld.tar.gz"

//其他文件存放目录
const ZddiOtherMenuName = "File/OtherFile/"

//设备信息目录存放目录
const DBMenuName = "File/DB/"

//设备信息文件存放路径
const DevicesInfoRootSqlPath = DBMenuName + "Devices.db"

//部署ddi任务存放路径
const InstallZddiTaskSqlPath = DBMenuName + "InstallZddi.db"

//put远程目录
const ZddiRemoteFilePath = "/root/"


//日志根目录
const LogPath = "Log/"

//接口日志目录
const GinLogPath = LogPath + "Gin.log"

//业务日志目录
const WorkLogPath = LogPath + "Work.log"

//安装DDI日志目录
const ZddiTaskLogPath = LogPath + "/ZddiTaskLog/"

//扫描设备日志目录
const DevicesLogPath = LogPath + "/DevicesLog/"

//备份还原日志目录
const RecoveryLogPath = LogPath + "/RecoveryLog/"

//制作license日志目录
const UploadMachineLicenseLogPath = LogPath + "/UploadMachineLicenseLogPath/"
//拨测生成表格列表
const DialingTestFilePath = "File/DialingTest/"

//下发参数默认的登录账号和密码
const ZddiApiDefaultUser = "admin"
const ZddiApiDefaultPasswd = "admin"

//下发参数默认的登录账号和密码
const ZddiAddNodeUser = "admin"
const ZddiAddNodePasswd = "admincns"

const SshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDIIJDt4g0l8p3gnJ70ckZUW/1yj8jXkcVMa8quxqBbySo4OAt9WIGn1VCfIyLe/3qkWOGbZ0yGp74GsYaON5t7dItCls/4duAOoKSxA8NjIsWWJx/BhMsi3HEJST6ZpdkVx5KJDDWd5bM7+1tD9webwIfIF4YfRr67zzH3O0AxPP9dYNMbFS/FTUdHrzdvFfae2UYX9nUCEr2I+xVPF+d0d2/iZYJikTwcIXy1LLbbgvp09Z0Avy+3CRseGTyIvFNWIa1k5nHDEdp8k5v4W7zoLLhjPcm06jQP/p0LInN+bEv2gPOM6cb3X6nEVv6CeUs/y5Vrsr1tSQPfYVxF21UEOCUK3xMrR95ELs1dGXYo4IPSGRXWkBmTEDpnNip/wvmGQKqTbJQj1ZVR0+KOU01mnupnuaOEJnrjwXkwwROyf2Abl79ebLFeNzIIET/ICBUT6XkbCAUqv6tgmSEZtITYNwJI+O3SkNFv8tV2+2uld/Sb11XpgDS9bD1sVjxXggs= Xarby@DESKTOP-Q3FA3PK\n"

//备份系统

const DebugLevel=logrus.DebugLevel
const BackupWorkDir = "/zdns_backup"

const LocalShellFile = ZddiOtherMenuName+"Shell/"
const LocalBackupShellFile = LocalShellFile+"zdns-recovery-tool"

const RemoteShellDir = "/usr/bin/"
const RemoteBackupToolShellFile = RemoteShellDir+"zdns-recovery-tool"
const RsyncFilePath = ZddiOtherMenuName+"Rsync/"
const RsyncFilePathCentos7X86 = "rsync-3.1.2-11.el7_9.x86_64.rpm"
const RsyncFilePathKylin10ARM = "rsync-3.1.3-6.ky10.aarch64.rpm"
const RsyncFilePathopenEulerX86 = "rsync-3.2.3-4.oe2203.x86_64.rpm"
const RsyncFilePathCentos6X86 = "rsync-3.0.6-12.el6.x86_64.rpm"
