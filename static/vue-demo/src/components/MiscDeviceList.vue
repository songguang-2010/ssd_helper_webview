<template>
  <div id="components-form-demo-advanced-search">
    <a-form layout="inline" class="ant-advanced-search-form" :form="form" @submit="onSearch">
      <a-form-item :label="`门店编码`">
        <a-input v-decorator="[
                `shop_no`
                ]" />
      </a-form-item>
      <a-form-item :label="`门店名称`">
        <a-input v-decorator="[
                `shop_name`
                ]" />
      </a-form-item>
      <a-form-item :label="`App版本`">
        <a-input v-decorator="[
                `app_version`
                ]" />
      </a-form-item>
      <!-- <a-form-item label="灰度设备" :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }">
        <a-select v-decorator="['canary', { initialValue: '0'}]" style="width: 270px;">
          <a-select-option value="0">不限</a-select-option>
          <a-select-option value="1">是</a-select-option>
          <a-select-option value="2">否</a-select-option>
        </a-select>
      </a-form-item>-->

      <a-form-item>
        <a-button type="primary" @click="onSearch">Search</a-button>
      </a-form-item>
    </a-form>

    <div v-if="loading" class="loading">Loading...</div>

    <div v-if="error" class="error">{{ error }}</div>

    <div v-if="data" class="content">
      <span v-if="filteredInfo">
        <span v-if="filteredInfo.app_env != '' && filteredInfo.app_env">
          <a-tag color="blue">灰度状态: {{appEnvMap[filteredInfo.app_env]}}</a-tag>
        </span>
        <span v-if="filteredInfo.is_canary != '' && filteredInfo.is_canary">
          <a-tag color="blue">灰度设备: {{isCanaryMap[filteredInfo.is_canary]}}</a-tag>
        </span>
      </span>

      <a-table :columns="columns" :dataSource="data" @change="handleChange">
        <span
          slot="network_type_str"
          slot-scope="text, record"
        >{{networkTypeMap[record.network_type]}}</span>
        <span slot="is_canary_str" slot-scope="text, record">{{isCanaryMap[record.is_canary]}}</span>
        <span slot="app_env_str" slot-scope="text, record">{{appEnvMap[record.app_env]}}</span>
        <span
          slot="update_time_str"
          slot-scope="text, record"
        >{{$moment(record.update_time*1000).format('YYYY-MM-DD HH:mm:ss')}}</span>
        <span slot="action" slot-scope="text, record">
          <span v-if="record.is_canary=='1'">
          <a-popconfirm title="确实要将该设备从灰度中取消么?" @confirm="() => cancelCanary(record.id)">
            <a>取消灰度</a>
          </a-popconfirm>
          </span>
          <span v-else>
          <a-popconfirm title="确实要将该设备设置为灰度设备么?" @confirm="() => setCanary(record.id)">
            <a>设为灰度</a>
          </a-popconfirm>
          </span>
        </span>
      </a-table>
    </div>
  </div>
</template>

<script>
// const columns = [];

