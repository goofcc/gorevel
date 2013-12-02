#!/usr/bin/env bash

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

# go install ?
bin/revel build gorevel ./exe

export GOPATH="$OLDGOPATH"

