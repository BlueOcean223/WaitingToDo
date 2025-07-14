import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue' // 引入ElementPlus图标库
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'


const app = createApp(App)
const pinia = createPinia()

// 注册ElementPlus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(router).use(ElementPlus).use(pinia).mount('#app')
