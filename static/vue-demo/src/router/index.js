import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import SkuSpecList from '@/components/SkuSpecList'
import SkuRequestList from '@/components/SkuRequestList'
import SkuResponseList from '@/components/SkuResponseList'
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
      name: 'Home',
      component: Home
    },
    {
      path: '/sku-spec-list',
      name: 'SkuSpecList',
      component: SkuSpecList
    },
    {
      path: '/sku-request-list',
      name: 'SkuRequestList',
      component: SkuRequestList
    },
    {
      path: '/sku-response-list',
      name: 'SkuResponseList',
      component: SkuResponseList
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
