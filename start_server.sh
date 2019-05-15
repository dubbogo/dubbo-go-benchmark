#!/usr/bin/env bash

#p代表协议
while getopts ":p:" opt
do
    case $opt in
        p)
        popt=$OPTARG
        ;;
        ?)
        echo "未知参数"
        exit 1;;
    esac
done

#判断 p参数是否输入，若未输入，则默认dubbo

if  [ ! -n "$popt" ] ;then
    popt="dubbo"
fi

sh stop.sh
export -n CONF_CONSUMER_FILE_PATH
export CONF_PROVIDER_FILE_PATH=$PWD/$popt/server.yml
export APP_LOG_CONF_FILE=$PWD/$popt/log.xml

cd ./$popt
go build .
./$popt &
echo "进程ID:"$!
echo $!>./pid

