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

  <!-- 点击 Terminal-1 填写SSH表单 -->
  <el-dialog v-model="dialogFormVisible" title="SSH" width="500">
    <el-form :model="form">
      <el-form-item label="HOST" :label-width="formLabelWidth" prop="host">
        <el-input v-model="form.HOST" autocomplete="off"  placeholder="please input your host" />
      </el-form-item>
      <el-form-item label="PORT" :label-width="formLabelWidth" prop="port">
        <el-input v-model="form.PORT" autocomplete="off"  placeholder="please input your port" />
      </el-form-item>
      <el-form-item label="USERNAME" :label-width="formLabelWidth" prop="username">
        <el-input v-model="form.USERNAME" autocomplete="off"  placeholder="please input your username" />
      </el-form-item>
      <el-form-item label="PASSWORD" :label-width="formLabelWidth" prop="password">
        <el-input v-model="form.PASSWORD" autocomplete="off"  placeholder="please input your password" />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogFormVisible = false">Cancel</el-button>
        <el-button type="primary" @click="addHost">
          Confirm
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
        USERNAME: '1',
        PASSWORD: '1'
      },
      selectedTerminal: '1'
    };
  },
  methods: {
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
    addHost() {
      this.dialogFormVisible = false;
      // Sending form data to the server
      // axios.post('/api/add_host', this.form)
      //   .then(response => {
      //     console.log('Server response:', response.data);
          
      //   })
      //   .catch(error => {
      //     console.error('Error sending data to server:', error);
      //   });
      this.$router.push('/index/Terminal-1');
    }
    
  }
};
</script>

<style scoped>
</style>
