<template>
  <Navbar :activeIndex="$route.path" />
  <div class="team-page">
    <div class="task-list">
      <TeamTask 
        v-for="task in tasks" 
        :key="task.id" 
        :task="task"
        @delete="handleDeleteTask"
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
  </div>
</template>

<script setup name="TeamPage">
import { ref,onMounted,onBeforeUnmount } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import TeamTask from '@/components/TeamTask.vue'
import Navbar from '@/components/Navbar.vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getTeamTaskList } from '@/api/task'

// 分页相关数据
const loading = ref(false)
const currentPage = ref(1)
const hasMore = ref(true)
const pageSize = ref(5) // 每页显示的任务数量

// 当前用户全局信息
const userStore = useUserStore()

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

const handleDeleteTask = (taskId) => {
  // TODO: 调用API删除任务
  tasks.value = tasks.value.filter(task => task.id !== taskId)
}

const handleInviteMember = (taskId) => {
  // TODO: 调用API邀请成员
  console.log('邀请成员到任务', taskId)
}

const handleCompleteTask = (taskId) => {
  // TODO: 调用API完成任务
  const task = tasks.value.find(t => t.id === taskId)
  if (task) {
    task.status = 1
  }
}

// 处理退出小组
const handleExitTeam = (taskId) => {
  // TODO: 调用API退出小组
  ElMessage.success('已退出小组')
}

const submitAddTask = async () => {
  try {
    await addTaskFormRef.value.validate()
    
    // TODO: 调用API添加任务
    const newTask = {
      id: Date.now(), // 临时ID，实际应由后端生成
      ...addTaskForm.value,
      status: 0,
      users: [] // 初始没有成员，或者包含当前用户
    }
    
    tasks.value.push(newTask)
    showAddTask.value = false
    addTaskFormRef.value.resetFields()
  } catch (error) {
    console.error('表单验证失败', error)
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
</style>