import './assets/main.css'
import { provideFluentDesignSystem, fluentCard, fluentButton } from '@fluentui/web-components'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)
provideFluentDesignSystem().register(fluentCard(), fluentButton())

app.use(createPinia())
app.use(router)

app.mount('#app')
