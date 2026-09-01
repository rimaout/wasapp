import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ui/ErrorMsg.vue'
import LoadingSpinner from './components/ui/LoadingSpinner.vue'

import 'bootstrap-icons/font/bootstrap-icons.css'
import './assets/dashboard.css'
import './assets/main.css'
import './assets/theme.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.use(router)
app.mount('#app')
