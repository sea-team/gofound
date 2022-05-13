<script>
import { h } from 'vue'
import { ElIcon } from 'element-plus'
export default {
  name: 'Icon',
  props: {
    name: {
      type: String,
      required: true
    },
    color:{
      type:String,
      required: true
    }
  },
  render() {
    return h(ElIcon, {},[h(this.name,{color:this.color})])
  },
}
</script>

<style scoped>

</style>
