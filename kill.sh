kill -9 ` ps -ef | grep start | grep -v grep | awk '{print $2}' `
killall main