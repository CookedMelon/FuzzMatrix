<template>
    <div style="height:100%;display:flex;flex-direction:row">
        <terminal-menu></terminal-menu>
        <!-- 这里menu好像子菜单第一次建立的ter-1会刷新掉 -->
        <terminal-editor></terminal-editor>
    </div>
</template>
<script>
import TerminalMenu from './Terminal-menu.vue';
import TerminalEditor from './Terminal-editor.vue';
/* eslint-disable vue/multi-word-component-names */

export default {
  name: 'Terminal-1',
    components: {
        'terminal-menu': TerminalMenu,
        'terminal-editor': TerminalEditor
    },

    data(){
      return{
        showTerminal: false
      };
    },
    methods: {
      handleNewTerminal() {
        this.showTerminal = true;
      }
    }
};
</script>
<style scoped>



</style>