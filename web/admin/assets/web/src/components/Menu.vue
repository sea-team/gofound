<template>
  <div style="text-align:center;width: 200px;color:white" v-if="!isCollapsed">
    <svg t="1652075159972" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"
         p-id="6654" width="50" height="50"
    >
      <path
          d="M501.76 402.432c-39.936 9.216-67.584 16.384-108.544 28.672-9.216 1.024-9.216 4.096-19.456-6.144-9.216-12.288-16.384-16.384-30.72-24.576-39.936-19.456-79.872-13.312-116.736 9.216-44.032 28.672-65.536 68.608-65.536 120.832 0 50.176 36.864 93.184 87.04 101.376 44.032 6.144 79.872-9.216 108.544-43.008 6.144-7.168 12.288-15.36 16.384-22.528H248.832c-13.312 0-16.384-7.168-12.288-19.456 7.168-19.456 22.528-52.224 31.744-68.608 1.024-4.096 6.144-9.216 15.36-9.216h234.496c1.024 18.432 1.024 34.816-1.024 52.224-7.168 46.08-24.576 89.088-52.224 125.952-46.08 59.392-105.472 98.304-182.272 108.544-61.44 7.168-120.832-4.096-171.008-43.008-50.176-37.888-73.728-83.968-81.92-143.36-7.168-68.608 12.288-133.12 53.248-186.368 46.08-59.392 105.472-96.256 179.2-110.592 59.392-12.288 117.76-4.096 169.984 30.72 34.816 22.528 58.368 52.224 73.728 90.112 6.144 5.12 3.072 8.192-4.096 9.216z m211.968 354.304c-58.368-1.024-111.616-16.384-156.672-56.32-37.888-31.744-61.44-73.728-68.608-123.904-12.288-71.68 7.168-138.24 52.224-192.512 46.08-61.44 102.4-93.184 179.2-105.472 65.536-12.288 126.976-6.144 182.272 31.744 50.176 34.816 80.896 80.896 90.112 142.336 12.288 87.04-13.312 156.672-73.728 216.064-43.008 43.008-95.232 68.608-153.6 80.896-18.432 5.12-35.84 5.12-51.2 7.168z m150.528-259.072c0-7.168 0-15.36-1.024-21.504-12.288-64.512-68.608-98.304-132.096-83.968-59.392 13.312-98.304 50.176-111.616 111.616-12.288 50.176 13.312 101.376 58.368 120.832 34.816 15.36 68.608 13.312 104.448-4.096 53.248-26.624 80.896-66.56 81.92-122.88z"
          fill="#2DBCAF" p-id="6655"
      ></path>
    </svg>
    <h1 style="margin-top:0px;">GoFound Admin</h1>
  </div>
  <el-menu
      :default-active="active+''"
      :collapse="isCollapsed"
      class="el-menu-vertical-demo"
      background-color="#191a22"
      text-color="#fff"
      active-text-color="#ffd04b"
  >

    <router-link v-for="(item,index) in menus" :to="{name:item.name}" :key="item.name">
      <el-menu-item :index="index">
        <Icon :name="item.icon" :color="item.color"/>
        <span v-text="item.label"></span>
      </el-menu-item>
    </router-link>

  </el-menu>
</template>

<script>
import { Document, Service } from '@element-plus/icons-vue'
import menus from '../menus'
import Icon from './Icon.vue'

export default {
  name: 'Menu',
  components: { Icon, Service, Document },
  props: {
    isCollapsed: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      menus: menus,
    }
  },
  computed: {
    active() {
      let active = 0
      this.menus.forEach((item, index) => {
        if (item.name === this.$route.name) {
          active = index
        }
      })
      return active
    },
  },
  methods: {
    openDocument() {
      window.open('https://gitee.com/tompeppa/gofound')
    },
  },
}
</script>

<style scoped>
.el-menu a{
  text-decoration: none;
}
</style>
