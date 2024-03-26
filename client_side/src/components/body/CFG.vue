<template>
    <div style="height:100%;display:flex;flex-direction:row">
        <d3-chart></d3-chart>
        <cfg-right></cfg-right>
    </div>
</template>
<script>
import '@element-plus/icons-vue'
import D3Chart from './D3CFG.vue';
import CFG_right from './CFG-right.vue';
export default {
    name: 'CFG',
    components: {
        'd3-chart': D3Chart,
        'cfg-right': CFG_right
    },
    mounted(){
    },
    methods:{
        
    }
}
</script>
<style scoped>



</style>