export default {
  name: "MiscDeviceList",
  data() {
    return {
      loading: false,
      data: [],
      error: null,
      // columns: columns,
      searchValue: {
        shop_no: "",
        shop_name: "",
        app_version: ""
        // canary: 0
      },
      form: this.$form.createForm(this, { name: "advanced_search" }),
      networkTypeMap: {
        1: "网线",
        2: "WIFI",
        3: "移动网络"
      },
      isCanaryMap: {
        0: "否",
        1: "是"
      },
      appEnvMap: {
        prod: "未灰度",
        canary: "已灰度"
      },
      filteredInfo: null
    };
  },
  computed: {
    columns() {
      let { filteredInfo } = this;
      filteredInfo = filteredInfo || {};
      const columns = [
        {
          title: "设备ID",
          dataIndex: "id",
          key: "id"
        },
        {
          title: "门店编号",
          dataIndex: "shop_no",
          key: "shop_no"
        },
        {
          title: "门店名称",
          dataIndex: "shop_name",
          key: "shop_name"
        },
        {
          title: "App版本号",
          dataIndex: "app_version",
          key: "app_version"
        },
        {
          title: "更新时间",
          dataIndex: "update_time",
          key: "update_time",
          scopedSlots: { customRender: "update_time_str" }
        },
        {
          title: "灰度状态",
          dataIndex: "app_env",
          key: "app_env",
          scopedSlots: { customRender: "app_env_str" },
          filters: [
            { text: "已灰度", value: "canary" },
            { text: "未灰度", value: "prod" }
          ],
          filteredValue: filteredInfo.app_env || null,
          onFilter: (value, record) => record.app_env.includes(value)
        },
        {
          title: "灰度设备",
          dataIndex: "is_canary",
          key: "is_canary",
          scopedSlots: { customRender: "is_canary_str" },
          filters: [
            { text: "是", value: "1" },
            { text: "否", value: "0" }
          ],
          filteredValue: filteredInfo.is_canary || null,
          onFilter: (value, record) => record.is_canary.includes(value)
        },
        {
          title: "网络类型",
          dataIndex: "network_type",
          key: "network_type",
          scopedSlots: { customRender: "network_type_str" },
          filters: [
            { text: "WIFI", value: "2" },
            { text: "网线", value: "1" },
            { text: "移动网络", value: "3" }
          ],
          filteredValue: filteredInfo.network_type || null,
          onFilter: (value, record) => record.network_type.includes(value)
        },
        {
          title: '操作',
          key: 'action',
          scopedSlots: { customRender: 'action' },
        }
      ];
      return columns;
    }
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
    cancelCanary(id) {
      const newData = [...this.data];
      const target = newData.filter(item => id === item.id)[0];
      if(target){
        let _this = this;
        this.cancelCanaryAction(id, (err, data) => {
          if (err) {
            _this.error = err.toString();
            alert(_this.error);
          } else {
            console.log("response data from go server: " + data);
            target.is_canary = 0;
            _this.data = newData;
            alert("取消成功");
          }
        });
      }
      
    },
    //设置为灰度设备
    setCanary(id) {
      const newData = [...this.data];
      const target = newData.filter(item => id === item.id)[0];
      if(target){
        let _this = this;
        this.setCanaryAction(id, (err, data) => {
          if (err) {
            _this.error = err.toString();
            alert(_this.error);
          } else {
            console.log("response data from go server: " + data);
            target.is_canary = 1;
            _this.data = newData;
            alert("设置成功");
          }
        });
      }
      
    },
    setCanaryAction(device_id, callback) {
      let _this = this;
      this.$ajax
        .post("/set-canary", {
          device_id: device_id
        })
        .then(function(res) {
          console.log("response from go server: " + res);
          callback(false, res.data);
        })
        .catch(function(error) {
          callback(error, false);
        });
    },
    cancelCanaryAction(device_id, callback) {
      let _this = this;
      this.$ajax
        .post("/cancel-canary", {
          device_id: device_id
        })
        .then(function(res) {
          console.log("response from go server: " + res);
          callback(false, res.data);
        })
        .catch(function(error) {
          callback(error, false);
        });
    },
    clearFilters() {
      console.log(this.filteredInfo);
      this.filteredInfo = null;
    },
    handleChange(pagination, filters, sorter) {
      console.log("Various parameters", pagination, filters, sorter);
      this.filteredInfo = filters;
      // this.sortedInfo = sorter;
    },
    onSearch() {
      this.fetchData();
    },
    fetchData() {
      this.error = this.post = null;
      this.loading = true;
      this.searchValue.shop_no = this.form.getFieldsValue().shop_no;
      this.searchValue.shop_name = this.form.getFieldsValue().shop_name;
      this.searchValue.app_version = this.form.getFieldsValue().app_version;
      // this.searchValue.canary = this.form.getFieldsValue().canary;
      this.getDeviceList((err, data) => {
        this.loading = false;
        if (err) {
          this.error = err.toString();
        } else {
          this.data = data;
          this.clearFilters();
        }
      });
    },
    getDeviceList(callback) {
      var shop_no = this.searchValue.shop_no;
      var shop_name = this.searchValue.shop_name;
      var app_version = this.searchValue.app_version;
      // var canary = this.searchValue.canary || 0;
      this.$ajax
        .get("/get-misc-devices", {
          params: {
            shop_no: shop_no,
            shop_name: shop_name,
            app_version: app_version
            // canary: canary
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
