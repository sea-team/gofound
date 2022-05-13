<template>
  <el-progress type="dashboard" :percentage="percent" :colors="colors"></el-progress>
</template>

<script>
export default {
  name: 'ProgressChat',
  props: {
    percent: {
      type: Number,
      default: 0,
    },
  },
  data(){
    return{
      colors: [
        { color: '#45c432', percentage: 30 },
        { color: '#ffd700', percentage: 50 },
        { color: '#e72048', percentage: 80 },
      ]
    }
  }
}
</script>

<style scoped>

</style>
