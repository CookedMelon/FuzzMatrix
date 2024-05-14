const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    hot: true, // 启用热重载
    proxy: {
      // 代理配置
      '/api': {
        target: 'http://localhost:6767', // 后端服务地址和端口
        changeOrigin: true, // 是否改变源地址
        pathRewrite: { '^/api': '' }, // 重写路径：去掉路径中开头的'/api'
        logLevel: 'debug' // 增加日志输出
      },
    },
  },
  configureWebpack:{
    devtool: 'source-map'
  }
})
