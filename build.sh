#!/bin/bash

# Auto-detect root to set GOPATH
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../../../.."

echo = Initialize root = $ROOT
export GOPATH=$ROOT
glide i

echo = Compile
webpack --display-modules

# This default target
go build github.com/cpo/hue-alarm \
  && mv hue-alarm hue-alarm.osx.x86 &

if [ "$1" == "dist" ] ; then
  # generic static artifacts
  rm artifacts.tar.gz
  tar czvf artifacts.tar.gz static/index.html static/bundle.js static/logo.png node_modules/semantic-ui/dist/semantic.min.css 

  # Beaglebone and other ARM linux 
  env GOOS=linux GOARCH=arm go build github.com/cpo/hue-alarm \
    && mv hue-alarm hue-alarm.linux.arm \
    && scp artifacts.tar.gz root@192.168.1.175:/hue-alarm/artifacts.tar.gz \
    && ssh root@192.168.1.175 -C 'cd /hue-alarm && tar zxvf ./artifacts.tar.gz' \
    && scp hue-alarm.linux.arm root@192.168.1.175:/hue-alarm/hue-alarm  

  # AMD64 linux systems
  env GOOS=linux GOARCH=amd64 go build github.com/cpo/hue-alarm \
    && mv hue-alarm hue-alarm.linux.x86 &
fi

wait
