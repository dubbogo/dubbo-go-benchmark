#!/usr/bin/env bash


# p代表协议，r代表参数大小（1等于300B），t代表间隔时间（单位ms）
while getopts ":p:r:t:" opt
do
    case $opt in
        p)
        popt=$OPTARG
        ;;
        r)
        ropt=$OPTARG
        ;;
        t)
        topt=$OPTARG
        ;;
        ?)
        echo "未知参数"
        exit 1;;
    esac
done

#判断 p、r和t参数是否输入，若未输入，则默认dubbo,2,30
if  [ ! -n "$popt" ] ;then
    popt="dubbo"
fi
if  [ ! -n "$ropt" ] ;then
    ropt=2
fi
if  [ ! -n "$topt" ] ;then
    topt=30
fi

sh stop.sh
export CONF_PROVIDER_FILE_PATH=$PWD/$popt/server/server.yml
export APP_LOG_CONF_FILE=$PWD/$popt/server/log.yml

cd ./$popt/server
go build .
./server &
echo "进程ID:"$!
echo $!>./pid
sleep 4
cd ../../
export -n CONF_PROVIDER_FILE_PATH
export CONF_CONSUMER_FILE_PATH=$PWD/$popt/continuity_testing/client.yml
export APP_LOG_CONF_FILE=$PWD/$popt/continuity_testing/log.yml

go run ./$popt/continuity_testing/*.go -t $topt  -r $ropt

