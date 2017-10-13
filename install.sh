#!/usr/bin/env

if [[ ! -f ./paste_table.go ]]; then
  echo "ERROR: Please install in the paste-table-go folder" 1>&2
  exit 1
fi

ORIGINAL_GOPATH=$GOPATH
GOPATH=$(pwd)

go get github.com/mkideal/cli
go build

GOPATH=$ORIGINAL_GOPATH
