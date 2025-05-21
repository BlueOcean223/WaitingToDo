<template>
  <!-- 未完成的任务卡片 -->
  <el-card class="task-card" v-if="task.status === 0">
    <div class="card-content">
      <h3>{{ task.title }}</h3>
      <p class="content">{{ truncatedContent }}</p>
      <el-tag type="info">{{ formatDate(task.ddl) }}</el-tag>
      <div class="actions">
        <el-button 
          type="success" 
          @click="handleComplete"
        >完成</el-button>
        <el-button type="warning" @click="handleChange">修改</el-button>
        <el-button type="danger" @click="handleDelete">删除</el-button>
      </div>
    </div>
  </el-card>

  <!-- 已完成的任务卡片 -->
  <el-card class="task-card" v-else>
    <div class="card-content">
      <h3>{{ task.title }}</h3>
      <p class="content">{{ truncatedContent }}</p>
      <el-tag type="info">{{ formatDate(task.ddl) }}</el-tag>
      <div class="actions">
        <span class="completed">该任务已完成！</span>
        <el-button type="danger" @click="handleDelete">删除</el-button>
      </div>
    </div>
  </el-card>
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
  computed: {
    truncatedContent() {
      const content = this.task.description || ''
      return content.slice(0, 20) + (content.length > 20 ? '...' : '')
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
    }
  }
}
</script>

<style scoped>
.task-card {
  margin: 10px auto;
  max-width: 800px;
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