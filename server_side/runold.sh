#!/bin/bash
export DEMO_PATH=$PWD/demo
export TOOL_PATH=$PWD/visfuzz
export VISFUZZ_BUILD=$PWD
cd $TOOL_PATH/visualization
python -m http.server 8000