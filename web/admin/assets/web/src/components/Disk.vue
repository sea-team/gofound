<template>
  <el-row>
    <el-col :span="12">
      <el-row>
        <el-col :span="12">总共</el-col>
        <el-col :span="12"><span v-text="disk.total"></span></el-col>
        <el-col :span="12">使用</el-col>
        <el-col :span="12"><span v-text="disk.used"></span></el-col>
        <el-col :span="12">剩余</el-col>
        <el-col :span="12"><span v-text="disk.free"></span></el-col>
        <el-col :span="12">文件系统</el-col>
        <el-col :span="12"><span v-text="disk.fsType"></span></el-col>
        <el-col :span="12">ROOT</el-col>
        <el-col :span="12"><span v-text="disk.path"></span></el-col>
      </el-row>
    </el-col>
    <el-col :span="12">
      <ProgressChat :percent="disk.usedPercent"></ProgressChat>
    </el-col>
  </el-row>
</template>

<script>
import ProgressChat from '../components/ProgressChat.vue'
export default {
  name: 'Disk',
  components: {
    ProgressChat
  },
  props: {
    disk: {
      type: Object,
      required: true
    }
  },
}
</script>

<style scoped>

</style>
