// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import axios from 'axios'
import 'babel-polyfill'

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
import 'ant-design-vue/dist/antd.css'

Vue.config.productionTip = false
Vue.prototype.$ajax = axios

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

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: {
    App
  },
  template: '<App/>'
})
