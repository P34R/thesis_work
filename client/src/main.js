import { createApp } from 'vue'
import router from './router';
import App from './App.vue';
import store from './store/index.js';
import * as cryptojs from 'crypto-js';
import * as jscrypto from 'jscrypto';
import * as ellipticjs from 'elliptic';

const app=createApp(App);

app.config.globalProperties.$cryptojs =cryptojs;
app.config.globalProperties.$jscrypto =jscrypto;
app.config.globalProperties.$elliptic =ellipticjs;
app.use(store);
app.use(router);
app.mount('#app');
