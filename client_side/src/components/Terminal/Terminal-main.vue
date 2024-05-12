<template>
  <div style="height:100%;display:flex;flex-direction:row">
      <terminal-menu></terminal-menu>
  </div>
</template>
<script>
  import TerminalMenu from './Terminal-menu.vue';
  export default {
    name: 'TERMINAL',
    components: {
        'terminal-menu': TerminalMenu
    },
    mounted(){
    },
    methods:{
    }
  }
</script>
<style scoped>



</style>