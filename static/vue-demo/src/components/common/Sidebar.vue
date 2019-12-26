<template>
  <a-layout-sider width="200" style="background: #fff">
    <a-menu
      mode="inline"
      :defaultSelectedKeys="['1']"
      :defaultOpenKeys="['sub1']"
      :style="{ height: '100%', borderRight: 0 }"
    >
      <a-menu-item key="1">
        <router-link to="/">
          <a-icon type="home" />Dashboard
        </router-link>
      </a-menu-item>
      <a-sub-menu key="sub1">
        <span slot="title">
          <a-icon type="shop" />商品规格管理
        </span>
        <a-menu-item key="2">
          <router-link to="/spec-list">规格列表</router-link>
        </a-menu-item>
      </a-sub-menu>
      <a-sub-menu key="sub2">
        <span slot="title">
          <a-icon type="laptop" />终端订单
        </span>
        <a-menu-item key="3">
          <router-link to="/ssd-order-list">订单记录</router-link>
        </a-menu-item>
      </a-sub-menu>
      <a-sub-menu key="sub3">
        <span slot="title">
          <a-icon type="laptop" />聚合订单
        </span>
        <a-menu-item key="4">
          <router-link to="/aos-order-list">订单记录</router-link>
        </a-menu-item>
      </a-sub-menu>
      <a-sub-menu key="sub4">
        <span slot="title">
          <a-icon type="laptop" />支付管理
        </span>
        <a-menu-item key="5">
          <router-link to="/tps-order-list">支付记录</router-link>
        </a-menu-item>
      </a-sub-menu>
      <a-sub-menu key="sub5">
        <span slot="title">
          <a-icon type="laptop" />设备管理
        </span>
        <a-menu-item key="6">
          <router-link to="/misc-device-list">设备列表</router-link>
        </a-menu-item>
      </a-sub-menu>
    </a-menu>
  </a-layout-sider>
</template>

<script>
export default {
  // data() {
  //     return {
  //         collapsed: false,
  //     };
  // },
  name: "Sidebar",
  props: {
    // msg: String
  }
};
</script>

<style>
</style>
