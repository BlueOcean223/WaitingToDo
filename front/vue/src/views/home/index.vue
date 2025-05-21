<template>
  <div class="home-container">
    <NavBar />
    
    <div class="main-content">
      <!-- 任务列表区域 -->
      <div class="task-list">
        <TaskCard 
          v-for="task in tasks" 
          :key="task.id" 
          :task="task"
          @complete="handleComplete"
          @change="handleChange"
          @delete="handleDelete"
        />
        <div v-if="loading" class="loading">加载中...</div>
        <!-- 没有更多数据时，提示用户已加载完 -->
        <div v-if="!hasMore && !loading" class="no-more">全部任务都在这里了</div>
      </div>
      
      <!-- 截止提醒区 -->
      <div class="deadline-reminder">
        <h3>赶不上ddl啦！</h3>
        <el-scrollbar>
          <div 
            v-for="task in urgentTasks"
            :key="task.id"
            class="urgent-task"
          >
            {{ task.title }} - {{ formatDate(task.ddl) }}
          </div>
          <div v-if="!urgentTasks.length" class="no-more">暂时不用慌</div>
        </el-scrollbar>
      </div>
    </div>

    <!-- 发布按钮 -->
    <el-button 
      type="primary" 
      circle 
      class="publish-btn"
      @click="handlePublish"
    >
      <el-icon><Plus /></el-icon>
    </el-button>
  </div>

<!--  添加任务模态框 -->
    <el-dialog v-model="showAddTask" title="添加任务" width="500px">
      <el-form 
        :model="addTaskForm" 
        label-width="80px" 
        label-position="left"
        :rules="addTaskRules"
        ref="addTaskForm"
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
            prefix-icon="Info"
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

  <!-- 修改任务模态框 -->
    <el-dialog v-model="showEditTask" title="修改任务" width="500px">
      <el-form 
        :model="editTaskForm" 
        label-width="80px" 
        label-position="left"
        :rules="addTaskRules"
        ref="editTaskForm"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="editTaskForm.title" placeholder="请输入标题"/>
        </el-form-item>
        <el-form-item label="内容" prop="description">
          <el-input 
            v-model="editTaskForm.description" 
            type="textarea" 
            :rows="8"
            placeholder="请输入任务内容"
            maxlength="500"
            show-word-limit
            prefix-icon="Info"
          />
        </el-form-item>
        <el-form-item label="截止时间" prop="ddl">
          <el-date-picker 
            v-model="editTaskForm.ddl" 
            type="datetime" 
            placeholder="请选择截止时间" 
            size="large"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditTask = false" size="large">取消</el-button>
        <el-button type="primary" @click="submitEditTask" size="large">完成</el-button>
      </template>
    </el-dialog>

  <!-- 确认窗口 -->
  <ConfirmDialog
    ref="confirmDialog"
    v-model="dialogVisible"
    :title="dialogTitle"
  />


</template>

<script>
import NavBar from '@/components/NavBar.vue'
import TaskCard from '@/components/TaskCard.vue'
import ConfirmDialog  from '@/components/ConfirmDialog.vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getList, add, remove, update } from '@/api/task'

