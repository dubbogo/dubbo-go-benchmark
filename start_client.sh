#!/usr/bin/env bash

# c代表压测并发数，n代表压测总数，p代表协议，r代表参数大小（1等于300B）
while getopts ":p:n:c:r:" opt
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
        r)
        ropt=$OPTARG
        ;;
        ?)
        echo "未知参数"
        exit 1;;
    esac
done

#判断 p、c、n和r参数是否输入，若未输入，则默认dubbo,1,1,2
if  [ ! -n "$copt" ] ;then
    copt=1
fi

if  [ ! -n "$nopt" ] ;then
    nopt=1
fi

if  [ ! -n "$popt" ] ;then
    popt="dubbo"
fi

if  [ ! -n "$ropt" ] ;then
    ropt=2
fi


export -n CONF_PROVIDER_FILE_PATH
export CONF_CONSUMER_FILE_PATH=$PWD/$popt/client/client.yml
export APP_LOG_CONF_FILE=$PWD/$popt/client/log.yml


go run ./$popt/client/*.go -c $copt  -n $nopt -r $ropt


