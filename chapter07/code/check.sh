#!/bin/bash

set -e

workdir=$PWD
wavm_dir=$workdir/../../wavm
app=fib

cd $workdir
wat2wasm $app.wat

cd $wavm_dir/cmd/wavm
go run main.go $workdir/$app.wasm

if [[ $? -ne 0 ]]; then
  echo "fail to run $app.wasm"
  exit 1
fi

echo "fine :)"
