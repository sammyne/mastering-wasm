#!/bin/bash

set -e

workdir=$PWD
wavm_dir=$workdir/../../wavm

cd $workdir
wat2wasm param.wat

cd $wavm_dir/cmd/wavm
go run main.go $workdir/param.wasm

if [[ $? -ne 0 ]]; then
  echo "fail to run param.wasm"
  exit 1
fi

echo "fine :)"
