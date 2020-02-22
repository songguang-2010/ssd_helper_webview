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
      <a-form-item :label="`门店编码`">
        <a-input v-decorator="[
                `shopCode`
                ]" />
      </a-form-item>
      <a-form-item :label="`日期`">
        <a-date-picker
          style="width: 100%"
          @change="onDateChange"
          :defaultValue="moment(dateCurrent, dateFormat)"
          :format="dateFormat"
        />
      </a-form-item>
      <a-form-item :label="`产品名称`">
        <a-input v-decorator="[
                `prodName`
                ]" />
      </a-form-item>

      <a-form-item>
        <a-button type="primary" @click="onSearch">Search</a-button>
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
import moment from "moment";

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
  name: "SkuSpecList",
  data() {
    return {
      loading: false,
      data: [],
      error: null,
      columns: columns,
      searchValue: {
        shopName: "",
        shopCode: "",
        date: "",
        prodName: ""
      },
      form: this.$form.createForm(this, { name: "advanced_search" }),
      dateFormat: "YYYY-MM-DD",
      dateCurrent: ""
    };
  },
  created() {
    this.dateCurrent = this.getNowFormatDate();
    // fetch the data when the view is created and the data is
    // already being observed
    this.fetchData();
  },
  watch: {
    // call again the method if the route changes
    $route: "fetchData"
  },
  methods: {
    moment,
    getNowFormatDate() {
      var date = new Date();
      var seperator1 = "-";
      var year = date.getFullYear();
      var month = date.getMonth() + 1;
      var strDate = date.getDate();
      if (month >= 1 && month <= 9) {
        month = "0" + month;
      }
      if (strDate >= 0 && strDate <= 9) {
        strDate = "0" + strDate;
      }
      var currentdate = year + seperator1 + month + seperator1 + strDate;
      return currentdate;
    },
    onDateChange(date, dateString) {
      console.log(date, dateString);
      this.searchValue.date = dateString;
    },
    onSearch() {
      this.fetchData();
    },
    fetchData() {
      this.error = this.post = null;
      this.loading = true;
      this.searchValue.shopName = this.form.getFieldsValue().shopName;
      this.searchValue.shopCode = this.form.getFieldsValue().shopCode;
      this.searchValue.prodName = this.form.getFieldsValue().prodName;
      //默认当天日期
      if (this.searchValue.date == "") {
        this.searchValue.date = this.dateCurrent;
      }
      // replace `getPost` with your data fetching util / API wrapper
      this.getSpecList((err, data) => {
        this.loading = false;
        if (err) {
          this.error = err.toString();
        } else {
          this.data = data;
        }
      });
    },
    getSpecList(callback) {
      var shopName = this.searchValue.shopName;
      var shopCode = this.searchValue.shopCode;
      var prodName = this.searchValue.prodName;
      var date = this.searchValue.date;
      this.$ajax
        .get("/get-sku-specs", {
          params: {
            shop_name: shopName,
            shop_code: shopCode,
            prod_name: prodName,
            date: date
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
