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
      <a-form-item :label="`分公司编码`">
        <a-input v-decorator="[
                `company_sale_id`
                ]" />
      </a-form-item>
      <a-form-item :label="`设备Model`">
        <a-input v-decorator="[
                `model`
                ]" />
      </a-form-item>
      <a-form-item :label="`排除测试数据`">
        <a-select 
          defaultValue="1"
          @change="handleTestExcludeChange"
        >
          <a-select-option value="1">
            是
          </a-select-option>
          <a-select-option value="0">
            否
          </a-select-option>
        </a-select>
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

      <div style="margin-bottom: 16px">
        <a-button type="primary" @click="batchSetCanary" :loading="loading">设为灰度设备</a-button>
        <a-button type="primary" @click="batchCancelCanary" :loading="loading">取消灰度设备</a-button>
        <span style="margin-left: 8px">
          <template v-if="hasSelected">{{`选择了 ${selectedRowKeys.length} 个条目`}}</template>
        </span>
      </div>

      <a-table
        :columns="columns"
        :dataSource="data"
        :pagination="pagination"
        @change="handleTableChange"
        :rowSelection="rowSelection"
      >
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
      <br />
      <!-- <a-pagination
        showSizeChanger
        :pageSize.sync="pageSize"
        :pageSizeOptions="pageSizeOptions"
        @showSizeChange="onShowSizeChange"
        :total="total"
        v-model="current"
      />-->
    </div>
  </div>
</template>

<script>
// const columns = [];
// const rowSelection = {
//     onChange: (selectedRowKeys, selectedRows) => {
//       console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
//     },
//     onSelect: (record, selected, selectedRows) => {
//       console.log(record, selected, selectedRows);
//       record.is_canary = 1

//       // const newData = [...this.data];
//       // const target = newData.filter(item => record.id === item.id)[0];
//       // target.is_canary = 0;
//       // this.data = newData;
//     },
//     onSelectAll: (selected, selectedRows, changeRows) => {
//       console.log(selected, selectedRows, changeRows);
//     },
//   };

//版本号过滤枚举结构
let filtersAppVersion = [
        // { text: "是", value: "1" },
        // { text: "否", value: "0" }
      ];

//列数据描述
const columns = function() {
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
      title: "分公司编码",
      dataIndex: "company_sale_id",
      key: "company_sale_id"
    },
    {
      title: "App版本号",
      dataIndex: "app_version",
      key: "app_version",
      filters: filtersAppVersion,
      filteredValue: filteredInfo.app_version || null,
      onFilter: (value, record) => record.app_version.includes(value),
      filterMultiple: false
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
      onFilter: (value, record) => record.app_env.includes(value),
      filterMultiple: false
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
      onFilter: (value, record) => record.is_canary.includes(value),
      filterMultiple: false
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
      title: "操作",
      key: "action",
      scopedSlots: { customRender: "action" }
    }
  ];
  return columns;
};

//列表达翻译
const columnsMaps = {
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
  }
};

//分页
const pagination = {
  pageSize: 20,
  current: 1,
  pageSizeOptions: ["10", "20", "50", "100", "200", "500"],
  total: 0,
  showSizeChanger: true
};

