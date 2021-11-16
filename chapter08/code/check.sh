#!/bin/bash

set -e

workdir=$PWD
wavm_dir=$workdir/../../wavm

wats=$(ls *.wat)

cd $workdir
for v in ${wats[@]}; do
  echo "checking $v ..."

  out=${v%.*}.wasm
  wat2wasm $v

  cd $wavm_dir/cmd/wavm
  go run main.go $workdir/$out

  if [[ $? -ne 0 ]]; then
    echo "fail to run $v"
    exit 1
  fi
done

echo "fine :)"
