#!/bin/bash
while [ 1 ]
do
    sleep 1
    basedir=`ps -ef | grep main | grep -v grep`
    if [ -z "$basedir" ] ;then
		nohup go run main.go &
	fi
done
