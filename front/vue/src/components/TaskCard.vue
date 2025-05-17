<template>
  <el-card class="task-card">
    <div class="card-content">
      <h3>{{ task.title }}</h3>
      <p class="content">{{ truncatedContent }}</p>
      <el-tag type="info">{{ formatDate(task.ddl) }}</el-tag>
      <div class="actions">
        <el-button 
          type="success" 
          :disabled="task.status === 1"
          @click="handleComplete"
        >完成</el-button>
        <el-button type="warning" @click="handleDelay">延期</el-button>
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
      this.$emit('complete', this.task.id)
    },
    handleDelay() {
      this.$emit('delay', this.task.id)
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
</style>