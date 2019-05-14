#!/usr/bin/env bash
if [  -f "pid" ]; then
pid=$(cat pid)
kill -9 $pid
fi
rm -f pid
