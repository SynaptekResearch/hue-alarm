#!/bin/bash

# Auto-detect root to set GOPATH
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo = Initialize root = $ROOT
export GOPATH=$ROOT
cd src/hue-alarm
glide i

echo = Compile
cd ../..

# This default target
go build hue-alarm \
  && mv hue-alarm hue-alarm.osx.x86 &

# Beaglebone and other ARM linux 
env GOOS=linux GOARCH=arm go build hue-alarm \
  && mv hue-alarm hue-alarm.linux.arm \
  && ssh root@192.168.1.175 -C mkdir /hue-alarm \
  || scp hue-alarm.linux.arm root@192.168.1.175:/hue-alarm/hue-alarm &

# AMD64 linux systems
env GOOS=linux GOARCH=amd64 go build hue-alarm \
  && mv hue-alarm hue-alarm.linux.x86 &

wait