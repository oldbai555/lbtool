#!/bin/sh

REDIS_PATH='/c/Users/baigege/Desktop/middleware/redis'

usage() {
  echo "Usage: sh 执行脚本.sh [start]"
  exit 1
}

#启动方法
start() {
  echo "start redis"
  cd $REDIS_PATH;
  ./redis-server.exe redis.windows.conf
}

#根据输入参数，选择执行对应方法，不输入则执行使用说明
case "$1" in
"start")
  start
  ;;
*)
  usage
  ;;
esac