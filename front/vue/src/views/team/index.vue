<template>
  <Navbar :activeIndex="$route.path" />
  <div class="team-page">
    <div class="header">
      <div></div>
      <el-button type="primary" @click="showJoinGroup = true">加入小组</el-button>
    </div>

    <div class="task-list">
      <TeamTask 
        v-for="task in tasks" 
        :key="task.id" 
        :task="task"
        @invite="handleInviteMember"
        @complete="handleCompleteTask"
        @exitTeam="handleExitTeam"
      />
      <div class="loading" v-if="loading">加载中...</div>
      <div class="no-more" v-if="!hasMore">没有更多任务了</div>
    </div>
    
    <!-- 发布按钮 -->
    <el-button 
      type="primary" 
      circle 
      class="publish-btn"
      @click="showAddTask = true"
    >
      <el-icon><Plus /></el-icon>
    </el-button>
    <!-- 添加任务模态框 -->
    <el-dialog v-model="showAddTask" title="添加任务" width="500px">
      <el-form 
        :model="addTaskForm" 
        label-width="80px" 
        label-position="left"
        :rules="addTaskRules"
        ref="addTaskFormRef"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="addTaskForm.title" placeholder="请输入标题"/>
        </el-form-item>
        <el-form-item label="内容" prop="description">
          <el-input 
            v-model="addTaskForm.description" 
            type="textarea" 
            :rows="8"
            placeholder="请输入任务内容"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="截止时间" prop="ddl">
          <el-date-picker 
            v-model="addTaskForm.ddl" 
            type="datetime" 
            placeholder="请选择截止时间" 
            size="large"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddTask = false" size="large">取消</el-button>
        <el-button type="primary" @click="submitAddTask" size="large">完成</el-button>
      </template>
    </el-dialog>

    <!-- 邀请成员对话框 -->
    <el-dialog v-model="inviteVisible" title="邀请好友加入小组任务" width="35%" @close="handleInviteClear">
      <!-- 邀请码部分 -->
      <div class="invite-code-section">
        <div class="invite-code-display">
          <span>小组邀请码：</span>
          <span class="invite-code">{{ currentInviteCode }}</span>
          <el-button type="primary" size="small" @click="copyInviteCode">
            复制
          </el-button>
        </div>
      </div>
      
      <el-divider>或者直接邀请好友</el-divider>
      
      <div class="invite-results">
        <div class="invite-card" v-for="friend in friends" :key="friend.id">
          <el-avatar :size="50" :src="picBaseUrl + friend.pic" />
          <div class="user-info">
              <span class="username">{{ friend.name }}</span>
          </div>
          <div>
            <el-button type="primary" size="small" @click="handleInviteFriend(friend.id)" v-if="isInvited(friend.id)">
            邀请好友
            </el-button>
            <el-button type="success" size="small" v-else>
            已在小组中
            </el-button>
          </div>
        </div>
      </div>
      <div class="friends-null" v-if="friendsNull">您暂时没有好友可以邀请！</div>
    </el-dialog>

    <!-- 加入小组对话框 -->
    <el-dialog v-model="showJoinGroup" title="加入小组" width="400px">
      <el-form :model="joinGroupForm" ref="joinGroupFormRef" :rules="joinGroupRules">
        <el-form-item 
          label="邀请码" 
          prop="inviteCode"
        >
          <el-input 
            v-model="joinGroupForm.inviteCode" 
            placeholder="请输入邀请码"
            maxlength="6"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showJoinGroup = false">取消</el-button>
          <el-button type="primary" @click="handleJoinGroup">加入</el-button>
        </span>
      </template>
    </el-dialog>

  </div>

  <!-- 确认操作弹窗 -->
  <ConfirmDialog
    ref="confirmDialog"
    v-model="dialogVisible"
    :title="dialogTitle"
  />
</template>

<script setup name="TeamPage">
import { ref,onMounted,onBeforeUnmount } from 'vue'
import TeamTask from '@/components/TeamTask.vue'
import Navbar from '@/components/Navbar.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getFriendList } from '@/api/friend'
import { getTeamTaskList, removeTeamTask, addTeamTask, completeTeamTask, inviteMember,getTeamTaskInviteCode,joinTeamByCode } from '@/api/task'

// 分页相关数据
const loading = ref(false)
const currentPage = ref(1)
const hasMore = ref(true)
const pageSize = ref(5) // 每页显示的任务数量

