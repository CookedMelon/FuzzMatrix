<template>
    <el-card id="right_side" style="
        background-color:darkseagreen;
        flex:1;
        height:100%;
        box-shadow: var(--el-box-shadow-light);
        border-radius: 8px;
        margin-left:10px;
        display:flex;
        flex-direction:row;
        font-size:larger;
        text-align:left;
        "
    >
        <p>
            <b>visiting frequency:</b>
            <br/>
            <br/>
            <img src="@/fig/frequency.png" style="width:100%">
        </p>
        <p>
            <b>contain basic block number:</b>
            <br/>
            <br/>
            <img src="@/fig/contain_bb.png" style="width:100%">
        </p>
        <p>
            <b>run time: </b>
            {{ detail_val.run_time }}
        </p>
        <p>
            <b>cycles done: </b>
            {{ detail_val.cycles_done }}
        </p>
        <p>
            <b>total path: </b>
            {{ detail_val.total_path }}
        </p>
        <p>
            <b>unique crashes: </b>
            {{ detail_val.unique_crashes }}
        </p>
        <p>
            <b>unique hangs: </b>
            {{ detail_val.unique_hangs }}
        </p>
        <p>
            <b>execs done: </b>
            {{ detail_val.execs_done }}
        </p>
        <p>
            <b>map density: </b>
            {{ detail_val.map_density }}
        </p>
        <p>
            <b>max depth: </b>
            {{ detail_val.max_depth }}
        </p>
    </el-card>
</template>
<script>
import '@element-plus/icons-vue'
import { inject } from 'vue'

export default {
    name: 'CFG-right',
    setup(){
        const detail_val = inject("detail_val");
        return{
            detail_val
        }
    },
    data() {
      return {
      };
    },
    components: {
    },
    mounted(){
    },
    methods:{
        
    }
}
</script>
<style scoped>



</style>