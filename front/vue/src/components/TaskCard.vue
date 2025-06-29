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

  <!-- 详情对话框 -->
  <el-dialog v-model="detailVisible" :title="task.title" width="550px">
    <div>
      <p style="white-space: pre-wrap;">{{ task.description }}</p>
      <el-tag :type="task.status === 0 ? 'info' : 'success'">截止时间：{{ formatDate(task.ddl) }}</el-tag>
      <p><strong>状态：</strong> {{ task.status === 0 ? '未完成' : '已完成' }}</p>
    </div>
    <template #footer>
      <el-button @click="detailVisible = false">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script>
export default {
  name: 'TaskCard',
  props: {
    task: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      hover: false,
      detailVisible: false
    }
  },
  methods: {
    formatDate(dateStr) {
      return new Date(dateStr).toLocaleString()
    },
    handleComplete() {
      this.$emit('complete', this.task)
    },
    handleChange() {
      this.$emit('change', this.task)
    },
    handleDelete() {
      this.$emit('delete', this.task.id)
    },
    showDetail() {
      this.detailVisible = true
    }
  }
}
</script>

<style scoped>
.task-card {
  margin: 10px auto;
  max-width: 800px;
  cursor: pointer;
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
.completed{
  color:rgb(201, 20, 120)
}
</style>