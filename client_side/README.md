# client_side

## How to use

The first half of the launch is the same as the previous work, but do not run the `python -m http.server 8000` command

1. Install LLVM (>=8), python3.

```bash
wget https://apt.llvm.org/llvm.sh
chmod +x llvm.sh
sudo ./llvm.sh 10
```

2. Download and compile VisFuzz:

```bash
git clone https://github.com/ChijinZ/VisFuzz.git
cd VisFuzz
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
```

3. Fuzz && visualize demo:
 
```bash
cd $DEMO_PATH/re2
sh compile.sh
cd repo
mkdir in
echo a>in/seed
nohup $TOOL_PATH/fuzz/afl/afl-fuzz -i in -o out ./app @@ &
python $TOOL_PATH/open_file_server.py 6767
```

## Project setup

```bash
npm install
```

### Compiles and hot-reloads for development

```bash
npm run serve
```

### Compiles and minifies for production

```bash
npm run build
```

### Use

open the browser with the link http://localhost:8080/index/