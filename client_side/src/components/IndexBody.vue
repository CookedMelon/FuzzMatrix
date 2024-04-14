<template>
  <div class="common-layout" style="width:100%;flex-grow: 1;">
    <el-container style="height: 100%;">
      <el-aside width="200px">

        <el-menu
            default-active="2"
            class="el-menu-vertical-demo"
            @open="handleOpen"
            @close="handleClose"
        >
            <el-menu-item index="1" v-on:click="$router.push('/index/Overview')">
              <el-icon><House /></el-icon>
              <template #title>Overview</template>
            </el-menu-item>
            <el-menu-item index="2" v-on:click="$router.push('/index/CFG')">
            <el-icon><View /></el-icon>
            <template #title>CFG</template>
            </el-menu-item>
            <el-menu-item index="3" v-on:click="$router.push('/index/Seeds')">
            <el-icon><setting /></el-icon>
            <template #title>Seeds</template>
            </el-menu-item> 
            <el-menu-item index="4" v-on:click="$router.push('/index/Terminal-main')">
              <el-icon><postcard /></el-icon>
              <template #title>Terminal</template>
            </el-menu-item>

            <!-- <el-submenu index="5" >
              <template>
                <slot name="title">
                <i class="el-icon-postcard"></i>
                <span>Terminal</span>
              </template> 
              <el-menu-item-group>
                <el-menu-item index="2-1">
                  <el-icon><plus /></el-icon>
                  <template #title>New Terminal</template>
                </el-menu-item>
              </el-menu-item-group>
              
            </el-submenu> -->

        </el-menu>
      </el-aside>

      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script>
  export default {
    methods: {
      handleOpen(key, keyPath) {
        console.log(key, keyPath);
      },
      handleClose(key, keyPath) {
        console.log(key, keyPath);
      }
    }
  }
</script>
<style scoped>
.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 200px;
  min-height: 400px;
}


</style>