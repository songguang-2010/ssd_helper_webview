<template>
  <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '280px' }">
    <div>
      <input type="text" v-model="loginForm.username" placeholder="用户名" />
      <input type="text" v-model="loginForm.password" placeholder="密码" />
      <button @click="login">登录</button>
      <button @click="logout">注销</button>
    </div>
  </a-layout-content>
</template>
 
<script>
import { mapMutations } from "vuex";
export default {
  name: "MiscDeviceList",
  data() {
    return {
      error: "",
      userToken: "",
      loginForm: {
        username: "",
        password: ""
      }
    };
  },

  methods: {
    ...mapMutations(["setToken", "removeToken"]),
    logout() {
      let _this = this;
      _this.userToken = "";
      _this.removeToken({});
      // _this.$router.push("/login");
    },
    login() {
      if (this.loginForm.username === "" || this.loginForm.password === "") {
        alert("账号或密码不能为空");
        return false;
      }
      let _this = this;
      // this.searchValue.canary = this.form.getFieldsValue().canary;
      this.loginAction((err, data) => {
        console.log(data);
        if (err) {
          _this.error = err.toString();
          alert(_this.error);
        } else {
          _this.userToken = "Bearer " + data.token;
          console.log(_this.userToken);
          // 将用户token保存到vuex中
          _this.setToken({ token: _this.userToken });
          console.log("success");
          _this.$router.push("/home");
          console.log("success");
          alert("登陆成功");
        }
      });
    },
    loginAction(callback) {
      let _this = this;
      this.$ajax
        .post("/login", {
          params: {
            username: _this.loginForm.username,
            password: _this.loginForm.password
          }
        })
        .then(function(res) {
          callback(false, res.data);
        })
        .catch(function(error) {
          callback(error, false);
        });
    }
  }
};
</script>