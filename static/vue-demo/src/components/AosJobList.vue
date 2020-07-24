<template>
  <div id="components-form-demo-advanced-search">
    <!-- <a-form layout="inline" class="ant-advanced-search-form" :form="form" @submit="onSearch">
      <a-form-item :label="`任务类型`">
        <a-select 
          defaultValue="1"
          @change="handleJobType"
        >
          <a-select-option value="1">
            订单详情同步到统计服务
          </a-select-option>
          <a-select-option value="2">
            校验POS订单同步数量
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item>
        <a-button type="primary" @click="onSearch">Search</a-button>
      </a-form-item>
    </a-form> -->

    <div v-if="loading" class="loading">Loading...</div>

    <div v-if="error" class="error">{{ error }}</div>

    <div v-if="data" class="content">

      <a-table
        :columns="columns"
        :dataSource="data"
        :pagination="pagination"
        @change="handleTableChange"
        :rowSelection="rowSelection"
      >
        <span
          slot="type_str"
          slot-scope="text, record"
        >{{typeMap[record.type]}}</span>
        <span slot="status_str" slot-scope="text, record">{{statusMap[record.status]}}</span>
        
        <span
          slot="create_time_str"
          slot-scope="text, record"
        >{{$moment(record.create_time*1000).format('YYYY-MM-DD HH:mm:ss')}}</span>
        <span
          slot="update_time_str"
          slot-scope="text, record"
        >{{$moment(record.update_time*1000).format('YYYY-MM-DD HH:mm:ss')}}</span>
        
      </a-table>

    </div>
  </div>
</template>

<script>

//列数据描述
const columns = function() {
  let { filteredInfo } = this;
  filteredInfo = filteredInfo || {};
  const columns = [
    {
      title: "任务ID",
      dataIndex: "id",
      key: "id"
    },
    {
      title: "参数",
      dataIndex: "params",
      key: "params"
    },
    {
      title: "描述",
      dataIndex: "description",
      key: "description"
    },
    {
      title: "创建时间",
      dataIndex: "create_time",
      key: "create_time",
      scopedSlots: { customRender: "create_time_str" }
    },
    {
      title: "更新时间",
      dataIndex: "update_time",
      key: "update_time",
      scopedSlots: { customRender: "update_time_str" }
    },
    {
      title: "状态",
      dataIndex: "status",
      key: "status",
      scopedSlots: { customRender: "status_str" },
      filters: [
        { text: "尚未执行", value: "0" },
        { text: "执行中", value: "1" },
        { text: "执行成功", value: "2" },
        { text: "执行失败", value: "3" }
      ],
      filteredValue: filteredInfo.status || null,
      onFilter: (value, record) => record.status.includes(value),
      filterMultiple: false
    },
    {
      title: "类型",
      dataIndex: "type",
      key: "type",
      scopedSlots: { customRender: "type_str" },
      filters: [
        { text: "订单详情同步到统计服务", value: "1" },
        { text: "校验POS订单同步数量", value: "2" }
      ],
      filteredValue: filteredInfo.type || null,
      onFilter: (value, record) => record.type.includes(value)
    }
  ];
  return columns;
};

//列表达翻译
const columnsMaps = {
  typeMap: {
    1: "订单详情同步到统计服务",
    2: "校验POS订单同步数量"
  },
  statusMap: {
    0: "尚未执行",
    1: "执行中",
    2: "执行成功",
    3: "执行失败"
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
  name: "AosJobList",
  data() {
    return {
      pagination: pagination,
      loading: false,
      data: [],
      error: null,
      searchValue: {
        // type: "1"
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
      this.fetchData();
    },
    // onShowSizeChange(current, pageSize) {
    //   console.log(current, pageSize);
    // },
    clearFilters() {
      console.log(this.filteredInfo);
      this.filteredInfo = null;
    },
    // handleJobType(value){
    //   this.searchValue.type = value
    //   console.log("type change")
    // },
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
    fetchData() {
      this.error = this.post = null;
      this.loading = true;

      this.getJobList((err, data) => {
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
    getJobList(callback) {
      console.log("Various parameters", this.pagination, this.filteredInfo);
      let _this = this;
    //   var type = this.searchValue.type;
      // var canary = this.searchValue.canary || 0;
      this.$ajax
        .get("/get-aos-jobs", {
          params: {
            // type: type,
            pageSize: _this.pagination.pageSize,
            pageNum: _this.pagination.current,
            status: ((_this.filteredInfo==null || _this.filteredInfo.status==null || _this.filteredInfo.status.length==0) ? 4 : _this.filteredInfo.status[0]),
            type: ((_this.filteredInfo==null || _this.filteredInfo.type==null || _this.filteredInfo.type.length==0) ? 0 : _this.filteredInfo.type[0])
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
