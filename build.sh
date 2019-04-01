#!/usr/bin/env bash

function main
{
        export GOPATH=$(pwd)
	    go build -a -o  zsmonitor  main.go
}

main "$@"
