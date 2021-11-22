#!/bin/bash

set -e

workdir=$PWD
wavm_dir=$workdir/../../wavm

wat=calc2.wat

cd $workdir

out=${wat%.*}.wasm
wat2wasm $wat

cd $wavm_dir/cmd/wavm
go run main.go $workdir/$out

echo "fine :)"
