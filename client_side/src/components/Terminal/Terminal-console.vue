<template>
  <div class="console" id="terminal"></div>
</template>
<script>
import Terminal from './Xterm'
export default {
name: 'WebConsole',
props: {
  terminal: {
    type: Object,
    default: () => ({})
  }
},
data () {
  return {
    term: null,
    terminalSocket: null
  }
},
methods: {
  runRealTerminal () {
    console.log('webSocket is finished')
  },
  errorRealTerminal () {
    console.log('error')
  },
  closeRealTerminal () {
    console.log('close')
  }
},
mounted () {
  console.log('pid : ' + this.terminal.pid + ' is on ready')
  let terminalContainer = document.getElementById('terminal')
  this.term = new Terminal()
  this.term.open(terminalContainer)
  // open websocket
  this.terminalSocket = new WebSocket('ws://127.0.0.1:3000/terminals')
  this.terminalSocket.onopen = this.runRealTerminal
  this.terminalSocket.onclose = this.closeRealTerminal
  this.terminalSocket.onerror = this.errorRealTerminal
  this.term.attach(this.terminalSocket)
  this.term._initialized = true
  console.log('mounted is going on')
},
beforeUnmount () {
  this.terminalSocket.close()
  this.term.destroy()
}
}
</script>

<style scoped>
.console#terminal {
  width: 100%;
  height: 100%;
}
</style>

