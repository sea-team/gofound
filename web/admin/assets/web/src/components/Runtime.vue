<template>

  <el-row>
    <el-col :span="12">
      <el-row>
        <el-col :span="12">os</el-col>
        <el-col :span="12"><span v-text="system.os"></span></el-col>

        <el-col :span="12">arch</el-col>
        <el-col :span="12"><span v-text="system.arch"></span></el-col>

        <el-col :span="12">go version</el-col>
        <el-col :span="12"><span v-text="system.version"></span></el-col>

        <el-col :span="12">pid</el-col>
        <el-col :span="12"><el-tag v-text="system.pid"></el-tag></el-col>

        <el-col :span="12">gzip</el-col>
        <el-col :span="12"><el-switch disabled v-model="system.enableGzip"></el-switch></el-col>

      </el-row>
    </el-col>
    <el-col :span="12">
      <el-row>
        <el-col :span="12">shard</el-col>
        <el-col :span="12"><span v-text="system.shard"></span></el-col>
        <el-col :span="12">debug</el-col>
        <el-col :span="12"><el-switch disabled v-model="system.debug"></el-switch></el-col>

        <el-col :span="12">GOMAXPROCS</el-col>
        <el-col :span="12"><span v-text="system.gomaxprocs"></span></el-col>
        <el-col :span="12">goroutines</el-col>
        <el-col :span="12"><span v-text="system.goroutines"></span></el-col>
      </el-row>
    </el-col>
    <el-divider content-position="left">认证</el-divider>
    <el-row>
      <el-col :span="12">状态</el-col>
      <el-col :span="12"><el-switch disabled v-model="system.enableAuth"></el-switch></el-col>
    </el-row>
    <el-alert type="warning" show-icon v-if="!system.enableAuth">
      未开启认证，将会导致数据不安全，建议开启，启动时增加参数：--auth=admin:123456
    </el-alert>

    <el-divider content-position="left">路径</el-divider>
    <el-row>
      <el-col :span="6">数据目录</el-col>
      <el-col :span="18">
        <el-tag style="margin-right: 10px;" type="success"><span v-text="system.dataSize"></span>MB</el-tag>
        <el-tag v-text="system.dataPath"></el-tag>
      </el-col>

      <el-col :span="6">词典目录</el-col>
      <el-col :span="18"><el-tag v-text="system.dictionaryPath"></el-tag></el-col>

      <el-col :span="6">启动路径</el-col>
      <el-col :span="18"><el-tag v-text="system.executable"></el-tag></el-col>
    </el-row>
    <el-divider content-position="left">数据库</el-divider>
    <el-row>
      <el-col :span="6">数据库数量</el-col>
      <el-col :span="18"><el-tag v-text="system.dbs"></el-tag></el-col>

      <el-col :span="6">索引数量</el-col>
      <el-col :span="18"><el-tag v-text="system.indexCount"></el-tag></el-col>

      <el-col :span="6">文档数量</el-col>
      <el-col :span="18"><el-tag v-text="system.documentCount"></el-tag></el-col>


    </el-row>

  </el-row>


</template>

<script>
export default {
  name: 'Runtime',
  props: {
    system: {
      type: Object,
      required: true,
    },
  },
}
</script>

<style scoped>

</style>
