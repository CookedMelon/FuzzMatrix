import { createApp , ref} from 'vue'
import App from './App.vue'
import router from './router'
import 'element-plus/dist/index.css'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const app = createApp(App)

// 创建全局变量
const detail_val=ref({
    "run_time":0,
    "cycles_done":0,
    "total_path":0,
    "unique_crashes":0,
    "unique_hangs":0,
    "map_density":0,
    "max_depth":0,
})
app.provide("detail_val", detail_val)


for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(ElementPlus)
app.use(router)

app.mount('#app')
