import Vue from 'vue'
import Router from 'vue-router'
import store from '../store'

import Layouts from '@/components/common/Layouts'
// import HelloWorld from '@/components/HelloWorld'
import SkuSpecList from '@/components/SkuSpecList'
import SkuRequestList from '@/components/SkuRequestList'
import SkuResponseList from '@/components/SkuResponseList'
import SsdOrderList from '@/components/SsdOrderList'
import AosOrderList from '@/components/AosOrderList'
import TpsOrderList from '@/components/TpsOrderList'
import MiscDeviceList from '@/components/MiscDeviceList'
import Login from '@/components/Login';

Vue.use(Router)

const Home = {
  template: '<div>home</div>'
}

const router = new Router({
  routes: [
    // {
    //   path: '/',
    //   redirect: '/login'
    // },
    {
      path: '/login',
      name: 'Login',
      component: Login,
      // 添加该字段requireAuth: true,表示进入需要认证
      meta: {
        requireAuth: false
      }
    },
    {
      path: '/',
      component: Layouts,
      children: [{
          path: '/home',
          name: 'Home',
          component: Home,
        },
        {
          path: '/sku-spec-list',
          name: 'SkuSpecList',
          component: SkuSpecList,
          // 添加该字段requireAuth: true,表示进入需要认证
          meta: {
            requireAuth: true
          }
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
    }
  ]
});

//路由过滤器
router.beforeEach((to, from, next) => {
  console.log("check each");
  // 判断该路由是否需要登录权限
  if (to.meta.requireAuth == false) {
    console.log("skip token");
    next();
  } else {
    console.log("check token");
    // 通过vuex state获取当前的token是否存在
    if (store.state.token) {
      console.log(store.state.token)
      console.log("next");
      next();
    } else {
      console.log("login");
      next('/login')
    }
  }
})

export default router;
