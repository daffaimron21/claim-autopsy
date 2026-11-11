import { createApp } from 'vue'
import App from './App.vue'
import { makeI18n, savedLocale } from './i18n'
import './style.css'

createApp(App).use(makeI18n(savedLocale())).mount('#app')
