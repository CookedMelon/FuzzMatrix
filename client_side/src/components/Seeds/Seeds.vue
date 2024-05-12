<template>
  <div class="container">
    <el-row type="flex" justify="center">
      <el-col :span="18">
        <el-table :data="tableData" style="width: 100%" height="500" :scrollbar-always-on="true">
          <el-table-column label="种子标志" prop="name"></el-table-column>
          <el-table-column label="内容" prop="content"></el-table-column>
          <el-table-column label="操作">
            <template v-slot:default="scope">
              <el-button size="mini" @click="handleEdit(scope.$index, scope.row)">修改</el-button>
              <el-button size="mini" type="danger" @click="handleDelete(scope.$index, scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
    </el-row>
    <el-row type="flex" justify="center">
      <el-col :span="18">
        <el-input type="textarea" :rows="5" placeholder="请输入内容" v-model="textarea"></el-input>
        <el-button type="primary" @click="onSubmit">提交</el-button>
      </el-col>
    </el-row>
  </div>
</template>

<script>

export default {
  name: 'SEEDS',
  data() {
    return {
      tableData: [],
      textarea: ''
    }
  },
  mounted() {
    this.fetchSeeds();
  },
  methods: {
    async fetchSeeds() {
      console.log(this.tableData)
      try {
        const response = await fetch('/api/seeds');
        const data = await response.json();
        this.tableData = data.seeds;
      } catch (error) {
        console.error('Error fetching seeds:', error);
      }
    },
    async handleEdit(index, row) {
    try {
        const newContent = await this.$prompt('请输入新的内容', '修改种子', {
            inputValue: row.content
        });
        await fetch(`/api/seeds/`, {  //修改请求URL
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                name: row.name,
                content: newContent
            })
        });
        this.fetchSeeds();
    } catch (error) {
        console.error('Error updating seed:', error);
    }
},
    async handleDelete(index, row) {
      try {
        await this.$confirm(`确定删除种子 "${row.name}" 吗?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        });
        await fetch(`/api/seeds/${row.name}`, { method: 'DELETE' });
        this.fetchSeeds();
      } catch (error) {
        console.error('Error deleting seed:', error);
      }
    },
    async onSubmit() {
    try {
        await fetch('/api/seeds/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                content: this.textarea  //只传递content
            })
        });
        this.textarea = '';
        this.fetchSeeds();  //重新获取种子列表
    } catch (error) {
        console.error('Error creating seed:', error);
    }
}
  }
}
</script>

<style scoped>
.container {
  margin-top: 20px;
}
.el-row {
  margin-bottom: 20px;
}
.el-button {
  margin-left: 10px;
}
</style>