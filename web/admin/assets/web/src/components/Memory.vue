<template>
  <el-row>
    <el-col :span="12">
      <el-row>
        <el-col :span="12">总共</el-col>
        <el-col :span="12"><span v-text="memory.total"></span>GB</el-col>
        <el-col :span="12">使用</el-col>
        <el-col :span="12"><span v-text="memory.used"></span>GB</el-col>
        <el-col :span="12">剩余</el-col>
        <el-col :span="12"><span v-text="memory.free"></span>GB</el-col>
        <el-col :span="12">GoFound使用</el-col>
        <el-col :span="12"><span v-text="memory.self"></span>GB</el-col>
      </el-row>
    </el-col>
    <el-col :span="12">
      <ProgressChat :percent="memory.usedPercent"></ProgressChat>
    </el-col>
  </el-row>
</template>

<script>
import ProgressChat from '../components/ProgressChat.vue'
export default {
  name: 'Memory',
  components:{ProgressChat},
  props: {
    memory: {
      type: Object,
      required: true
    }
  },

}
</script>

<style scoped>

</style>
