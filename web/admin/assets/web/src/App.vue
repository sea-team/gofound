<script>
import Menu from './components/Menu.vue'
import { Expand, Fold } from '@element-plus/icons-vue'
import GC from './components/GC.vue'
import { useDark, useToggle } from '@vueuse/core'
const _isDark = useDark()
const toggleDark = useToggle(_isDark)
import router from './router'
export default {
  components: {
    GC,
    Fold,
    Expand,
    Menu,
  },
  data() {
    return {
      isCollapsed: false,
      isDark: _isDark,
    }
  },
  computed: {
    name() {
      let name = ''
      router.options.routes.forEach(route => {
        if (route.name === this.$route.name) {
          name = route.label
        }
      })
      return name
    },
  },
  created() {

  },
  methods: {
    toggle() {

      toggleDark()
    },
  }
}
</script>

<template>
  <el-container>
    <el-aside class="menubar">
      <Menu :isCollapsed="isCollapsed"></Menu>
    </el-aside>
    <el-container>

      <el-header
          style="display:flex;justify-content: space-between;border-bottom: 1px solid var(--el-card-border-color);"
      >
        <div style="display:flex;align-items: center">
          <el-link @click="isCollapsed=!isCollapsed">
            <el-icon :size="26">
              <fold v-if="!isCollapsed"/>
              <expand v-else/>
            </el-icon>

          </el-link>
          <span style="margin-left:10px;" v-text="name"></span>
        </div>
        <div style="display:flex;align-items:center">
          主题切换：<el-switch v-model="isDark">切换主题</el-switch>
          <span style="margin-right:10px;"></span>
          <GC/>
        </div>
      </el-header>
      <el-main class="main">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<style>

</style>
