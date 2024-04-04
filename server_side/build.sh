#!/bin/bash
export DEMO_PATH=$PWD/demo
export TOOL_PATH=$PWD/visfuzz
cd $TOOL_PATH/fuzz
mkdir build
cd build
cmake ../llvm/ .
make
export VISFUZZ_BUILD=$PWD
cd $TOOL_PATH/fuzz/afl
make
cd $DEMO_PATH/re2
bash compile.sh
cd repo
mkdir in
echo a>in/seed