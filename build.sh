#!/bin/bash

# Auto-detect root to set GOPATH
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo = Initialize root = $ROOT
export GOPATH=$ROOT
cd src/hue-alarm
glide i

echo = Compile
cd ../..

# Change GOOS and GOARCH to cross compile for other platforms!
env GOOS=linux GOARCH=arm go build hue-alarm

# This will upload the ARM binary to my beaglebone. Change this!
ssh root@192.168.1.175 -C mkdir /hue-alarm
scp hue-alarm root@192.168.1.175:/hue-alarm/hue-alarm