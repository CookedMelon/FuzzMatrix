<template>
  <div class="container">
    <terminal-main></terminal-main>
    <my-terminal></my-terminal>
  </div>
</template>

<script>
import Console from './Console'
import TerminalMain from './Terminal-main.vue'

export default {
  name: 'WebSSH',
  data () {
    return {
      terminal: {
        pid: 1,
        name: 'terminal',
        cols: 400,
        rows: 400
      }
    }
  },
  components: {
    'my-terminal': Console,
    'terminal-main': TerminalMain
  }
}
</script>