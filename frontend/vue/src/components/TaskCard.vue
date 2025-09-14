<template>
  <!-- 未完成的任务卡片 -->
  <el-card 
    class="task-card" 
    v-if="task.status === 0"
    @click="showDetail"
    @mouseenter="hover = true"
    @mouseleave="hover = false"
    :style="{ transform: hover ? 'scale(1.03)' : 'scale(1)', transition: 'transform 0.3s ease' }"
  >
    <div class="card-content">
      <h3>{{ task.title }}</h3>
      <p class="content">{{ task.description }}</p>
      <el-tag type="info">{{ formatDate(task.ddl) }}</el-tag>
      <div class="actions">
        <el-button 
          type="success" 
          @click.stop="handleComplete"
        >完成</el-button>
        <el-button type="warning" @click.stop="handleChange">修改</el-button>
        <el-button type="danger" @click.stop="handleDelete">删除</el-button>
      </div>
    </div>
  </el-card>

  <!-- 已完成的任务卡片 -->
  <el-card 
    class="task-card" 
    v-else
    @click="showDetail"
    @mouseenter="hover = true"
    @mouseleave="hover = false"
    :style="{ transform: hover ? 'scale(1.03)' : 'scale(1)', transition: 'transform 0.3s ease' }"
  >
    <div class="card-content">
      <h3>{{ task.title }}</h3>
      <p class="content">{{ task.description }}</p>
      <el-tag :type="task.status === 0 ? 'info' : 'success'">{{ formatDate(task.ddl) }}</el-tag>
      <div class="actions">
        <span class="completed">该任务已完成！</span>
        <el-button type="danger" @click.stop="handleDelete">删除</el-button>
      </div>
    </div>
  </el-card>

  <!-- 详情对话框（复用 TaskDetail） -->
  <el-dialog v-model="detailVisible" :title="task.title" width="550px" class="task-detail-dialog">
    <TaskDetail :task="task" />
    <template #footer>
      <el-button type="primary" @click="detailVisible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { InfoFilled } from '@element-plus/icons-vue'
import TaskDetail from '@/components/TaskDetail.vue'

// 定义props
const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})

// 定义emits
const emit = defineEmits(['complete', 'change', 'delete'])

// 响应式数据
const hover = ref(false)
const detailVisible = ref(false)

// 方法
const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

const handleComplete = () => {
  emit('complete', props.task)
}

const handleChange = () => {
  emit('change', props.task)
}

const handleDelete = () => {
  emit('delete', props.task.id)
}

const showDetail = () => {
  detailVisible.value = true
}
</script>

<style scoped>
.task-card {
  margin: 10px auto;
  max-width: 800px;
  cursor: pointer;
  transition: all 0.3s ease;
}
.task-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
.content {
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.actions {
  margin-top: 15px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
.completed {
  color: rgb(201, 20, 120);
  font-weight: bold;
}

/* 详情对话框样式（壳） */
.task-detail-dialog {
  border-radius: 8px;
}
</style>