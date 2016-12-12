#!/bin/bash

CURDIR=$(cd `dirname $0`; pwd)
PRJHOME=$(dirname $CURDIR)

export GOPATH="$PRJHOME"
export GOBIN="$PRJHOME/bin"

go "$@"