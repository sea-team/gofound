<template>
  <el-button type="primary" @click="gc()">释放内存</el-button>
</template>

<script>
import api from '../api'

export default {
  name: 'GC',
  methods:{
    gc(){
      api.gc().then(res=>{
        if(res.data.state){
          this.$message.success('释放成功')
        }else {
          this.$message.error(res.data.message)
        }
      })
    }
  }
}
</script>

<style scoped>

</style>
