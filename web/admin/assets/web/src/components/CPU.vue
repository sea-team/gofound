<template>
  <el-row>
    <el-col :span="12">
      <el-row>
        <el-col :span="12">核心</el-col>
        <el-col :span="12"><span v-text="cpu.cores"></span></el-col>
        <el-col :span="12">名称</el-col>
        <el-col :span="12"><span v-text="cpu.modelName"></span></el-col>
      </el-row>
    </el-col>
    <el-col :span="12">
      <ProgressChat :percent="cpu.usedPercent"></ProgressChat>
    </el-col>
  </el-row>
</template>

<script>
import ProgressChat from '../components/ProgressChat.vue'
export default {
  name: 'CPU',
  components: {
    ProgressChat
  },
  props: {
    cpu: {
      type: Object,
      required: true
    }
  },
}
</script>

<style scoped>

</style>
