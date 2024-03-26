<template>
  <div class="common-layout" style="width:100%;flex-grow: 1;">
    <el-container style="height: 100%;">
      <el-aside width="200px">
        <el-menu
            default-active="1"
            class="el-menu-vertical-demo"
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
    name: 'MyBody',
    components: {
    }
}
</script>
<style scoped>
.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 200px;
  min-height: 400px;
}


</style>