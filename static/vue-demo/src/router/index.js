import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import SpecList from '@/components/SpecList'
import SsdOrderList from '@/components/SsdOrderList'
import AosOrderList from '@/components/AosOrderList'
import TpsOrderList from '@/components/TpsOrderList'
import MiscDeviceList from '@/components/MiscDeviceList'

Vue.use(Router)

const Home = {
  template: '<div>home</div>'
}

export default new Router({
  routes: [{
      path: '/',
      name: 'SSDHelper',
      component: Home
    },
    {
      path: '/spec-list',
      name: 'SpecList',
      component: SpecList
    },
    {
      path: '/ssd-order-list',
      name: 'SsdOrderList',
      component: SsdOrderList
    },
    {
      path: '/aos-order-list',
      name: 'AosOrderList',
      component: AosOrderList
    },
    {
      path: '/tps-order-list',
      name: 'TpsOrderList',
      component: TpsOrderList
    },
    {
      path: '/misc-device-list',
      name: 'MiscDeviceList',
      component: MiscDeviceList
    }
  ]
})
