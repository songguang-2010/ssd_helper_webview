// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import moment from "moment";
import App from './App'
import router from './router'
import axios from 'axios'
// import VueAxios from 'vue-axios'
import 'babel-polyfill'
import store from './store'

import Button from 'ant-design-vue/lib/button'
import Layout from 'ant-design-vue/lib/layout'
import Menu from 'ant-design-vue/lib/menu'
import Icon from 'ant-design-vue/lib/icon'
import Breadcrumb from 'ant-design-vue/lib/breadcrumb'
import Table from 'ant-design-vue/lib/table'
import Input from 'ant-design-vue/lib/input'
import Form from 'ant-design-vue/lib/form'
import Col from 'ant-design-vue/lib/col'
import Row from 'ant-design-vue/lib/row'
import DatePicker from 'ant-design-vue/lib/date-picker'
import Select from 'ant-design-vue/lib/select'
import Tag from 'ant-design-vue/lib/tag'
import 'ant-design-vue/dist/antd.css'

Vue.use(Button)
Vue.use(Layout)
Vue.use(Menu)
Vue.use(Icon)
Vue.use(Breadcrumb)
Vue.use(Table)
Vue.use(Input)
Vue.use(Form)
Vue.use(Col)
Vue.use(Row)
Vue.use(DatePicker)
Vue.use(Select)
Vue.use(Tag)

// axios 配置
axios.defaults.timeout = 5000
axios.defaults.baseURL = '/api'
if (process.env.NODE_ENV == "production") {
  axios.defaults.baseURL = '/'
}
console.log(process.env.NODE_ENV)

//初始化运行时, 清除token信息
store.commit("removeToken");

// 添加请求拦截器，在请求头中加token
axios.interceptors.request.use(
  config => {
    if (store.state.token) {
      config.headers.Authorization = store.state.token;
    }

    return config;
  },
  error => {
    return Promise.reject(error);
  });

// http response 拦截器
axios.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // 401 清除token信息
          store.commit("removeToken");

          // 只有在当前路由不是登录页面才跳转
          router.currentRoute.path !== 'login' && router.push('/login')
          // router.replace({
          //   path: '/login',
          //   query: {
          //     redirect: router.currentRoute.path
          //   },
          // })
      }
    }
    // console.log(JSON.stringify(error));//console : Error: Request failed with status code 402
    return Promise.reject(error.response.data)
  },
)

// Vue.use(axios, VueAxios)

Vue.config.productionTip = false
Vue.prototype.$ajax = axios
Vue.prototype.$moment = moment

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: {
    App
  },
  template: '<App/>'
})
