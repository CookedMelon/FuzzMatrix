#!/bin/bash
export DEMO_PATH=$PWD/demo
export TOOL_PATH=$PWD/visfuzz
export GO_PATH=$PWD/go

cd $TOOL_PATH/fuzz
mkdir -p build
cd build
cmake ../llvm/ .
make
export VISFUZZ_BUILD=$PWD
cd $TOOL_PATH/fuzz/afl
make

cd $DEMO_PATH/re2
bash compile.sh
cd repo
mkdir -p in
echo a>in/seed

cd $GO_PATH
go mod tidy
go build -o server