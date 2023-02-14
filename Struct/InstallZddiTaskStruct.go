package Struct

import (
	"TestPlatform/Const"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type InstallZddiTaskStruct struct {
	TaskName       string `json:"task_name" gorm:"primaryKey"`
	Nodes          string `json:"nodes"`
	LicenseVersion string `json:"license_version"`
	BuildPkg       string `json:"build_pkg"`
	ZddiPkg        string `json:"zddi_pkg"`
	CreateTime     string `json:"create_time"`
	StatusMsg      string `json:"status_msg"`
}

type BatchInstallZddiTaskStruct struct {
	BatchSsh []InstallZddiTaskStruct `json:"batch_tasks"`
}

func (TaskDB *InstallZddiTaskStruct) DBCreateTask(logrus logrus.Logger) (string, error) {
	ZddiInstallInfoDB, open_db_err := gorm.Open(sqlite.Open(Const.InstallZddiTaskSqlPath), &gorm.Config{})
	if open_db_err != nil {
		logrus.Debug("open db file", open_db_err)
		return "open db file " + Const.InstallZddiTaskSqlPath, open_db_err
	} else {
		_, flag, find_err := TaskDB.DBFindTask()
		if find_err != nil {
			return "find task err", find_err
		}
		if flag == false {
			ZddiInstallInfoDB.Create(TaskDB)
			return "check db succ !!!", nil
		} else {
			return "exist task " + TaskDB.TaskName, errors.New("exist task " + TaskDB.TaskName)
		}
		return TaskDB.TaskName + " insert db succ ...", nil
	}
}

func (TaskDB *InstallZddiTaskStruct) DBFindTask() (string, bool, error) {
	if len(TaskDB.TaskName) == 0 {
		return "task name is nil", false, errors.New("task name is nil")
	}
	ZddiInstallInfoDB, open_db_err := gorm.Open(sqlite.Open(Const.InstallZddiTaskSqlPath), &gorm.Config{})
	if open_db_err != nil {
		return "open db file " + Const.InstallZddiTaskSqlPath, false, open_db_err
	}
	if ZddiInstallInfoDB.First(&InstallZddiTaskStruct{}, "task_name", TaskDB.TaskName).RowsAffected == 1 {
		return "exist task " + TaskDB.TaskName, true, nil
	} else {
		return "not exist task " + TaskDB.TaskName, false, nil
	}
}

func (TaskDB *InstallZddiTaskStruct) DelFindTask() (string, error) {
	_, flag, find_err := TaskDB.DBFindTask()
	if find_err != nil {
		return "find task err", find_err
	}
	if flag == true {
		ZddiInstallInfoDB, open_db_err := gorm.Open(sqlite.Open(Const.InstallZddiTaskSqlPath), &gorm.Config{})
		if open_db_err != nil {
			return "open db file " + Const.InstallZddiTaskSqlPath, open_db_err
		}
		del_task := ZddiInstallInfoDB.First(&InstallZddiTaskStruct{}, "task_name", TaskDB.TaskName)
		del_task.Delete(1)
		return "sql del " + TaskDB.TaskName + " succ", nil
	} else {
		return "sql not find task " + TaskDB.TaskName, errors.New("sql not find task " + TaskDB.TaskName)
	}
}

func (TaskDB *InstallZddiTaskStruct) DelAllTask() (string, error) {
	ZddiInstallInfoDB, open_db_err := gorm.Open(sqlite.Open(Const.InstallZddiTaskSqlPath), &gorm.Config{})
	if open_db_err != nil {
		return "open db file " + Const.InstallZddiTaskSqlPath, open_db_err
	} else {
		ZddiInstallInfoDB.Where("1 = 1").Delete(&InstallZddiTaskStruct{})
		return "sql del All task succ", nil
	}

}

func (TaskDB *InstallZddiTaskStruct) UpdateTaskMsg(msg string) (string, error) {
	_, flag, find_err := TaskDB.DBFindTask()
	if find_err != nil {
		return "find task err", find_err
	}
	if flag == true {
		ZddiInstallInfoDB, open_db_err := gorm.Open(sqlite.Open(Const.InstallZddiTaskSqlPath), &gorm.Config{})
		if open_db_err != nil {
			return "open db file " + Const.InstallZddiTaskSqlPath, open_db_err
		}
		del_task := ZddiInstallInfoDB.First(&InstallZddiTaskStruct{}, "task_name", TaskDB.TaskName)
		TaskDB.StatusMsg = msg
		del_task.Updates(&TaskDB)
		return "sql Update " + TaskDB.TaskName + " succ", nil
	} else {
		return "sql not find task " + TaskDB.TaskName, errors.New("sql not find task " + TaskDB.TaskName)
	}
}
