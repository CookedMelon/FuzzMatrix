<template>
  <div class="container">
    <terminal-main></terminal-main>
    <my-terminal></my-terminal>
  </div>
</template>

<script>
import Console from './Console'
import TerminalMain from './Terminal-main.vue'

export default {
  name: 'Terminal-1',
  components: {
    'my-terminal': Console,
    'terminal-main': TerminalMain
  },
  mounted() {
    // 连接到 WebSocket 服务器
    const socket = new WebSocket("ws://localhost:3000/terminals");

    socket.onopen = () => {
      console.log("WebSocket 连接成功！");
      // 连接建立后可以发送消息或者进行其他操作
    };

    socket.onmessage = (event) => {
      console.log("收到消息：", event.data);
      // 处理从服务器接收到的消息
    };

    socket.onerror = (error) => {
      console.error("WebSocket 连接错误：", error);
      // 处理连接错误
    };

    socket.onclose = () => {
      console.log("WebSocket 连接已关闭");
      // 处理连接关闭
    };
  }
}
</script>
