<template>
  <div id="components-form-demo-advanced-search">
    <a-form layout="inline" class="ant-advanced-search-form" :form="form" @submit="onSearch">
      <a-form-item :label="`门店名称`">
        <a-input v-decorator="[
                `shopName`
                ]" />
      </a-form-item>
      <a-form-item :label="`手机号码`">
        <a-input v-decorator="[
                `phone`
                ]" />
      </a-form-item>
      <a-form-item :label="`订单编号`">
        <a-input v-decorator="[
                `orderNo`
                ]" />
      </a-form-item>
      <a-form-item :label="`日期`">
        <a-date-picker
          style="width: 100%"
          @change="onDateChange"
          :defaultValue="$moment(dateCurrent, dateFormat)"
          :format="dateFormat"
        />
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
        <span
          slot="create_time_str"
          slot-scope="text, record"
        >{{$moment(record.create_time*1000).format('YYYY-MM-DD HH:mm:ss')}}</span>
      </a-table>
    </div>
  </div>
</template>

<script>
// import moment from "moment";

const columns = [
  {
    title: "ID",
    dataIndex: "id",
    key: "id"
  },
  {
    title: "订单编号",
    dataIndex: "order_no",
    key: "order_no"
  },
  {
    title: "门店名称",
    dataIndex: "shop_name",
    key: "shop_name"
  },
  {
    title: "门店编码",
    dataIndex: "shop_id",
    key: "shop_id"
  },
  {
    title: "客户手机号",
    dataIndex: "customer_phone",
    key: "customer_phone"
  },
  {
    title: "订单金额",
    dataIndex: "goods_amount",
    key: "goods_amount"
  },
  {
    title: "活动金额",
    dataIndex: "activity_amount",
    key: "activity_amount"
  },
  {
    title: "优惠券金额",
    dataIndex: "coupon_amount",
    key: "coupon_amount"
  },
  {
    title: "实收金额",
    dataIndex: "final_amount",
    key: "final_amount"
  },
  {
    title: "门店折扣金额",
    dataIndex: "discount_amount",
    key: "discount_amount"
  },
  {
    title: "支付状态",
    dataIndex: "pay_status",
    key: "pay_status"
  },
  {
    title: "支付类型",
    dataIndex: "pay_type",
    key: "pay_type"
  },
  {
    title: "实付金额",
    dataIndex: "pay_amount",
    key: "pay_amount"
  },
  {
    title: "实付优惠金额",
    dataIndex: "pay_reduce_amount",
    key: "pay_reduce_amount"
  },
  {
    title: "退款状态",
    dataIndex: "refund_status",
    key: "refund_status"
  },
  {
    title: "代金券金额",
    dataIndex: "voucher_amount",
    key: "voucher_amount"
  },
  {
    title: "创建时间",
    dataIndex: "create_time",
    key: "create_time",
    scopedSlots: { customRender: "create_time_str" }
  }
];

export default {
  name: "SsdOrderList",
  data() {
    return {
      loading: false,
      data: [],
      error: null,
      columns: columns,
      searchValue: {
        orderNo: "",
        shopName: "",
        phone: "",
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
    // moment,
    //日期格式化
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
    //   时间格式化
    dateFormat(time) {
      var date = new Date(time);
      var year = date.getFullYear();
      /* 在日期格式中，月份是从0开始的，因此要加0
       * 使用三元表达式在小于10的前面加0，以达到格式统一  如 09:11:05
       * */
      var month =
        date.getMonth() + 1 < 10
          ? "0" + (date.getMonth() + 1)
          : date.getMonth() + 1;
      var day = date.getDate() < 10 ? "0" + date.getDate() : date.getDate();
      var hours =
        date.getHours() < 10 ? "0" + date.getHours() : date.getHours();
      var minutes =
        date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes();
      var seconds =
        date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds();
      // 拼接
      return (
        year +
        "-" +
        month +
        "-" +
        day +
        " " +
        hours +
        ":" +
        minutes +
        ":" +
        seconds
      );
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
      this.searchValue.orderNo = this.form.getFieldsValue().orderNo;
      this.searchValue.shopName = this.form.getFieldsValue().shopName;
      this.searchValue.phone = this.form.getFieldsValue().phone;
      //默认当天日期
      if (this.searchValue.date == "") {
        this.searchValue.date = this.dateCurrent;
      }
      this.getOrderList((err, data) => {
        this.loading = false;
        if (err) {
          this.error = err.toString();
        } else {
          this.data = data;
        }
      });
    },
    getOrderList(callback) {
      var orderNo = this.searchValue.orderNo;
      var shopName = this.searchValue.shopName;
      var phone = this.searchValue.phone;
      var date = this.searchValue.date;
      this.$ajax
        .get("http://localhost:39493/get-ssd-orders", {
          params: {
            order_no: orderNo,
            shop_name: shopName,
            phone: phone,
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
