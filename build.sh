#!/usr/bin/env bash
ORIGIN_PATH=$(pwd)
cd $(dirname $0)
SCRIPT_PATH=$(pwd)

BINARY_PATH=$SCRIPT_PATH/file-manager
go build -i -o $BINARY_PATH cmd/filemanager/*
if [[ "$?" != "0" ]]
then
    cd $ORIGIN_PATH
    exit 1
fi
cd $ORIGIN_PATH
echo "build finished ($BINARY_PATH)"
