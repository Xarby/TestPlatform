# TestPlatform

###测试组工具平台
仅供学习使用，本项目采用go-gin框架\

实现功能:\
    1.一键集群化部署环境(支持SCP/FTP方式获取源rpm包文件,部署+激活+组建集群)\
    2.使用rsync对完整的(centos6/7,kylin10,openEuler22)进行备份/还原\
    3.支持对多设备进行获取CPU/磁盘/内存/SN码等等设备信息获取并生成表格\


集群化部署环境整体流程：
    1.校验
        （1）对版本包获取的源（FTP服务器或者SCP主机登录），并校验文件是否存在
        （2）检查所有部署的环境是否为干净环境（通过ssh登录并查询是否安装rpm）
        （3）检查输入的架构是否合法，例如两个master则算作异常
    2.每部署一个节点开启一个grouting并发执行
        需要grouting原因是因为由server 发送IO是一定的但是每个节点安装rpm是可以异步安装，提升部署效率
    3.每个grouting完成安装之后就会去检测所有服务是否正常启动，并且sync.WaitGroup异步累计器减一
    4.单所有节点基础服务安装之后 通过对master节点发送post请求纳管slave/backup 或者组HA来实现集群化。


