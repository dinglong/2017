#!/bin/bash

if [ ${#} != 2 ]; then
	echo "usage: ${0} <install | clean>"
	exit 1
fi

if [ ${1} == "install" ]; then
	GOPATH=${PWD} go install microiot.com/center
elif [ ${1} == "clean" ]; then
	GOPATH=${PWD} go clean -i -n -x microiot.com/center
fi
