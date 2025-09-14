<template>
    <Navbar :active-index="$route.path" />
  <div class="message-container">
    <div class="message-header">
      <h2>消息通知</h2>
      <div class="header-buttons">
        <el-switch
          v-model="showFriendInfo"
          active-text="显示好友信息"
          inactive-text="关闭"
          class="friend-info-switch"
        />
        <el-button type="primary" @click="markAllAsRead" :disabled="unreadCount === 0">
          一键已读
        </el-button>
      </div>
    </div>
    
    <div class="message-list">
      <el-card
        v-for="message in messages"
        :key="message.id"
        @click="showSendMessageUser(message.from_id)"
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
            <el-button type="success" size="small" @click.stop="handleAccept(message)">
              接受
            </el-button>
            <el-button type="danger" size="small" @click.stop="handleReject(message)">
              拒绝
            </el-button>
          </template>
          
          <!-- 普通消息或已处理的消息的按钮 -->
          <template v-else>
            <el-button
              type="primary"
              size="small"
              @click.stop="markAsRead(message)"
              :disabled="message.is_read"
            >
              已读
            </el-button>
            <el-button type="danger" size="small" @click.stop="handleDelete(message)">
              删除
            </el-button>
          </template>
        </div>
      </el-card>
      <div class="loading" v-if="loading">加载中...</div>
      <div class="no-more" v-if="!hasMore">没有更多消息了</div>
    </div>
  </div>

  <!-- 确认窗口 -->
  <ConfirmDialog
    ref="confirmDialog"
    v-model="dialogVisible"
    :title="dialogTitle"
  />

  <!-- 展示邮件发送者信息窗口 -->
  <el-dialog 
    v-model="isShowFromUser" 
    title="是这位用户给你发送的消息" 
    width="500px"
    :close-on-click-modal="false"
    v-if="showFriendInfo"
  >
    <div class="sender-info-container">
      <div class="sender-avatar">
        <el-avatar :size="120" :src="sendMessageUser.avatar">
          <img src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png" />
        </el-avatar>
      </div>
      <div class="sender-details">
        <div class="detail-item">
          <span class="detail-label">昵称：</span>
          <span class="detail-value">{{ sendMessageUser.name || '未知用户' }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">邮箱：</span>
          <span class="detail-value">{{ sendMessageUser.email || '未提供邮箱' }}</span>
        </div>
        <div class="detail-item" v-if="sendMessageUser.description">
          <span class="detail-label">简介：</span>
          <span class="detail-value">{{ sendMessageUser.description }}</span>
        </div>
      </div>
    </div>
    
    <template #footer>
      <el-button type="primary" @click="isShowFromUser = false" size="large">关闭</el-button>
    </template>
  </el-dialog>

</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import Navbar from '@/components/Navbar.vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { useMessageStore } from '@/stores/message'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { getMessageList, updateMessage,deleteMessage,readAllMessage, handleRequest } from '@/api/message'
import { getUserById } from '@/api/user'

// 当前用户信息
const userStore = useUserStore()
// 未读取消息全局管理
const messageStore = useMessageStore()

// 分页相关数据
const loading = ref(false)
const currentPage = ref(1)
const hasMore = ref(true)
const pageSize = ref(5)

// 消息数据
const messages = ref([])

// 确认窗口相关数据
const confirmDialog = ref({})
const dialogVisible = ref(false)
const dialogTitle = ref('')

// 是否显示好友信息
const showFriendInfo = ref(true)

// 滚动事件监听
const handleScroll = () => {
  const { scrollTop, scrollHeight, clientHeight } = document.documentElement
  if (scrollTop + clientHeight >= scrollHeight - 100 && !loading.value && hasMore.value) {
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
  if (index !== -1) {
    timeString = timeString.substring(0, index)
  }
  return timeString
}

// 标记单条消息为已读
const markAsRead = async (message) => {
  message.is_read = 1
  // 调用API更新消息状态
  const res = await updateMessage(message)
  if(res.data.status === 1){
    ElMessage.success('消息已标记为已读')
    // 更新未读消息数量
    messageStore.readMessage()
  }else{
    ElMessage.error(res.data.message)
  }
}

// 一键已读所有消息
const markAllAsRead = async() => {
  // 如果用户还有未处理的好友请求或组队邀请，则不能全部标记为已读
  var ok = false
  messages.value.forEach(msg => {
    if((msg.type === 1 || msg.type === 2) && msg.is_read === 0 ){
      ok = true
      return
    }
  })
  if (ok) {
    ElMessage.warning('请先处理好友请求或组队邀请')
    return
  }
  // 将所有消息标记为已读
  messages.value.forEach(msg => {
    if (!msg.is_read) {
      msg.is_read = 1
    }
  })
  // 调用API批量更新消息状态
  const res = await readAllMessage(userStore.userInfo.id)
  if (res.data.status === 1) {
    ElMessage.success('所有消息已标记为已读')
    // 更新全局的未读消息数量
    messageStore.readAllMessage()
  }else{
    ElMessage.error(res.data.message)
  }
}

// 删除消息
const handleDelete = async (message) => {
  dialogTitle.value = '确认要删除该消息吗？'
  dialogVisible.value = true
  const isConfirmed = await confirmDialog.value.confirm()
  if (isConfirmed){
    // 调用API删除消息
    const res = await deleteMessage(message.id)
    if(res.data.status === 1){
      ElMessage.success('消息已删除')
      // 从消息队列中删除该消息
      messages.value = messages.value.filter(item => item.id !== message.id)
      // 如果该消息为未读消息，则更新未读消息数
      if(message.isRead === 0){
        messageStore.readMessage()
      }
    }else{
      ElMessage.error(res.data.message)
    }
  }
  console.log(isConfirmed)
  dialogVisible.value = false
}

// 处理接受请求
const handleAccept = async (message) => {
  // 调用API处理接受请求
  const data = {
    ...message,
    request_action: 1
  }
  const res = await handleRequest(data)

  if(res.data.status === 1){
    // 根据消息类型执行不同操作
    if (message.type === 1) {
      ElMessage.success('好友请求已接受')
    } else if (message.type === 2) {
      ElMessage.success('小队邀请已接受')
    }
  }else{
    ElMessage.error(res.data.message)
  }
  // 标记为已读
  message.is_read = 1
  messageStore.readMessage()
}

// 处理拒绝请求
const handleReject = async (message) => {
  // 调用API处理接受请求
  const data = {
    ...message,
    request_action: 0
  }
  const res = await handleRequest(data)

  if(res.data.status === 1){
    // 根据消息类型执行不同操作
    if (message.type === 1) {
      ElMessage.warning('好友请求已拒绝')
    } else if (message.type === 2) {
      ElMessage.warning('小队邀请已拒绝')
    }
    // 标记为已读
    message.is_read = 1
    messageStore.readMessage()
  }else{
    console.log(res.data.message)
    ElMessage.error('处理请求失败')
  }
}

// 是否展示
const isShowFromUser = ref(false)
// 信息存储
const sendMessageUser = ref({})
// 点击时展示发送该消息的用户信息
const showSendMessageUser = async (from_id) => {
  if (!showFriendInfo.value) return
  
  isShowFromUser.value = true

  const res = await getUserById(from_id)
  if(res.data.status === 1){
    sendMessageUser.value = res.data.data
    sendMessageUser.value.avatar = import.meta.env.VITE_PIC_BASE_URL + res.data.data.pic
  }else{
    ElMessage.error('展示失败')
    console.log(res.data.message)
  }
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

.header-buttons {
  display: flex;
  align-items: center;
}

.friend-info-switch {
  margin-right: 15px; 
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
.sender-info-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.sender-avatar {
  margin-bottom: 24px;
}

.sender-avatar .el-avatar {
  border: 3px solid #f0f0f0;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
}

.sender-avatar .el-avatar:hover {
  transform: scale(1.05);
}

.sender-details {
  width: 100%;
  padding: 0 20px;
}

.detail-item {
  display: flex;
  margin-bottom: 16px;
  line-height: 1.6;
}

.detail-label {
  font-weight: 500;
  color: #606266;
  min-width: 60px;
  text-align: right;
  margin-right: 12px;
}

.detail-value {
  flex: 1;
  color: #303133;
  word-break: break-word;
}
</style>