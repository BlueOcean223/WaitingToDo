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
          <el-button type="primary" size="small" @click="handleAddFriend(searchResult.id)">
            添加好友
          </el-button>
        </div>
      </div>
      <div class="search-null" v-if="searchNull">该邮箱未注册！</div>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import NavBar from '@/components/NavBar.vue'
import FriendCard from '@/components/FriendCard.vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user.js'
import { getFriendList, searchUserInfoByEmail } from '@/api/friend'

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

const handleAddFriend = (userId) => {
  // TODO: 调用后端API添加好友
  console.log('添加好友:', userId)
  addFriendDialogVisible.value = false
  // 可以在这里刷新好友列表
}

const handleSendMessage = (friendId) => {
  // TODO: 处理发送消息逻辑
  console.log('发送消息给:', friendId)
}

const handleDeleteFriend = (friendId) => {
  // TODO: 调用后端API删除好友
  console.log('删除好友:', friendId)
  // 可以在这里刷新好友列表
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