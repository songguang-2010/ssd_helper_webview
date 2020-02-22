<template>
  <div id="components-form-demo-advanced-search">
    <a-form layout="inline" class="ant-advanced-search-form" :form="form" @submit="onSearch">
      <a-form-item :label="`门店编码`">
        <a-input v-decorator="[
                `shopNo`
                ]" />
      </a-form-item>
      <a-form-item :label="`门店名称`">
        <a-input v-decorator="[
                `shopName`
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
    title: "报货日期",
    dataIndex: "date_request",
    key: "date_request"
  },
  {
    title: "收货日期",
    dataIndex: "date_response",
    key: "date_response"
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
    title: "产品数量",
    dataIndex: "prod_number",
    key: "prod_number"
  }
];

export default {
  name: "SkuResponseList",
  data() {
    return {
      loading: false,
      data: [],
      error: null,
      columns: columns,
      searchValue: {
        shopNo: "",
        shopName: "",
        prodName: "",
        date: ""
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
      this.searchValue.shopNo = this.form.getFieldsValue().shopNo;
      this.searchValue.prodName = this.form.getFieldsValue().prodName;
      //默认当天日期
      if (this.searchValue.date == "") {
        this.searchValue.date = this.dateCurrent;
      }
      // replace `getPost` with your data fetching util / API wrapper
      this.getList((err, data) => {
        this.loading = false;
        if (err) {
          this.error = err.toString();
        } else {
          this.data = data;
        }
      });
    },
    getList(callback) {
      var shopNo = this.searchValue.shopNo;
      var shopName = this.searchValue.shopName;
      var prodName = this.searchValue.prodName;
      var date = this.searchValue.date;
      this.$ajax
        .get("/get-sku-responses", {
          params: {
            shop_name: shopName,
            shop_no: shopNo,
            prod_name: prodName,
            date_response: date
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
