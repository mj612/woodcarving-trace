import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import './assets/styles/global.css'
import * as api from './api'

Vue.config.productionTip = false

// 如果已有token，尝试获取用户信息
if (store.getters.isLoggedIn) {
  store.dispatch('initUserInfo')
}

Vue.use(ElementUI)

// 将API挂载到Vue原型上
Vue.prototype.$api = api

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
