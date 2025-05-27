<template>
  <NavBar :activeIndex="$route.path" />
  <div class="friends-container">
    <div class="header">
      <h2>好友列表</h2>
      <el-button type="primary" @click="showAddFriendDialog">添加好友</el-button>
    </div>

    <div class="friends-grid">
      <FriendCard 
        v-for="friend in friends" 
        :key="friend.id"
        :friend="friend"
        @send-message="handleSendMessage"
        @delete-friend="handleDeleteFriend"
        @click="goToFriendDetail(friend.id)"
      />
    </div>

    <!-- 添加好友对话框 -->
    <el-dialog v-model="addFriendDialogVisible" title="添加好友" width="50%" @close="handleSearchClear">
      <el-input 
        v-model="searchEmail" 
        placeholder="请输入好友邮箱" 
        clearable
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button :icon="Search" @click="handleSearch" />
        </template>
      </el-input>

      <div class="search-results" v-if="!searchNull && isSearch">
        <div class="search-card">
          <el-avatar :size="50" :src="picBaseUrl + searchResult.pic" />
          <div class="user-info">
              <span class="username">{{ searchResult.name }}</span>
          </div>
          <div>
            <el-button type="primary" size="small" @click="handleAddFriend(searchResult.id)" v-if="searchResult.isFriend === 2">
            添加好友
            </el-button>
            <el-button type="success" size="small" v-if="searchResult.isFriend === 1">
            已是好友
            </el-button>
            <el-button type="warning" size="small" v-if="searchResult.isFriend === 0">
            待同意
            </el-button>
          </div>
        </div>
      </div>
      <div class="search-null" v-if="searchNull">该邮箱未注册！</div>
    </el-dialog>
  </div>


  <!-- 给好友捎句话窗口 -->
  <el-dialog v-model="showAddMessage" title="给好友捎句话" width="500px">
    <el-input
      v-model="messageContent"
      type="textarea"
      :rows="8"
      placeholder="请输入要捎带的话"
      maxlength="500"
      show-word-limit
      style="margin: 20px 0"
    />
    <template #footer>
      <el-button @click="showAddMessage = false" size="large">取消</el-button>
      <el-button type="primary" @click="submitAddMessage" size="large">完成</el-button>
    </template>
  </el-dialog>

  <!-- 确认窗口 -->
  <ConfirmDialog
    ref="confirmDialog"
    v-model="dialogVisible"
    :title="dialogTitle"
  />
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import NavBar from '@/components/NavBar.vue'
import FriendCard from '@/components/FriendCard.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { ElMessage, ElStep } from 'element-plus'
import { useUserStore } from '@/stores/user.js'
import { getFriendList, searchUserInfoByEmail, sendAddFriendRequest,deleteFriend } from '@/api/friend'
import { addMessage } from '@/api/message'

const router = useRouter()
const userStore = useUserStore()
// 头像公共前缀
const picBaseUrl = import.meta.env.VITE_PIC_BASE_URL

// 好友列表数据
const friends = ref([])

// 添加好友对话框相关
const addFriendDialogVisible = ref(false)
const searchEmail = ref('')
const searchResult = ref(null)
const searchNull  = ref(false)
const isSearch = ref(false)

// 确认窗口数据
const confirmDialog = ref({})
const dialogVisible = ref(false)
const dialogTitle = ref('')

// 给好友捎句话窗口数据
const showAddMessage = ref(false)
const messageContent = ref('')
const messageToFriendId = ref('')

const showAddFriendDialog = () => {
  addFriendDialogVisible.value = true
  searchEmail.value = ''
  searchResult.value = []
}

onMounted(() => {
  fetchFriends()
})

// 获取好友列表数据
const fetchFriends = async () => {
  const res = await getFriendList(userStore.userInfo.id)
  if(res.data.status === 1){
    friends.value = res.data.data
  }else{
    ElMessage.error(res.data.message)
  }
}

// 搜索好友
const handleSearch = async() => {
  // 输入框为空
  if (!searchEmail.value.trim()) return
  // 检查邮箱合法性
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.com$/
  if (!emailRegex.test(searchEmail.value)) {
    ElMessage.error('请输入正确的邮箱格式')
    return
  }
  // 用户不能搜索自己
  if (searchEmail.value === userStore.userInfo.email) {
    ElMessage.error('请不要搜索自己')
    return
  }
  
  // 调用后端API搜索用户
  const res = await searchUserInfoByEmail(searchEmail.value)
  isSearch.value = true
  if(res.data.status === 1){
    searchNull.value = (res.data.data === null)
    console.log(searchNull.value)
    searchResult.value = res.data.data
  }else{
    ElMessage.error(res.data.message)
  }
}
// 关闭搜索窗口时，清空数据
const handleSearchClear = () => {
  searchNull.value = false
  searchResult.value = []
  isSearch.value = false
}

const handleAddFriend = async (friendId) => {
  // 调用后端API添加好友
  const data = {
    user_id: userStore.userInfo.id,
    friend_id: friendId
  }
  const res = await sendAddFriendRequest(data)
  if(res.data.status === 1){
    ElMessage.success('发送好友请求成功！')
  }else{
    ElMessage.error(res.data.message)
  }
  // 关闭对话框
  addFriendDialogVisible.value = false
}

const handleSendMessage = (friendId) => {
  // 显示发送消息的对话框
  showAddMessage.value  = true
  messageToFriendId.value = friendId
}

const submitAddMessage = async () => { 
  if (messageContent.value === ''){
    ElMessage.error('发送内容不能为空！')
    return
  }
  const title = `好友${userStore.userInfo.name}给你发来消息`
  const data = {
    title: title,
    description: messageContent.value,
    from_id: userStore.userInfo.id,
    to_id: messageToFriendId.value,
    type: 0
  }

  const res = await addMessage(data)
  if(res.data.status === 1){
    ElMessage.success('发送成功！')
    messageContent.value = ''
    showAddMessage.value = false
    messageToFriendId.value = ''
  }else {
    ElMessage.error('发送失败！')
    console.log(res.data.message)
  }
}

const handleDeleteFriend = async (friendId) => {
  // 弹出确认窗口，询问用户是否要删除该好友
  dialogTitle.value = '确认要删除该好友吗？'
  dialogVisible.value = true
  const isConfirmed = await confirmDialog.value.confirm()

  if (isConfirmed) {
    // 调用后端API删除好友
    const res = await deleteFriend(userStore.userInfo.id,friendId)
    if (res.data.status === 1 ){
      ElMessage.success('删除好友成功')
      friends.value = friends.value.filter(friend => friend.id !== friendId)
    }else{
      ElMessage.error('删除好友失败')
      console.log(res.data.message)
    }
  }

}

const goToFriendDetail = (friendId) => {
  router.push(`/friendDetail/${friendId}`)
}
</script>

<style scoped>
.friends-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.friends-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); 
  gap: 20px;
}

.search-results {
  margin-top: 20px;
}

.search-card {
  display: flex;
  align-items: center;
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  margin-bottom: 10px;
}

.search-card .el-avatar {
  margin-right: 15px;
}

.search-card .user-info {
  flex: 1;
}

.search-card .username {
  font-weight: bold;
}
.search-null{
    text-align: center;
  padding: 20px;
  color: #999;
  font-size: 14px;
}
</style>