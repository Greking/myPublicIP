#!/bin/bash 


DIR="$( cd "$( dirname "$0"  )" && pwd  )"

OLDGOPATH=$GOPATH

export GOPATH=$DIR:$OLDGOPATH

go install publicIP


export GOPATH=$OLDGOPATH