//导出对象
export default {
  name: "MiscDeviceList",
  data() {
    return {
      pagination: pagination,
      loading: false,
      data: [],
      error: null,
      searchValue: {
        shop_no: "",
        shop_name: "",
        app_version: "",
        company_sale_id: "",
        model: "",
        // 是否排除测试数据
        test_exclude: "1"
        // canary: 0
      },
      form: this.$form.createForm(this, { name: "advanced_search" }),
      ...columnsMaps,
      filteredInfo: null,
      selectedRowKeys: [],
      selectedRows: []
      // rowSelection,
    };
  },
  computed: {
    hasSelected() {
      return this.selectedRowKeys.length > 0;
    },
    rowSelection() {
      const { selectedRowKeys, selectedRows } = this;
      return {
        selectedRowKeys,
        selectedRows,
        // onChange: this.onSelectChange,
        onChange: (selectedRowKeys, selectedRows) => {
          console.log(
            `selectedRowKeys: ${selectedRowKeys}`,
            "selectedRows: ",
            selectedRows
          );
          this.selectedRowKeys = selectedRowKeys;
          this.selectedRows = selectedRows;
        }
      };
    },
    columns: columns
  },
  created() {
    // fetch the data when the view is created and the data is
    // already being observed
    this.routeChange();
  },
  watch: {
    // call again the method if the route changes
    $route: "routeChange"
    // pageSize(val) {
    //   console.log('pageSize', val);
    // },
    // current(val) {
    //   console.log('current', val);
    // },
  },
  methods: {
    routeChange(){
      this.fetchCanaryAppVersions();
      this.fetchData();
    },
    // onShowSizeChange(current, pageSize) {
    //   console.log(current, pageSize);
    // },
    clearFilters() {
      console.log(this.filteredInfo);
      this.filteredInfo = null;
    },
    handleTestExcludeChange(value){
      this.searchValue.test_exclude = value
      console.log("test exclude change")
    },
    // table change event
    handleTableChange(pagination, filters, sorter) {
      console.log("Various parameters", pagination, filters, sorter);
      this.filteredInfo = filters;
      // this.sortedInfo = sorter;
      const pager = { ...this.pagination };
      pager.current = pagination.current;
      pager.pageSize = pagination.pageSize;
      this.pagination = pager;
      this.fetchData();
    },
    onSearch() {
      this.fetchData();
    },
    fetchCanaryAppVersions() {
      this.error = null;
      this.loading = true;
      this.getCanaryAppVersions((err, data) => {
        this.loading = false;
        if (err) {
          this.error = err.toString();
        } else {
          filtersAppVersion = [];
          if(data.length>0){
            for(var i=0;i<data.length;i++){
              filtersAppVersion[i] = {
                text: data[i].app_version,
                value: data[i].app_version
              };
            }
          }
          console.log("filtersAppVersion:", filtersAppVersion)
        }
      });
    },
    getCanaryAppVersions(callback) {
      let _this = this;
      this.$ajax
        .get("/get-canary-app-versions")
        .then(function(res) {
          callback(false, res.data);
        })
        .catch(function(error) {
          callback(error, false);
        });
    },
    fetchData() {
      this.error = this.post = null;
      this.loading = true;
      this.searchValue.shop_no = this.form.getFieldsValue().shop_no;
      this.searchValue.shop_name = this.form.getFieldsValue().shop_name;
      this.searchValue.app_version = this.form.getFieldsValue().app_version;
      this.searchValue.company_sale_id = this.form.getFieldsValue().company_sale_id;
      this.searchValue.model = this.form.getFieldsValue().model;
      // this.searchValue.test_exclude = this.form.getFieldsValue().test_exclude;
      // this.searchValue.canary = this.form.getFieldsValue().canary;
      this.getDeviceList((err, data) => {
        this.loading = false;
        this.selectedRowKeys = [];
        // this.clearFilters();
        if (err) {
          this.error = err.toString();
        } else {
          this.data = data.list;
          const pagination = { ...this.pagination };
          pagination.total = data.total;
          this.pagination = pagination;
        }
      });
    },
    getDeviceList(callback) {
      console.log("Various parameters", this.pagination, this.filteredInfo);
      let _this = this;
      var shop_no = this.searchValue.shop_no;
      var shop_name = this.searchValue.shop_name;
      var app_version = this.searchValue.app_version;
      var company_sale_id = this.searchValue.company_sale_id;
      var model = this.searchValue.model;
      var test_exclude = this.searchValue.test_exclude;
      // var canary = this.searchValue.canary || 0;
      this.$ajax
        .get("/get-misc-devices", {
          params: {
            shop_no: shop_no,
            shop_name: shop_name,
            app_version: app_version,
            company_sale_id: company_sale_id,
            model: model,
            test_exclude: test_exclude,
            pageSize: _this.pagination.pageSize,
            pageNum: _this.pagination.current,
            is_canary: ((_this.filteredInfo==null || _this.filteredInfo.is_canary==null || _this.filteredInfo.is_canary.length==0) ? 2 : _this.filteredInfo.is_canary[0]),
            app_env: ((_this.filteredInfo==null || _this.filteredInfo.app_env==null || _this.filteredInfo.app_env.length==0) ? 'all' : _this.filteredInfo.app_env[0])
          }
        })
        .then(function(res) {
          callback(false, res.data);
        })
        .catch(function(error) {
          callback(error, false);
        });
    },
    batchSetCanary() {
      if (this.selectedRowKeys.length == 0) {
        alert("请至少选择一个条目");
        return false;
      }

      this.loading = true;
      let selectedIds = "";
      for (var i = 0; i < this.selectedRows.length; i++) {
        selectedIds += "," + this.selectedRows[i].id;
      }
      //去除前端空格及逗号
      selectedIds = selectedIds.replace(/^(\s|,)+/g, "");

      let _this = this;
      this.batchSetCanaryAction(selectedIds, (err, data) => {
        _this.loading = false;
        _this.selectedRowKeys = [];
        if (err) {
          _this.error = err.toString();
          alert(_this.error);
        } else {
          console.log("response data from go server: " + data);
          for (var i = 0; i < _this.selectedRows.length; i++) {
            _this.selectedRows[i].is_canary = 1;
          }
          alert("设置成功");
          _this.fetchData()
        }
      });
    },
    batchCancelCanary() {
      if (this.selectedRowKeys.length == 0) {
        alert("请至少选择一个条目");
        return false;
      }

      this.loading = true;
      let selectedIds = "";
      for (var i = 0; i < this.selectedRows.length; i++) {
        selectedIds += "," + this.selectedRows[i].id;
      }
      //去除前端空格及逗号
      selectedIds = selectedIds.replace(/^(\s|,)+/g, "");

      let _this = this;
      this.batchCancelCanaryAction(selectedIds, (err, data) => {
        _this.loading = false;
        _this.selectedRowKeys = [];
        if (err) {
          _this.error = err.toString();
          alert(_this.error);
        } else {
          console.log("response data from go server: " + data);
          for (var i = 0; i < _this.selectedRows.length; i++) {
            _this.selectedRows[i].is_canary = 0;
          }
          alert("取消成功");
          _this.fetchData()
        }
      });
    },
    batchSetCanaryAction(device_ids, callback) {
      let _this = this;
      this.$ajax
        .post("/set-canary-batch", {
          device_ids: device_ids
        })
        .then(function(res) {
          console.log("response from go server: " + res);
          callback(false, res.data);
        })
        .catch(function(error) {
          callback(error, false);
        });
    },
    batchCancelCanaryAction(device_ids, callback) {
      let _this = this;
      this.$ajax
        .post("/cancel-canary-batch", {
          device_ids: device_ids
        })
        .then(function(res) {
          console.log("response from go server: " + res);
          callback(false, res.data);
        })
        .catch(function(error) {
          callback(error, false);
        });
    },
    cancelCanary(id) {
      const newData = [...this.data];
      const target = newData.filter(item => id === item.id)[0];
      if (target) {
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
            _this.fetchData()
          }
        });
      }
    },
    //设置为灰度设备
    setCanary(id) {
      const newData = [...this.data];
      const target = newData.filter(item => id === item.id)[0];
      if (target) {
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
            _this.fetchData()
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
    }
  }

  // props: {
  //     content: String
  // }
};
</script>

<style>
</style>
