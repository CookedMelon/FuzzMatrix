#!/bin/bash
export DEMO_PATH=$PWD/demo
export TOOL_PATH=$PWD/visfuzz
export VISFUZZ_BUILD=$PWD
cd $DEMO_PATH/re2
cd repo
nohup $TOOL_PATH/fuzz/afl/afl-fuzz -i in -o out ./app @@ &
pip install -r $TOOL_PATH/requirements.txt
python $TOOL_PATH/open_file_server.py 6767