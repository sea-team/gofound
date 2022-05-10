<template>
  <div id="dashboard">
    <el-card style="margin-bottom:20px;">
      <div style="display:flex;">
        <div>
          <el-label>数据库：</el-label>
          <el-select v-model="currentDB">
            <el-option
                v-for="item in  databases "
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </div>
        <template v-if="db">
          <div><span>存储路径：</span>
            <el-tag v-text="db.IndexPath"></el-tag>
          </div>
        </template>
        <div v-cloak>
          <el-button style="margin-left:10px;" type="primary" @click="search()">刷新</el-button>
        </div>
      </div>

    </el-card>
    <el-card>
      <template #header>

        <div class="header">

          <div style="display:flex;justify-content:space-between">
            <div>
              <el-input @keyup.enter="search()" v-model="params.query" placeholder="输入关键字,Enter">
                <template #append>
                  <el-button type="primary" @click="search()">查询</el-button>
                </template>
              </el-input>
            </div>

            <div style="margin-left: 10px;">
              <el-label>关键词高亮：</el-label>
              <el-switch
                  v-model="params.highlight"
              />
            </div>
          </div>

          <div v-if="data&&data.pageCount>1">
            <el-pagination @size-change="sizeChange"
                           @current-change="currentChange"
                           layout="total, sizes, prev, pager, next, jumper"
                           small="small" :page-sizes="[10,20,30,50,100, 200, 300, 500]" background
                           :page-size="params.pageSize" :total="data.total"
            />
          </div>
        </div>

      </template>
      <div v-if="data">
        <el-alert type="success" style="margin-bottom:10px;">
          <div>
            查询耗时：{{ data.time }}ms，找到：{{ data.total }}个结果，每页{{ params.limit }}条，总共{{ data.pageCount }}页
          </div>
          <div>
            <b>分词结果：</b>
          </div>
          <div>
            <el-tag style="margin-right: 10px;" v-for="item in data.words" v-text="item" :key="item"></el-tag>
          </div>
        </el-alert>
      </div>


      <el-table :stripe="true" v-loading="loading" :data="tableData" style="width: 100%">

        <el-table-column fixed type="expand" prop="document" label="Document" width="100">
          <template #default="scope">
            <json-viewer
                :value="scope.row.document"
                :expand-depth=5
                copyable
                boxed
                sort
            ></json-viewer>
          </template>
        </el-table-column>

        <el-table-column prop="id" label="ID" width="120"/>
        <el-table-column prop="score" label="Score" width="80">
          <template #default="scope">
            <el-tag v-text="scope.row.score"></el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="text" label="Text">
          <template #default="scope">
            <span v-html="scope.row.text"></span>
          </template>
        </el-table-column>


        <el-table-column fixed="right" prop="operation" label="Operation" width="100">
          <template #default="scope">
            <el-link @click="updateRow(scope.row)" type="primary" style="margin-right:10px;">更新</el-link>
            <el-link @click="deleteRow(scope.row)" type="danger">删除</el-link>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="data&&data.pageCount>1" style="margin-top:10px;">
        <el-pagination @size-change="sizeChange"
                       @current-change="currentChange"
                       layout="total, sizes, prev, pager, next, jumper"
                       small="small" :page-sizes="[10,20,30,50,100, 200, 300, 500]" background
                       :page-size="params.pageSize" :total="data.total"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import api from '../api'
import jsonViewer from 'vue-json-viewer'

export default {
  name: 'dashboard',
  components: { jsonViewer },
  data() {
    return {
      dbs: {},
      currentDB: '',
      loading: false,
      params: {
        query: '',
        page: 1,
        limit: 10,
        highlight: true,
        order: 'DESC',
      },
      data: null,
    }
  },
  watch: {
    currentDB(val) {
      this.search()
    },
    "params.highlight"(val) {
      this.search()
    }
  },
  computed: {
    databases() {
      let rs = []
      for (let db in this.dbs) {
        rs.push({
          label: db,
          value: db,
        })
      }
      return rs
    },
    db() {
      if (!this.currentDB || !this.dbs) return null
      return this.dbs[this.currentDB]
    },
    tableData() {
      if (!this.data) return null
      return this.data.documents
    },
  },
  created() {
    this.getDatabases()
  },
  methods: {
    sizeChange(size) {
      this.params.limit = size
      this.$nextTick(() => this.search())
    },
    currentChange(page) {
      this.params.page = page
      this.$nextTick(() => this.search())
    },
    getDatabases() {
      api.getDatabase().then(res => {
        if (!res.status) {
          this.$message.error(res.message)
        }
        this.dbs = res.data.data
        //选中第一项
        this.$nextTick(() => {
          this.currentDB = this.databases[0].value
        })
      })
    },
    search() {
      this.loading = true
      api.query(this.currentDB, this.params).then(res => {
        if (!res.status) {
          this.$message.error(res.message)
        }
        console.log(res.data.data)
        this.data = res.data.data
      }).finally(() => {
        this.loading = false
      })
    },
    updateRow(row) {

    },
    deleteRow(row) {
      let self = this
      //删除
      //弹出确认框
      this.$confirm('确定删除吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        api.remove(this.currentDB, row.id).then(res => {
          console.log(res.data)
          if (!res.data.state) {
            this.$message.error(res.data.message)
            return
          }
          this.$message.success('删除成功')
          self.search()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除',
        })
      })
    },
  },
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
}
</style>
