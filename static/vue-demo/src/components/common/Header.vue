<template>
  <a-layout-header class="header">
    <div class="logo">
      <img src="../../assets/logo.png" />
      <h1>{{ msg }}</h1>
    </div>
    <a-menu
      theme="dark"
      mode="horizontal"
      :defaultSelectedKeys="['2']"
      :style="{ lineHeight: '64px', float: 'right' }"
    >
      <a-menu-item key="1">xx用户</a-menu-item>
      <!-- <a-menu-item key="2">nav 2</a-menu-item> -->
      <!-- <a-menu-item key="3">nav 3</a-menu-item> -->
    </a-menu>
  </a-layout-header>
</template>

<script>
export default {
  // data() {
  //     return {
  //         collapsed: false,
  //     };
  // },
  name: 'Header',
  props: {
    msg: String
  }
}
</script>

<style>
#components-layout-demo-top-side-2 .logo {
  width: auto;
  height: 31px;
  /* background: rgba(255, 255, 255, 0.2); */
  margin: 16px 28px 16px 0;
  float: left;
}
#components-layout-demo-top-side-2 .logo img {
  width: 120px;
  height: 31px;
  /* background: rgba(255, 255, 255, 0.2); */
  margin: 0;
  float: left;
}
#components-layout-demo-top-side-2 .logo h1 {
  width: auto;
  height: 31px;
  margin: 0 0 0 20px;
  padding: 0 0 0 0;
  float: left;
  font-style: none;
  font-size: 28px;
  line-height: 28px;
  color: white;
}
</style>
