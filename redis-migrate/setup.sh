#!/bin/bash
# -*- coding: utf-8 -*-
 
start() {
    cd data
    for i in `seq 6379 6380`; 
    do
        redis-server --daemonize yes --port $i --pidfile $i.pid --dbfilename $i.rdb
    done
}

stop() {
    cd data
    for i in `ls *.pid`;
    do
        kill `cat $i`
    done
}

clean() {
    rm *.pid *.rdb
}


case $1 in
start)
    start
    ;;
stop)
    stop
    ;;
clean)
    clean
    ;;
*)
    echo "./$0 start|stop|clean"
esac
