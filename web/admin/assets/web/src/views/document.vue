<template>

  <el-row>
    <el-col :span="12">
      <el-card>
        <template #header>
          码云
        </template>
        <el-link type="success" target="_blank" href="https://gitee.com/tompeppa/gofound">
          https://gitee.com/tompeppa/gofound
        </el-link>
      </el-card>
    </el-col>

    <el-col :span="12">
      <el-card>
        <template #header>
          Github
        </template>
        <el-link type="success" target="_blank" href="https://github.com/newpanjing/gofound">
          https://github.com/newpanjing/gofound
        </el-link>
      </el-card>
    </el-col>

    <el-col :span="12">
      <el-card>
        <template #header>
          社区
        </template>
        <el-link type="success" target="_blank" href="https://simpleui.72wo.com">
          https://simpleui.72wo.com
        </el-link>
      </el-card>
    </el-col>

    <el-col :span="12">
      <el-card>
        <template #header>
          QQ群
        </template>
        群号：
        <el-tag>556102631</el-tag>
      </el-card>
    </el-col>

  </el-row>


</template>

<script>
export default {
  name: 'document.vue',
}
</script>

<style scoped>
.el-col{
  padding:10px;
}
</style>
