<template>
  <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '280px' }">
    <router-view></router-view>
  </a-layout-content>
</template>

<script>
export default {
  // data() {
  //     return {
  //         collapsed: false,
  //     };
  // },
  name: 'Content',
  props: {
    content: String
  }
}
</script>

<style>
</style>
