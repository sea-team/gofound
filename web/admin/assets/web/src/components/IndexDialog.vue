<template>
  <el-dialog draggable v-model="visible" @close="$emit('close')" title="添加索引">
    <el-alert type="success" style="margin-bottom: 10px;">请保持ID的唯一，如果存在将会更新数据。</el-alert>
    <el-form ref="form" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="数据库">
        <el-tag v-text="db"></el-tag>
      </el-form-item>
      <el-form-item label="ID" prop="id">
        <el-input v-model="form.id" :disabled="data!=null" placeholder="索引id"></el-input>
      </el-form-item>
      <el-form-item label="索引文本" prop="text">
        <el-input type="textarea" v-model="form.text" placeholder="索引文本" :rows="5"></el-input>
      </el-form-item>

      <el-form-item label="JSON文档" prop="document">
        <el-input type="textarea" v-model="form.document" placeholder="JSON文本" :rows="5"></el-input>
        <el-link type="success" @click="example()">填入示例</el-link>
      </el-form-item>

    </el-form>
    <template #footer>
      <el-button type="primary" @click="save()">确定</el-button>
    </template>
  </el-dialog>
</template>

<script>
import api from '../api'

export default {
  name: 'IndexDialog',
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    db: {
      type: String,
      default: '',
    },
    data: {
      type: Object,
      default: () => null,
    },
  },
  data() {
    return {
      rules: {
        id: [
          { required: true, message: '请输入ID', trigger: 'blur' },
        ],
        text: [
          { required: true, message: '请输入索引文本', trigger: 'blur' },
        ],
        document: [
          { required: true, message: '请输入JSON文档', trigger: 'blur' },
        ],
      },
      form: {
        id: '',
        text: '',
        document: '',
      },
    }
  },
  watch: {
    data(val) {
      if (val) {
        this.form = val
        if(val.originalText) {
          this.form.text = val.originalText
        }

        this.form.document = JSON.stringify(val.document)
      } else {
        this.form = {
          id: '',
          text: '',
          document: '',
        }
      }
    },
  },
  methods: {
    example() {
      this.form.document = JSON.stringify({ name: '张三', age: 18 })
    },
    save() {
      this.$refs.form.validate(valid => {
        if (valid) {
          //校验json文档
          let data = {
            id: parseInt(this.form.id),
            text: this.form.text,
            document: this.form.document,
          }
          try {
            data.document = JSON.parse(this.form.document)
          } catch (e) {
            this.$message.error('JSON文档格式错误')
            return
          }
          api.addIndex(this.db, data).then(({ data }) => {
            console.log(data)
            if (data.state) {
              this.$message.success('添加成功!')
              this.$emit('success')
            } else {
              this.$message.error(data.message)
            }
          })
        }
      })
    },
  },
}
</script>

<style scoped>

</style>
