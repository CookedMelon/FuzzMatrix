<template>
    <div class="container">
      <Header></Header>
  
      <Body></Body>
    </div>
  </template>
  
<script>
  import Header from './IndexHead.vue'
  import Body from './IndexBody.vue'
  export default {
    name: 'IndexPage',
    components: {
      Header,Body
    }
  }
</script>
<style scoped>
.container {
  display: flex;
  flex-direction: column;
  height: 100vh; /* 使用视窗的完整高度 */
}
</style>