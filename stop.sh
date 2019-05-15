#!/usr/bin/env bash
cd ./dubbo/
if [  -f "pid" ]; then
pid=$(cat pid)
kill -9 $pid
fi
rm -f pid
rm -f ./dubbo

cd ../jsonrpc/
if [  -f "pid" ]; then
pid=$(cat pid)
kill -9 $pid
fi
rm -f pid
rm -f ./jsonrpc
