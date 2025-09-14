<template>
  <el-menu
    :default-active="activeIndex"
    class="el-menu-demo"
    mode="horizontal"
    :ellipsis="false"
    :router="true"
  >
    <!-- 左侧菜单项 -->
    <el-menu-item index="/">主页</el-menu-item>
    <el-menu-item index="/team">小组</el-menu-item>
    <el-menu-item index="/friend">好友</el-menu-item>
    
    <!-- 右侧用户信息区域 -->
    <div class="flex-grow"></div>
    
    <el-menu-item index="/notice" class="notification-icon">
      <el-badge :value="unreadCount" :max="99" :hidden="unreadCount === 0">
        <el-icon><Bell /></el-icon>
        <span>通知</span>
      </el-badge>
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

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bell } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { useMessageStore } from '@/stores/message'
import { getUnreadMessageCount } from '@/api/message'

// 定义props
const props = defineProps({
  activeIndex: {
    type: String,
    default: '/'
  }
})

// 获取路由实例
const router = useRouter()

// 获取store实例
const userStore = useUserStore()
const messageStore = useMessageStore()

// 计算属性
const user = computed(() => userStore.userInfo)
const userName = computed(() => user.value.name || '未登录')
const avatarUrl = computed(() => {
  const baseUrl = import.meta.env.VITE_PIC_BASE_URL
  return `${baseUrl}${user.value.pic || ''}`
})
const unreadCount = computed(() => messageStore.unreadCount)

// 方法
const logout = () => {
  //  清除用户信息
  localStorage.removeItem('user')
  //  清除token
  localStorage.removeItem('token')
  // 清除全局用户信息
  userStore.clearUserInfo()
  // 跳转到登录页
  router.push('/login')
}

const getUnreadCount = async () => {
  const res = await getUnreadMessageCount(user.value.id)
  if(res.data.status === 1){
    messageStore.setUnreadCount(res.data.data)
  }else{
    console.log(res.data.message)
  }
}

// 生命周期钩子
onMounted(() => {
  getUnreadCount()
})
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
  position: relative;
}

.el-menu--horizontal {
  height: 45px;
  line-height: 45px;
}

.el-menu--horizontal > .el-menu-item,
.el-menu--horizontal > .el-sub-menu {
  height: 45px;
  line-height: 45px;
}

.el-menu--horizontal .el-sub-menu__title {
  height: 45px;
  line-height: 45px;
}
</style>