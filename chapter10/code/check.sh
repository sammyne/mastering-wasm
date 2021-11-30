#!/bin/bash

set -e

workdir=$PWD
src=$workdir/../../chapter01/code/hello-world/webapp/hello-world.wasm
wavm_dir=$workdir/../../wavm/cmd/wavm

#wasm-objdump -x $src

cd $wavm_dir
go run main.go $src
