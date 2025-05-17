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
        />
        <div v-if="loading" class="loading">加载中...</div>
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
</template>

<script>
import NavBar from '@/components/NavBar.vue'
import TaskCard from '@/components/TaskCard.vue'
import { Plus } from '@element-plus/icons-vue'

export default {
  name: 'HomeView',
  components: {
    NavBar,
    TaskCard,
    ElIconPlus: Plus
  },
  data() {
    return {
      tasks: [{title: 'Task 1',description: '测试数据测试数据',ddl: '2025/05/18', status: 0}],
      urgentTasks: [],
      loading: false,
      currentPage: 1
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
        this.loading = true
        // 调用API接口
        // const res = await this.$api.getTasks(...)
        // this.tasks = [...this.tasks, ...res.data]
        this.currentPage++
      } finally {
        this.loading = false
      }
    },
    handleScroll() {
      const { scrollTop, clientHeight, scrollHeight } = document.documentElement
      if (scrollTop + clientHeight >= scrollHeight - 100 && !this.loading) {
        this.fetchTasks()
      }
    },
    formatDate(dateStr) {
      return new Date(dateStr).toLocaleString()
    },
    handlePublish() {
      // 处理发布逻辑
    }
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
</style>