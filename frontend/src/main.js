import { createApp } from 'vue'
import App from './App.vue'
import { loadConfig } from "./lib/config/config.js";
import './index.css'

loadConfig();
createApp(App).mount('#app');
