import Vue from 'vue'
import Vuex from 'vuex'
import Cookies from 'js-cookie'
import { login, getUserInfo } from '@/api'

Vue.use(Vuex)

const TokenKey = 'woodcarving_token'

export default new Vuex.Store({
  state: {
    token: Cookies.get(TokenKey) || '',
    userInfo: null,
    role: ''
  },
  
  getters: {
    token: state => state.token,
    userInfo: state => state.userInfo,
    role: state => state.role,
    isLoggedIn: state => !!state.token
  },
  
  mutations: {
    SET_TOKEN(state, token) {
      state.token = token
      Cookies.set(TokenKey, token)
    },
    
    SET_USER_INFO(state, userInfo) {
      state.userInfo = userInfo
      state.role = userInfo.role
    },
    
    CLEAR_USER(state) {
      state.token = ''
      state.userInfo = null
      state.role = ''
      Cookies.remove(TokenKey)
    }
  },
  
  actions: {
    // 登录
    async login({ commit }, loginForm) {
      try {
        const res = await login(loginForm)
        const { token, userInfo } = res.data
        
        commit('SET_TOKEN', token)
        commit('SET_USER_INFO', userInfo)
        
        return Promise.resolve(res)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    // 获取用户信息
    async getUserInfo({ commit }) {
      try {
        const res = await getUserInfo()
        commit('SET_USER_INFO', res.data)
        return Promise.resolve(res)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    // 登出
    logout({ commit }) {
      commit('CLEAR_USER')
      return Promise.resolve()
    },
    
    // 初始化用户信息
    initUserInfo({ commit, dispatch }) {
      return dispatch('getUserInfo').catch(() => {
        // 如果获取用户信息失败，清除登录状态
        commit('CLEAR_USER')
      })
    }
  }
})
