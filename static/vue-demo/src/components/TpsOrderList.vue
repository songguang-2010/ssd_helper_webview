<template>
  <div id="components-form-demo-advanced-search">
    <a-form layout="inline" class="ant-advanced-search-form" :form="form" @submit="onSearch">
      <a-form-item :label="`订单号码`">
        <a-input v-decorator="[
                `order_no`
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
        <span slot="pay_type_str" slot-scope="text, record">{{payTypeMap[record.pay_type]}}</span>
        <span slot="pay_status_str" slot-scope="text, record">{{payStatusMap[record.pay_status]}}</span>
        <span
          slot="third_pay_status_str"
          slot-scope="text, record"
        >{{thirdPayStatusMap[record.pay_status]}}</span>
        <span
          slot="rev_status_str"
          slot-scope="text, record"
        >{{revStatusMap[record.reverse_status]}}</span>
        <span
          slot="third_rev_status_str"
          slot-scope="text, record"
        >{{thirdRevStatusMap[record.third_reverse_status]}}</span>
        <span slot="ref_status_str" slot-scope="text, record">{{refStatusMap[record.refund_status]}}</span>
        <span
          slot="third_ref_status_str"
          slot-scope="text, record"
        >{{thirdRefStatusMap[record.third_refund_status]}}</span>
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
    title: "订单编号",
    dataIndex: "order_code",
    key: "order_code"
  },
  {
    title: "三方交易单号",
    dataIndex: "third_trade_no",
    key: "third_trade_no"
  },
  {
    title: "支付方式",
    dataIndex: "pay_type",
    key: "pay_type",
    scopedSlots: { customRender: "pay_type_str" }
  },
  {
    title: "支付状态",
    dataIndex: "pay_status",
    key: "pay_status",
    scopedSlots: { customRender: "pay_status_str" }
  },
  {
    title: "三方支付状态",
    dataIndex: "third_pay_status",
    key: "third_pay_status",
    scopedSlots: { customRender: "third_pay_status_str" }
  },
  {
    title: "退单状态",
    dataIndex: "reverse_status",
    key: "reverse_status",
    scopedSlots: { customRender: "rev_status_str" }
  },
  {
    title: "三方退单状态",
    dataIndex: "third_reverse_status",
    key: "third_reverse_status",
    scopedSlots: { customRender: "third_rev_status_str" }
  },
  {
    title: "退款状态",
    dataIndex: "refund_status",
    key: "refund_status",
    scopedSlots: { customRender: "ref_status_str" }
  },
  {
    title: "三方退款状态",
    dataIndex: "third_refund_status",
    key: "third_refund_status",
    scopedSlots: { customRender: "third_ref_status_str" }
  }
];

export default {
  name: "TpsOrderList",
  data() {
    return {
      loading: false,
      data: [],
      error: null,
      columns: columns,
      searchValue: {
        order_no: ""
      },
      form: this.$form.createForm(this, { name: "advanced_search" }),
      payTypeMap: {
        1: "微信付款码",
        9: "支付宝付款码"
      },
      payStatusMap: {
        0: "未支付",
        1: "支付成功",
        2: "支付结果未知",
        3: "支付失败"
      },
      thirdPayStatusMap: {
        0: "未支付",
        1: "支付成功",
        2: "支付结果未知",
        3: "支付失败"
      },
      revStatusMap: {
        0: "未撤销",
        1: "撤销成功",
        2: "撤销失败"
      },
      thirdRevStatusMap: {
        0: "未撤销",
        1: "撤销成功",
        2: "撤销失败"
      },
      refStatusMap: {
        0: "未退款",
        1: "退款中",
        2: "退款成功",
        3: "退款失败"
      },
      thirdRefStatusMap: {
        0: "未退款",
        1: "退款中",
        2: "退款成功",
        3: "退款失败"
      }
    };
  },
  created() {
    // fetch the data when the view is created and the data is
    // already being observed
    // this.fetchData();
  },
  //   watch: {
  //     // call again the method if the route changes
  //     $route: "fetchData"
  //   },
  methods: {
    onSearch() {
      this.fetchData();
    },
    fetchData() {
      this.error = this.post = null;
      this.loading = true;
      this.searchValue.order_no = this.form.getFieldsValue().order_no;
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
      var order_no = this.searchValue.order_no;
      this.$ajax
        .get("http://localhost:39493/get-tps-orders", {
          params: {
            order_no: order_no
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
