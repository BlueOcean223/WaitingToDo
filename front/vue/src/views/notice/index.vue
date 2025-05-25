<template>
    <Navbar :active-index="$route.path" />
  <div class="message-container">
    <div class="message-header">
      <h2>消息通知</h2>
      <el-button type="primary" @click="markAllAsRead" :disabled="unreadCount === 0">
        一键已读
      </el-button>
    </div>
    
    <div class="message-list">
      <div class="loading" v-if="loading">加载中...</div>
      <div class="no-more" v-if="!hasMore">没有更多消息了</div>
      <el-card
        v-for="message in messages"
        :key="message.id"
        class="message-card"
        :class="{ 'unread': !message.is_read }"
      >
        <template #header>
          <div class="card-header">
            <span class="title">{{ message.title }}</span>
            <span class="time">{{ formatTime(message.send_time) }}</span>
          </div>
        </template>
        
        <div class="card-content">
          <p>{{ message.description }}</p>
        </div>
        
        <div class="card-actions">
          <!-- 好友请求或小队邀请的按钮 -->
          <template v-if="(message.type === 1 || message.type === 2) && !message.is_read">
            <el-button type="success" size="small" @click="handleAccept(message)">
              接受
            </el-button>
            <el-button type="danger" size="small" @click="handleReject(message)">
              拒绝
            </el-button>
          </template>
          
          <!-- 普通消息或已处理的消息的按钮 -->
          <template v-else>
            <el-button
              type="primary"
              size="small"
              @click="markAsRead(message)"
              :disabled="message.is_read"
            >
              已读
            </el-button>
            <el-button type="danger" size="small" @click="deleteMessage(message)">
              删除
            </el-button>
          </template>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import Navbar from '@/components/Navbar.vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getMessageList } from '@/api/message'

// 当前用户信息
const userStore = useUserStore()

// 分页相关数据
const loading = ref(false)
const currentPage = ref(1)
const hasMore = ref(true)
const pageSize = ref(5)

// 模拟消息数据
const messages = ref([])

// 滚动事件监听
const handleScroll = () => {
  const { scrollTop, scrollHeight, clientHeight } = document.documentElement
  if (scrollTop + clientHeight >= scrollHeight - 100 && !loading && hasMore) {
    fetchMessages()
  }
}

// 获取消息数据
const fetchMessages = async () => {
  try{
    if(!hasMore.value) return

    loading.value = true
    // 调用API接口
    const res = await getMessageList(currentPage.value,pageSize.value,userStore.userInfo.id)
    if(res.data.status === 1){
      // 查询成功
      if(res.data.data == null){
        hasMore.value = false
        return
      }
      messages.value.push(...res.data.data)
      currentPage.value++
    }else{
      // 查询失败
      ElMessage.error(res.data.message)
    }
  }finally{
    loading.value = false
  }
}

// 计算未读消息数量
const unreadCount = computed(() => {
  return messages.value.filter(msg => !msg.is_read).length
})

// 格式化时间显示
const formatTime = (timeString) => {
  // 将.后面的全部内容忽略
  const index = timeString.indexOf('.')
  return timeString.substring(0, index)
}

// 标记单条消息为已读
const markAsRead = (message) => {
  message.is_read = 1
  // TODO: 调用API更新消息状态
  emit('update-unread', unreadCount.value)
  ElMessage.success('消息已标记为已读')
}

// 一键已读所有消息
const markAllAsRead = () => {
  messages.value.forEach(msg => {
    if (!msg.is_read) {
      msg.is_read = 1
    }
  })
  // TODO: 调用API批量更新消息状态
  emit('update-unread', unreadCount.value)
  ElMessage.success('所有消息已标记为已读')
}

// 删除消息
const deleteMessage = (message) => {
  const index = messages.value.findIndex(msg => msg.id === message.id)
  if (index !== -1) {
    const wasUnread = !messages.value[index].is_read
    messages.value.splice(index, 1)
    // TODO: 调用API删除消息
    if (wasUnread) {
      emit('update-unread', unreadCount.value)
    }
    ElMessage.success('消息已删除')
  }
}

// 处理接受请求
const handleAccept = (message) => {
  // 根据消息类型执行不同操作
  if (message.type === 'friend_request') {
    ElMessage.success('好友请求已接受')
  } else if (message.type === 'team_invite') {
    ElMessage.success('小队邀请已接受')
  }
  
  // 标记为已读
  message.is_read = 1
  // TODO: 调用API处理接受请求
  emit('update-unread', unreadCount.value)
}

// 处理拒绝请求
const handleReject = (message) => {
  // 根据消息类型执行不同操作
  if (message.type === 'friend_request') {
    ElMessage.warning('好友请求已拒绝')
  } else if (message.type === 'team_invite') {
    ElMessage.warning('小队邀请已拒绝')
  }
  
  // 标记为已读
  message.is_read = 1
  // TODO: 调用API处理拒绝请求
  emit('update-unread', unreadCount.value)
}

// 组件挂载时获取消息数据
onMounted(() => {
  // 监听滚动事件
  window.addEventListener('scroll', handleScroll)
  // 获取消息数据
  fetchMessages()
})

onBeforeUnmount(() => {
  // 组件销毁时取消滚动事件监听
  window.removeEventListener('scroll', handleScroll)
})

// 定义事件
const emit = defineEmits(['update-unread'])
</script>

<style scoped>
.message-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.message-card {
  transition: all 0.3s ease;
}

.message-card.unread {
  border-left: 4px solid #409eff;
  background-color: #f5f9ff;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header .title {
  font-weight: bold;
}

.card-header .time {
  font-size: 0.8em;
  color: #999;
}

.card-content {
  margin-bottom: 15px;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
.loading {
  text-align: center;
  padding: 20px;
  color: #666;
}
.no-more {
  text-align: center;
  padding: 20px;
  color: #999;
  font-size: 14px;
}
</style>