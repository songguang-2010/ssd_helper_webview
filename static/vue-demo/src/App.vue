<template>
  <a-layout id="components-layout-demo-top-side-2">
    <Header msg="Welcome to Your SSD Helper App" />
    <router-view></router-view>
  </a-layout>
</template>
<script>
import Header from "@/components/common/Header.vue";
// import Sidebar from "@/components/common/Sidebar.vue";
// import Breadcrumb from "@/components/common/Breadcrumb.vue";
// import Content from "@/components/common/Content.vue";

export default {
  name: "app",
  components: {
    Header
    // Sidebar,
    // Breadcrumb,
    // Content
  },
  // methods: {
  //   onMenuClick: function (name) {
  //     this.content = name
  //   }
  // },
  data() {
    return {
      // collapsed: false,
      // content: "<div style='color: #ff0000;'>Dashboard</div>"
    };
  }
};
</script>

<style>
</style>
