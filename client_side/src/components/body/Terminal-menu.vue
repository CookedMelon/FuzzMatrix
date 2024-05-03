<template>
    <div class="terminal-menu">
        <el-container style="height: 100%;">
            <el-aside width="200px">
      
              <el-menu
                  default-active="1"
                  class="el-menu-vertical-demo"
                  @select="handleSelect"
              >
                  <el-menu-item index="1" >
                    <el-icon><plus /></el-icon>
                    <template #title>New Terminal</template>
                  </el-menu-item>
                  
                  <!-- <el-menu-item index="2" v-on:click="$router.push('/index/Terminal-1')">
                    <el-icon><postcard /></el-icon>
                    <template #title>Terminal-1</template>
                  </el-menu-item> -->

                  <!-- 动态生成的终端菜单项 -->
                  <el-menu-item 
                    v-for="(item, index) in terminalItems" :key="index" :index="'term-' + (index + 1)"
                    v-on:click="$router.push('/index/Terminal-'+ (index + 1))"
                    v-to="'/index/Terminal-'+ (index + 1)"
                    @click="dialogFormVisible = true"   
                    > 
                    <el-icon><postcard /></el-icon>
                    {{ item }}                    
                  </el-menu-item>

              </el-menu>
            </el-aside>
      
            <el-main>
              <router-view></router-view>
            </el-main>
        </el-container>
        <!-- <div>
            <NewTerminal></NewTerminal>
        </div> -->
        <ul>
            <li @click="toggleTerminal">New Terminal</li>
            
        </ul>
    </div>
  </template>
  
<script>
// import NewTerminal from './Terminal-1.vue'
export default {
  data() {
    return {
      activeIndex: '1',
      terminalItems: [
        'Terminal-1'
      ], // 用于存储动态生成的终端菜单项
    };
  },
  methods: {
    handleSelect(key) {
      if (key === '1') {
        // 当点击 "New Terminal" 时，追加一个新的终端菜单项
        const newTerminalIndex = this.terminalItems.length + 1;
        this.terminalItems.push(`Terminal-${newTerminalIndex}`);
        this.activeIndex = `term-${this.terminalItems.length - 1}`; // 更新默认激活的菜单项
      } else {
        // 处理其他菜单项的点击事件
        console.log('Selected:', key);
      }
    }
  }
};
</script>
  
<style scoped>
</style>