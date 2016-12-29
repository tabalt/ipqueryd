#! /bin/bash

WORKSPACE=$(cd $(dirname $0)/; pwd)
cd $WORKSPACE

### config ###

APPS="ipqueryd"
INIT_DIRS="logs tmp"


### common function ###

function _make_pid_file() {
    app=$1
    echo "$WORKSPACE/tmp/$app.pid"
}

function _get_pid() {
    app=$1
    pid_file=`_make_pid_file $app`
    if [ -f $pid_file ]; then 
        pid=`cat $pid_file`
        kill -0 $pid >/dev/null 2>&1
        if [ $? = 0 ]; then 
            echo $pid
        else
            rm -f $pid_file
            echo 0
        fi
    else
        echo 0
    fi
}


### action ###

function useage() {
    echo "Useage: $0 \$cmd [\$app]";
    exit 1;
}

function init() {
    app=$1
    for dir in $INIT_DIRS ; do
        init_path=$WORKSPACE/$dir
        if [ ! -d $init_path ]; then 
            mkdir -p $init_path && chmod -R 777 $init_path
        fi
    done
}

function build() {
    app=$1
    cd $WORKSPACE/src
    ./go.sh install main/$app
    cd $WORKSPACE
}

function status() {
    app=$1
    pid=`_get_pid $app`
    if [ "$pid" = 0 ]; then 
        echo "$app not running."
    else
        echo "$app with pid $pid is running."
    fi
}

function start() {
    app=$1

    pid=`_get_pid $app`
    if [ $pid -gt 0 ]; then 
        echo "$app with pid $pid has running!"
        return 1
    fi

    ($WORKSPACE/bin/$app -c $WORKSPACE/conf/$app.json >> $WORKSPACE/logs/$app.log 2>&1 &)
}

function stop() {
    app=$1

    pid=`_get_pid $app`
    if [ $pid -eq 0 ]; then 
        echo "$app not running!"
        return 1
    fi

    kill -9 $pid
    rm -f `_make_pid_file $app`
}

function restart() {
    app=$1
    stop $app && start $app
}

function run() {
    app=$1
    init $app
    build $app
    start $app
}

### main ###

APP_LIST=($APPS)
APP=${APP_LIST[0]}
for (( i = 0; i < ${#APP_LIST[@]}; i++ )); do
    app=${APP_LIST[$i]}
    if [ $app = "$2" ]; then
        APP=$app
    fi
done

case $1 in
    "init")
        init $APP
    ;;
    "build")
        build $APP
        ;;
    "status")
        status $APP
        ;;
    "start")
        start $APP
        ;;
    "stop")
        stop $APP
        ;;
    "restart")
        restart $APP
        ;;
    "run")
        run $APP
        ;;
    *)
        useage
        ;;
esac