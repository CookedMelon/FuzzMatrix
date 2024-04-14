<template>
  <div ref="terminal" class="terminal"></div>
</template>

<script>
/* eslint-disable vue/multi-word-component-names */
import { Terminal } from "xterm"
import { FitAddon } from 'xterm-addon-fit'
import "xterm/css/xterm.css"
export default {
  data() {
    return {
      term: "", // 保存terminal实例
      rows: 40,
      cols: 100
    }
  },
  mounted() {
    this.initXterm()
  },
  methods: {
    initXterm() {
      let _this = this
      let term = new Terminal({
        rendererType: "canvas", //渲染类型
        rows: _this.rows, //行数
        cols: _this.cols, // 不指定行数，自动回车后光标从下一行开始
        convertEol: true, //光标将设置为下一行的开头
        disableStdin: false, //输入
        cursorBlink: true, //光标闪烁
        selectToClipboard: true,
        theme: {
          foreground: "#ECECEC", //字体
          background: "#000000", //背景色
          cursor: "help", //设置光标
          lineHeight: 20
        }
      })
      // 创建terminal实例
      term.open(this.$refs["terminal"])
      // 换行并输入起始符 $
      term.prompt = () => {
        term.write("\r\n\x1b[33m$\x1b[0m ")
      }
      term.prompt()
      // canvas背景全屏
      const fitAddon = new FitAddon()
      term.loadAddon(fitAddon)
      fitAddon.fit()

      window.addEventListener("resize", resizeScreen)
      function resizeScreen() {
        try { // 窗口大小改变时，触发xterm的resize方法使自适应
          fitAddon.fit()
        } catch (e) {
          console.log("e", e.message)
        }
      }
      _this.term = term
      _this.runFakeTerminal()
    },
    runFakeTerminal() {
      let _this = this
      let term = _this.term
      if (term._initialized) return
      // 初始化
      term._initialized = true
      term.writeln('Welcome to FuzzMatrix.')
      term.writeln('A user-friendly fuzzing vulnerability mining system.')
      term.prompt()
      // 添加事件监听器，支持输入方法
      term.onKey(e => {
        const printable = !e.domEvent.altKey && !e.domEvent.altGraphKey && !e.domEvent.ctrlKey && !e.domEvent.metaKey
        if (e.domEvent.keyCode === 13) {
          term.prompt()
        } else if (e.domEvent.keyCode === 8) { // back 删除的情况
          if (term._core.buffer.x > 2) {
            term.write('\b \b')
          }
        } else if (printable) {
          term.write(e.key)
        }
        console.log(1,'print', e.key)
      })
      term.onData(key => {  // 粘贴的情况
        if(key.length > 1) term.write(key)
      })
    }
    
  }
}
</script>

<style scoped>
.terminal {
  width: 100%;
  height: 100%;
}
</style>