export default {
  name: 'HomeView',
  components: {
    NavBar,
    TaskCard,
    ConfirmDialog,
    ElIconPlus: Plus
  },
  data() {
    return {
      tasks: [{title: 'Task 1',description: '测试数据测试数据',ddl: '2025/05/18', status: 0},{title: 'Task 2',description: '测试数据测试数据',ddl: '2025/06/01', status: 1}],
      urgentTasks: [],
      loading: false,
      currentPage: 1,
      showAddTask: false, // 是否显示添加任务表单
      showEditTask: false, // 是否显示修改任务表单
      hasMore: true, // 是否有更多数据
      addTaskForm:{
        title: '',
        description: '',
        ddl: ''
      },
      addTaskRules: { 
        title: [
          { required: true, message: '请输入任务标题', trigger: 'blur' } // 失去焦点时触发
        ],
        description: [
          { required: true, message: '请输入任务描述', trigger: 'blur' }
        ],
        ddl: [
          { required: true, message: '请选择任务截止时间', trigger: ['blur','change'] } // 失去焦点与改变时均触发
        ]
      },
      editTaskForm: {
        id: '',
        title: '',
        description: '',
        ddl: ''
      },
      confirmDialog: {},
      dialogVisible: false,
      dialogTitle: '',
    }
  },
  mounted() {
    window.addEventListener('scroll', this.handleScroll)
    this.fetchTasks()
  },
  beforeDestroy() {
    window.removeEventListener('scroll', this.handleScroll)
  },
  methods: {
    async fetchTasks() {
      try {
        // 如果没有更多数据，直接返回
        if (!this.hasMore) return;

        this.loading = true
        // 调用API接口
        const res = await getList(this.currentPage,5)
        if(res.data.status === 1){
          // 如果是第一页，直接赋值
          if (this.currentPage === 1){
            this.tasks = res.data.data
          }else{// 否则追加数据
            this.tasks = [...this.tasks, ...res.data.data]
          }
          
          // 判断是否还有更多数据
          this.hasMore = this.tasks.length < res.data.data.count
          this.currentPage++
        }else{
          ElMessage.error(res.data.message)
        }
      } finally {
        this.loading = false
      }
    },
    handleScroll() {
      const { scrollTop, clientHeight, scrollHeight } = document.documentElement
      if (scrollTop + clientHeight >= scrollHeight - 100 && !this.loading && this.hasMore) {
        this.fetchTasks()
      }
    },
    formatDate(dateStr) {
      return new Date(dateStr).toLocaleString()
    },
    handlePublish() {
      // 显示添加任务模态框
      this.showAddTask = true
    },
    // 提交新增任务
    async submitAddTask(){
      try{
        // 检查表格参数是否合法
        await this.$refs.addTaskForm.validate()

        // 将数据发送给后端
        const data = {
          ...this.addTaskForm,
          type: 0
        }
        const res = await add(data)
        if(res.data.status === 1){
          // 添加成功
          ElMessage.success(res.data.message)
          this.showAddTask = false
          this.addTaskForm = {
            title: '',
            description: '',
            deadline: ''
          }
          // 刷新任务列表
          this.currentPage = 1
          this.hasMore = true
          this.fetchTasks()
        }else{
          // 添加失败
          ElMessage.error(res.data.message)
        }
      }catch(error){
        if(error instanceof Error){
          ElMessage.error('任务发布失败: ' + error.message)
        }
      }
    },
    // 提交修改任务
    async submitEditTask(){ 
      // 检查表格参数是否合法
      try{
        await this.$refs.editTaskForm.validate()
        // 将数据发送给后端
        const data = {
          ...this.editTaskForm,
          type: 0
        }
        const res = await update(data)
        if(res.data.status === 1){
          // 修改成功
          ElMessage.success(res.data.message)
          this.showEditTask = false
          // 更新任务列表中的当前任务信息
          for(let i = 0; i < this.tasks.length; i++){
            if(this.tasks[i].id === data.id){
              this.tasks[i].title = data.title
              this.tasks[i].description = data.description
              this.tasks[i].ddl = data.ddl
              break
            }
          }
        }else{
          ElMessage.error(res.data.message)
        }
      }catch(error){
        if(error instanceof Error){
          ElMessage.error('任务修改失败: ' + error.message)
        }
      }
    },
    // 完成任务
    async handleComplete(task){
      this.dialogTitle = '这个任务搞定了吗？'
      this.dialogVisible = true
      const isConfirmed = await this.$refs.confirmDialog.confirm()
      if (isConfirmed){
        // 用户确认执行操作
        task.status = 1
        const res = await update(task)
        if(res.data.status === 1){
          // 任务完成
          ElMessage.success(res.data.message)
          // 将任务列表中的该任务更新
          this.tasks = this.tasks.map(item => item.id === task.id ? task : item)
        }else{
          // 任务未完成
          ElMessage.error(res.data.message)
        }
      }
    },
    // 修改任务
    handleChange(task){
      this.editTaskForm.id = task.id
      this.editTaskForm.title = task.title
      this.editTaskForm.description = task.description
      this.editTaskForm.ddl = task.ddl
      this.showEditTask = true
    },
    // 删除任务
    async handleDelete(id){
      console.log(id)
      this.dialogTitle = '确定要删除该任务吗？'
      this.dialogVisible = true
      const isConfirmed = await this.$refs.confirmDialog.confirm()
      if (isConfirmed){
        // 用户确认执行操作
        const res = await remove(id)
        if(res.data.status === 1){
          // 删除成功
          ElMessage.success(res.data.message)
          // 刷新任务列表,删除该任务
          this.tasks = this.tasks.filter(item => item.id !== id)
        }else{
          // 删除失败
          ElMessage.error(res.data.message)
        }
      }
    },
  }
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
}
.main-content {
  display: flex;
  gap: 20px;
  padding: 20px;
}
.task-list {
  flex: 4;
}
.deadline-reminder {
  flex: 1;
  background: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  text-align: center;
}
.publish-btn {
  position: fixed;
  bottom: 40px;
  right: 40px;
  width: 56px;
  height: 56px;
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