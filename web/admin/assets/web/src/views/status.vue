<template>
  <div v-loading="loading">
    <div style="text-align: center">
      <el-button @click="getStatus" icon="Refresh" type="primary">刷新</el-button>
    </div>
    <el-row v-if="data">

      <el-col :span="12">
        <el-row>
          <el-col :span="24">
            <el-card>
              <template #header>内存</template>
              <Memory :memory="data.memory"></Memory>
            </el-card>

          </el-col>
          <el-col :span="24">
            <el-card>
              <template #header>CPU</template>
              <CPU :cpu="data.cpu"></CPU>
            </el-card>
          </el-col>

          <el-col :span="24">
            <el-card>
              <template #header>磁盘</template>
              <Disk :disk="data.disk"></Disk>
            </el-card>
          </el-col>


        </el-row>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>运行时</template>
          <Runtime :system="data.system"></Runtime>
        </el-card>
      </el-col>



    </el-row>
  </div>
</template>

<script>
import api from '../api'
import ProgressChat from '../components/ProgressChat.vue'
import Memory from '../components/Memory.vue'
import CPU from '../components/CPU.vue'
import Disk from '../components/Disk.vue'
import Runtime from '../components/Runtime.vue'

export default {
  name: 'document.vue',
  components: { Runtime, Disk, CPU, Memory, ProgressChat },
  data() {
    return {
      data: null,
      loading: false,
    }
  },
  created() {
    this.getStatus()
  },
  methods: {
    getStatus() {
      let self = this
      self.loading = true

      api.getStatus().then(({ data }) => {

        if (data.state) {
          this.data = data.data
        } else {
          this.$message.error(data.message)
        }
      }).finally(() => self.loading = false)
    },
  },
}
</script>

<style>
.el-col {
  padding: 10px;
}
</style>
