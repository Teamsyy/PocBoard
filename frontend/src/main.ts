import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './style.css'
import Vue3DraggableResizable from 'vue3-draggable-resizable'
import 'vue3-draggable-resizable/dist/Vue3DraggableResizable.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.component('Vue3DraggableResizable', Vue3DraggableResizable)

app.mount('#app')