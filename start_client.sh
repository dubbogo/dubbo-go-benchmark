#!/usr/bin/env bash

#c 代表压测并发数，n代表压测总数，p代表协议
while getopts ":p:n:c:" opt
do
    case $opt in
        c)
        copt=$OPTARG
        ;;
        n)
        nopt=$OPTARG
        ;;
        p)
        popt=$OPTARG
        ;;
        ?)
        echo "未知参数"
        exit 1;;
    esac
done

#判断 p、c和n参数是否输入，若未输入，则默认dubbo,1,1
if  [ ! -n "$copt" ] ;then
    copt=1
fi

if  [ ! -n "$nopt" ] ;then
    nopt=1
fi

if  [ ! -n "$popt" ] ;then
    popt="dubbo"
fi


export -n CONF_PROVIDER_FILE_PATH
export CONF_CONSUMER_FILE_PATH=$PWD/$popt/client/client.yml
export APP_LOG_CONF_FILE=$PWD/$popt/client/log.xml


go run ./$popt/client/*.go -c $copt  -n $nopt


