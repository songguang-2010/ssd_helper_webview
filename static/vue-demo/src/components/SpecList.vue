<template>
  <div id="components-form-demo-advanced-search">
    <!-- <label style="background: #ff0000;">门店名称</label> -->
    <!-- <a-input-search v-model="searchValue" style="width: 400px" @search="onSearch" enterButton></a-input-search> -->
    <!-- <a-input addonBefore="门店名称" v-model="searchValue" style="width: 400px">
        <a-icon slot="addonAfter" type="search" v-on:click="onSearch" />
    </a-input>-->
    <!-- <br /><br/> -->

    <a-form layout="inline" class="ant-advanced-search-form" :form="form" @submit="onSearch">
      <a-form-item :label="`门店名称`">
        <a-input v-decorator="[
                `shopName`
                ]" />
      </a-form-item>

      <a-form-item>
        <a-button type="primary" html-type="submit">Search</a-button>
      </a-form-item>
    </a-form>

    <div v-if="loading" class="loading">Loading...</div>

    <div v-if="error" class="error">{{ error }}</div>

    <div v-if="data" class="content">
      <a-table :columns="columns" :dataSource="data">
        <!-- <span slot="customTitle"><a-icon type="smile-o" /> Name</span> -->
      </a-table>
    </div>
  </div>
</template>

<script>
const columns = [
  {
    title: "ID",
    dataIndex: "id",
    key: "id"
  },
  {
    title: "门店编码",
    dataIndex: "store_code",
    key: "store_code"
  },
  {
    title: "门店名称",
    dataIndex: "store_name",
    key: "store_name"
  },
  {
    title: "产品名称",
    dataIndex: "prod_name",
    key: "prod_name"
  },
  {
    title: "产品编码",
    dataIndex: "prod_code",
    key: "prod_code"
  },
  {
    title: "销售单位比例",
    dataIndex: "sale_unit_ratio",
    key: "sale_unit_ratio"
  },
  {
    title: "报货单位比例",
    dataIndex: "purchase_unit_ratio",
    key: "purchase_unit_ratio"
  },
  {
    title: "销售单位",
    dataIndex: "sale_unit",
    key: "sale_unit"
  },
  {
    title: "报货单位",
    dataIndex: "purchase_unit",
    key: "purchase_unit"
  }
];

export default {
  name: "SpecList",
  data() {
    return {
      loading: false,
      data: [],
      error: null,
      columns: columns,
      searchValue: "",
      form: this.$form.createForm(this, { name: "advanced_search" })
    };
  },
  created() {
    // fetch the data when the view is created and the data is
    // already being observed
    this.fetchData();
  },
  watch: {
    // call again the method if the route changes
    $route: "fetchData"
  },
  methods: {
    onSearch() {
      this.fetchData();
    },
    fetchData() {
      this.error = this.post = null;
      this.loading = true;
      // window.console.log(this.form.getFieldsValue().shopName)
      // if (this.searchValue == undefined) {
      //     this.searchValue = ""
      // }
      this.searchValue = this.form.getFieldsValue().shopName;
      // replace `getPost` with your data fetching util / API wrapper
      this.getPost(this.searchValue, (err, data) => {
        this.loading = false;
        if (err) {
          this.error = err.toString();
        } else {
          this.data = data;
        }
      });
    },
    getPost(shopName, callback) {
      this.$ajax
        .get("http://localhost:39493/get-specs", {
          params: {
            shop_name: shopName
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

  // props: {
  //     content: String
  // }
};
</script>

<style>
</style>
