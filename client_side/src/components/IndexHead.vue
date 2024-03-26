<template> 
<el-menu :default-active="1" class="top el-menu-demo " mode="horizontal"
    background-color="#000"
    text-color="#fff"
    active-text-color="#ffd04b"
    @select="handleSelect">
    <img src="../assets/logo.png" class="logo" style="height:50%"/>
    <h1 class="title" style="font-size: 24px;"><span style="color:#57caeb">FuzzVis</span> &nbsp; A user-friendly fuzzing vulnerability mining system</h1>
</el-menu>
    <!-- <div :style="{'background-image': 'url(' + require('../assets/logo.png') + ')', 'width': '100%', 'height': '75px', 'background-size': 'cover'}">  -->
<!-- <h1>{{ msg }}</h1>
    </div>  -->
</template>

<script>
import '@element-plus/icons-vue'
import '@/css/common.css'

export default {
    name: "MyHeader",
    props: {},
    mounted(){
    }
}

// const activeIndex = ref('1')
</script>
<style scoped>
.top {
    height:75px;
    background-color: #000000;
    display:flex;
    border-bottom:0;
}
.logo{
    height: 60px;
    display:inline-block;
    align-items: center;
    margin: auto 10px;
}
.title{
    color:#435ebe;
    margin: auto 0;
}
.right-container{
    display: flex;
    margin-left: auto;
    margin-right: 10px;
}
/* .right-item{
    color:#7b7b7b;
    margin: auto 30px;
    text-decoration: none;
    margin-left: auto;
    font-size: larger;
    text-decoration: underline
} */
.right-item{
    color:#7b7b7b;
    margin: auto 30px;
    text-decoration: none;
    margin-left: auto;
}
</style>