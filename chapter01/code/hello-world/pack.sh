#!/bin/bash

outDir=$PWD/webapp
targetDir=$PWD/target/wasm32-unknown-unknown/release

cargo build --release

cp $targetDir/*.wasm $outDir
