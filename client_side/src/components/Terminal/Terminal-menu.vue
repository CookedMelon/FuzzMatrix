<template>
  <div class="terminal-menu">
      <el-container style="height: 100%;">
          <el-aside width="200px">
    
            <el-menu
                class="el-menu-vertical-demo"
                @select="handleSelect"
                :default-active="selectedTerminal"
            >
                <el-menu-item index="1" >
                  <el-icon><plus /></el-icon>
                  <template #title>New Terminal</template>
                </el-menu-item>

                <!-- 动态生成的终端菜单项 -->
                <el-menu-item 
                  v-for="(item, index) in terminalItems" :key="index" :index="'term-' + (index + 1)"
                  @click="handleTerminalClick(index)"
                > 
                  <el-icon><postcard /></el-icon>
                  {{ item }}                    
                </el-menu-item>

            </el-menu>
          </el-aside>
          
          <!-- 保存的历史记录 -->
          <el-main>
            <router-view></router-view>
            <div v-if="history.length > 0" class="history-ssh-info">
              <div class="history-title">History Hosts</div> 
              <div class="history-card-container">
                <el-card
                  v-for="(item, index) in history"
                  :key="index"
                  class="history-card"
                  shadow="hover"
                  style="width: 300px; margin-bottom: 20px;" 
                >
                  <!-- 主体悬浮框 -->
                  <div v-bind:key="index" class="clearfix">
                    <span>Host: {{ item.HOST }}</span>
                  </div>
                  <div>Port: {{ item.PORT }}</div>
                  <div>Username: {{ item.USERNAME }}</div>
                  <div class="history-footer">
                    <!-- 删除按钮 -->
                    <el-button type="danger" @click="removeHost(index)">Remove</el-button>
                    <!-- 连接按钮 -->
                    <el-button type="primary" @click="redirectToTerminal(index)">Connect</el-button>
                  </div>
                </el-card>
              </div>
            </div>

          </el-main>
      </el-container>
      
      <!-- SSH信息列表 -->
      <div class="saved-ssh-info">
        <h2>Saved SSH Information</h2>
        <ul>
          <li v-for="(item, index) in savedSSHInfo" :key="index">
            <div>Host: {{ item.HOST }}</div>
            <div>Port: {{ item.PORT }}</div>
            <div>Username: {{ item.USERNAME }}</div>
            <el-button @click="redirectToTerminal(index)">Connect</el-button>
          </li>
        </ul>
      </div>

      <ul>
          <li @click="toggleTerminal">New Terminal</li>
          
      </ul>

  </div>

  <!-- 点击 Terminal-1 填写SSH表单 -->
  <el-dialog v-model="dialogFormVisible" title="SSH" width="500">
    <el-form :model="form">
      <el-form-item label="HOST" :label-width="formLabelWidth" prop="host">
        <el-input v-model="form.HOST" autocomplete="off"  placeholder="please input your host" />
        <template #label>
          <span style="color: red">*</span> HOST
        </template>
      </el-form-item>

      <el-form-item label="PORT" :label-width="formLabelWidth" prop="port">
        <el-input v-model="form.PORT" autocomplete="off"  placeholder="please input your port" />
        <template #label>
          <span style="color: red">*</span> PORT
        </template>
      </el-form-item>

      <el-form-item label="USERNAME" :label-width="formLabelWidth" prop="username">
        <el-input v-model="form.USERNAME" autocomplete="off"  placeholder="please input your username" />
        <template #label>
          <span style="color: red">*</span> USERNAME
        </template>
      </el-form-item>

      <el-form-item label="PASSWORD" :label-width="formLabelWidth" prop="password">
        <el-input v-model="form.PASSWORD" autocomplete="off"  placeholder="please input your password" show-password/>
      </el-form-item>

      <el-form-item label="KEY" :label-width="formLabelWidth" prop="key">
        <el-input v-model="form.KEY" autocomplete="off"  placeholder="please input your key" />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          Cancel
        </el-button>
        <el-button type="primary" @click="connectHost()">
          Connect
        </el-button>
        <el-button type="primary" @click="saveHost()">
          Save
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script>
// import axios from 'axios';

export default {
  data() {
    return {
      terminalItems: [
        'Terminal-1'
      ],// 用于存储动态生成的终端菜单项
      dialogFormVisible: false,
      formLabelWidth: '140px',
      form: {
        HOST: '1.1.1.1',
        PORT: '22',
        USERNAME: '',
        PASSWORD: '',
        KEY: ''
      },
      selectedTerminal: '1',
      history: [], // 用于存储历史SSH信息
    };
  },
  created() {
    this.loadHistory();
  },
  methods: {
    loadHistory() {
      // 在这里加载历史记录的逻辑
      const savedHistory = localStorage.getItem('history');
      if (savedHistory) {
        this.history = JSON.parse(savedHistory);
      }
    },
    handleSelect(key) {
      if (key === '1') {
        this.selectedTerminal = '';
        // 当点击 "New Terminal" 时，追加一个新的终端菜单项
      const newTerminalIndex = this.terminalItems.length + 1;
      this.terminalItems.push(`Terminal-${newTerminalIndex}`);
      this.activeIndex = `term-${this.terminalItems.length - 1}`; // 更新默认激活的菜单项
      } else {
        console.log('Selected:', key);
        this.selectedTerminal = key;
      }
    },
    handleTerminalClick(index) {
      console.log('Terminal clicked:', index);
      this.dialogFormVisible = true;
      this.selectedTerminal = `term-${index + 1}`;
    },

    connectHost() {
      //表单点击连接
      if (this.form.HOST == '' || this.form.PORT == '' || this.form.USERNAME == '' || (this.form.PASSWORD == '' && this.form.KEY == '')) {
        this.$message({
          message: 'please fill in all required fields',
          type: 'error',
        })
      } else {
        this.dialogFormVisible = false;
        this.$router.push('/index/Terminal-1');
        // Sending form data to the server
        // axios.post('/api/add_host', this.form)
        //   .then(response => {
        //     console.log('Server response:', response.data);
            
        //   })
        //   .catch(error => {
        //     console.error('Error sending data to server:', error);
        //   });
      }
      
    },
    saveHost() {
      //点击保存
      if (this.form.HOST == '' || this.form.PORT == '' || this.form.USERNAME == '' || (this.form.PASSWORD == '' && this.form.KEY == '')) {
        this.$message({
          message: 'please fill in all required fields',
          type: 'error',
        })
      } else {
        localStorage.setItem('savedFormData', JSON.stringify(this.form));
        this.history.push({ ...this.form }); // 将表单数据添加到历史记录中
        localStorage.setItem('history', JSON.stringify(this.history));
        this.$message({
          message: 'Form saved successfully',
          type: 'success',
        })
      }
    },
    redirectToTerminal(index) {
      localStorage.setItem('history', JSON.stringify(this.history));
      this.$router.push('/index/Terminal-' + (index + 1));
    },
    removeHost(index) {
      // 删除主机信息
      this.history.splice(index, 1);
      localStorage.setItem('history', JSON.stringify(this.history));
    }
  
  }
};
</script>

<style scoped>
.history-title {
  text-align: left; /* 设置标题靠左对齐 */
  margin-bottom: 50px; /* 设置标题下方间距 */
}

.history-card-container {
  display: flex; /* 设置为 Flex 容器 */
  flex-wrap: wrap; /* 控制子元素换行 */
}

.delete-button {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 20px;
  height: 20px;
  line-height: 20px;
  text-align: center;
  cursor: pointer;
  background-color: #f56c6c;
  color: #fff;
  border-radius: 50%;
}

</style>