// 当前用户全局信息
const userStore = useUserStore()

// 图片基础URL
const picBaseUrl = import.meta.env.VITE_PIC_BASE_URL

// 确认窗口信息
const confirmDialog = ref(null)
const dialogVisible = ref(false)
const dialogTitle = ref('')

// 任务列表
const tasks = ref([])

const showAddTask = ref(false)
const addTaskFormRef = ref(null)

const addTaskForm = ref({
  title: '',
  description: '',
  ddl: ''
})

const addTaskRules = {
  title: [
    { required: true, message: '请输入任务标题', trigger: 'blur' },
    { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入任务内容', trigger: 'blur' }
  ],
  ddl: [
    { required: true, message: '请选择截止时间', trigger: 'change' }
  ]
}

// 邀请好友加入小组相关数据
const inviteVisible = ref(false)
const friends = ref([])
const friendsNull = ref(false)
const inviteTaskId = ref(0)
const currentInviteCode = ref('')

// 加入小组相关数据
const showJoinGroup = ref(false)
const joinGroupForm = ref({
  inviteCode: ''
})
const joinGroupFormRef = ref(null)

// 加入小组表单验证规则
const joinGroupRules = {
  inviteCode: [
    { required: true, message: '请输入邀请码', trigger: 'blur' },
    { pattern: /^[A-Za-z0-9]{6}$/, message: '邀请码必须是6位字母数字组合', trigger: 'blur' }
  ]
}

// 滚动事件监听
const handleScroll = () => {
  const { scrollTop, scrollHeight, clientHeight } = document.documentElement
  if (scrollTop + clientHeight >= scrollHeight - 100 && !loading.value && hasMore.value) {
    fetchTasks()
  }
}


// 获取任务列表
const fetchTasks = async () => {
  // 没有更多数据了
  if (!hasMore.value) return

  loading.value = true
  try{
    const res = await getTeamTaskList(currentPage.value, pageSize.value,userStore.userInfo.id)
    if(res.data.status === 1){
      // 没有更多数据
      if (res.data.data == null){
        hasMore.value = false
        return
      }

      tasks.value.push(...res.data.data)
      currentPage.value++
    }else{
      ElMessage.error("获取任务列表失败")
      console.log(res.data.message)
    }
  }finally{
    loading.value = false
  }
}


// 复制邀请码
const copyInviteCode = () => {
  navigator.clipboard.writeText(currentInviteCode.value)
    .then(() => {
      ElMessage.success('邀请码已复制到剪贴板')
    })
    .catch(() => {
      ElMessage.error('复制失败，请手动复制')
    })
}

// 处理加入小组
const handleJoinGroup = async () => {
  if (!joinGroupForm.value.inviteCode) {
    ElMessage.warning('请输入邀请码')
    return
  }
  
  try{
    // 验证邀请码格式
    await joinGroupFormRef.value.validateField('inviteCode')

    // 加入小组
    const res = await joinTeamByCode(joinGroupForm.value.inviteCode)
    if (res.data.status === 1) {
      ElMessage.success('成功加入小组')
      showJoinGroup.value = false
      joinGroupForm.value.inviteCode = ''
      // 重新加载任务列表
      currentPage.value = 1
      tasks.value = []
      fetchTasks()
    } else {
      ElMessage.error(res.data.message || '加入失败')
    }
  } catch (error) {
    ElMessage.error('表单填写格式不正确')
  }
}

// 显示邀请对话框
const handleInviteMember = async (taskId) => {
  inviteVisible.value = true
  inviteTaskId.value = taskId

  // 获取当前任务的邀请码
  const res0 = await getTeamTaskInviteCode(taskId)
  if(res0.data.status === 1) {
    currentInviteCode.value = res0.data.data
  } else {
    currentInviteCode.value = '加载邀请码失败'
  }

  // 加载好友数据
  const res = await getFriendList(userStore.userInfo.id)
  if(res.data.status === 1){
    if(res.data.data === null){
      friendsNull.value = true
      return
    }
    friends.value = res.data.data
  }else{
    friendsNull.value = true
    ElMessage.error('显示好友列表异常')
  }
}

// 邀请好友加入小组
const handleInviteFriend = async (friendId) => { 
  // 向好友发送邀请信息
  const data = {
    task_id: inviteTaskId.value,
    user_id: friendId
  }

  const res = await inviteMember(data)
  if(res.data.status === 1){
    ElMessage.success(res.data.message)
  }else{
    ElMessage.error('发送邀请失败！')
  }
}

// 检查好友是否已经在小组中
const isInvited = (friendId) => {
  // 先找到任务对应的小组成员
  const members = tasks.value.find( task => task.id === inviteTaskId.value).users
  return !members.some(member => member.id === friendId)
}

// 关闭邀请对话框
const handleInviteClear = () => {
  inviteVisible.value = false
  inviteTaskId.value = 0
  friends.value = []
  friendsNull.value = false
}

const handleCompleteTask = async (taskId) => {
  // 调用API完成任务
  const data = {
    task_id:  taskId,
    user_id:  userStore.userInfo.id,
    status:   1
  }
  const res = await completeTeamTask(data)
  if(res.data.status === 1){
    ElMessage.success('成功完成了属于您的部分')
    // 更新任务完成状态
    tasks.value = tasks.value.map(task => {
      if(task.id === taskId){
        var count = 0
        task.users.map(user => {
          if(user.id === userStore.userInfo.id){
            user.status = 1
          }
          count += user.status
          return user
        })
        // 小组成员全部完成了对应部分，则任务完成
        task.status = (count === task.users.length ?  1 : 0)
      }
      return task
    })
  }else{
    ElMessage.error('完成任务失败')
    console.log(res.data.message)
  }
}

// 处理退出小组
const handleExitTeam = async (taskId) => {
  // 调用API退出小组
  dialogVisible.value = true
  dialogTitle.value = '退出小组同时也会删除任务记录，您确定要退出吗？'
  // 等待用户操作
  const isConfirm = await confirmDialog.value.confirm()

  if(isConfirm){
    const res = await removeTeamTask(taskId,userStore.userInfo.id)
    if(res.data.status === 1){
      ElMessage.success('退出小组成功')
      tasks.value = tasks.value.filter(task => task.id !== taskId)
    }else{
      ElMessage.error('退出小组失败')
    }
    dialogVisible.value = false
  }
}

// 添加任务
const submitAddTask = async () => {
  try {
    await addTaskFormRef.value.validate()
    
    // 调用API添加任务
    const data = {
      user_id: userStore.userInfo.id,
      ...addTaskForm.value,
      type:   1,
      status: 0,
    }
    
    const res = await addTeamTask(data)
    if(res.data.status === 1){
      ElMessage.success('添加任务成功')

      // 重新加载任务列表
      currentPage.value = 1
      hasMore.value = true
      tasks.value = []
      fetchTasks()
      // 重置相关数据
      showAddTask.value = false
      addTaskFormRef.value.resetFields()
      addTaskForm.value = {
        title: '',
        description: '',
        ddl: '',
      }
    }else{
      ElMessage .error('添加任务失败')
    }
  } catch (error) {
    ElMessage.error('表单填写格式不正确')
  }
}

onMounted(() => {
  // 监听滚动事件
  window.addEventListener('scroll', handleScroll)
  // 初始获取任务数据
  fetchTasks()
})

onBeforeUnmount(() => {
  // 组件销毁时取消滚动事件监听
  window.removeEventListener('scroll', handleScroll)
})

</script>

<style scoped>
.team-page {
  max-width: 1500px;
  margin: 0 auto;
  padding: 20px;
  position: relative;
}

.task-list {
  display: flex;
  flex-direction: column;
  align-items: center; /* 使卡片居中 */
  gap: 20px;
  margin-top: 20px;
}

.task-list > * {
  width: 100%;
  max-width: 650px; /* 控制卡片的宽度 */
}

.publish-btn {
  position: fixed;
  right: 40px;
  bottom: 40px;
  width: 60px;
  height: 60px;
  font-size: 24px;
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

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.invite-code-section {
  text-align: center;
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f8f8f8;
  border-radius: 8px;
}

.invite-code-display {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
}

.invite-code {
  font-size: 18px;
  font-weight: bold;
  color: #409eff;
  padding: 8px 16px;
  background-color: #ecf5ff;
  border-radius: 4px;
}

.invite-results {
  margin-top: 20px;
}

.invite-card {
  display: flex;
  align-items: center;
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  margin-bottom: 10px;
}

.invite-card .el-avatar {
  margin-right: 15px;
}

.invite-card .user-info {
  flex: 1;
}

.invite-card .username {
  font-weight: bold;
}
.friends-null{
  text-align: center;
  padding: 20px;
  color: #999;
  font-size: 14px;
}
</style>