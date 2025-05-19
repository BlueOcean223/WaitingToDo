<template>
  <el-menu
    :default-active="activeIndex"
    class="el-menu-demo"
    mode="horizontal"
    :ellipsis="false"
    :router="true"
  >
    <!-- 左侧菜单项 -->
    <el-menu-item index="/home">主页</el-menu-item>
    <el-menu-item index="/team">小组</el-menu-item>
    <el-menu-item index="/friend">好友</el-menu-item>
    
    <!-- 右侧用户信息区域 -->
    <div class="flex-grow"></div>
    
    <el-menu-item index="/notice" class="notification-icon">
      <el-icon><Bell /></el-icon>
      通知
    </el-menu-item>
    
    <el-sub-menu index="user-menu">
      <template #title>
        <span class="username">{{ userName }}</span>
        <el-avatar :src="avatarUrl" size="default"></el-avatar>
      </template>
      <el-menu-item index="/profile">个人资料</el-menu-item>
      <el-menu-item index="/logout" @click="logout()">退出登录</el-menu-item>
    </el-sub-menu>
  </el-menu>
</template>

<script>
import { Bell } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

export default {
  name: 'NavBar',
  components: {
    Bell
  },
  props: {
    activeIndex: {
      type: String,
      default: '/home'
    }
  },
  computed: {
    user() {
      return useUserStore().userInfo // 从全局状态中获取用户信息
    },
    userName() {
      return this.user.name || '未登录'
    },
    avatarUrl() {
      const baseUrl = import.meta.env.VITE_PIC_BASE_URL || 'http://192.168.163.129:9000'
      return `${baseUrl}${this.user.pic || ''}`
    }
  },
  methods: {
    logout() {
        //  清除用户信息
        localStorage.removeItem('user')
        //  清除token
        localStorage.removeItem('token')
        this.$router.push('/login')
    },
  }
}
</script>

<style scoped>
.el-menu-demo {
  display: flex;
  padding: 0 20px;
}

.flex-grow {
  flex-grow: 1;
}

.username {
  margin-right: 10px;
  font-size: 14px;
}

.notification-icon {
  margin-right: 10px;
}

.el-menu--horizontal > .el-menu-item,
.el-menu--horizontal > .el-sub-menu {
  height: 60px;
  line-height: 60px;
}

.el-menu--horizontal .el-sub-menu__title {
  height: 60px;
  line-height: 60px;
}
</style>