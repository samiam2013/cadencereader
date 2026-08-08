#!/bin/bash

OS=$(uname)
if [ $OS != 'Darwin' ]; then 
    echo "not running on macos, is this prod?"; 
    exit 1; 
fi

echo "compilation time:";
TIMEFORMAT=%R;
time go build .. ;
if [ $? -ne 0 ]; then 
    echo "build failed, check error above";
    exit 1;
fi

set -a
source .env
set +a

echo "starting broswer";
(sleep 1; open http://localhost:8080) &
./cadencereader
