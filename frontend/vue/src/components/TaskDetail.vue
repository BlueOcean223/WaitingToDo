<template>
  <div class="detail-content" v-if="task">
    <div class="description-section">
      <h4>任务描述</h4>
      <p class="description-text">{{ task.description }}</p>
    </div>

    <div class="info-section">
      <div class="info-item">
        <span class="info-label">截止时间：</span>
        <el-tag :type="task.status === 0 ? 'info' : 'success'">{{ formatDate(task.ddl) }}</el-tag>
      </div>

      <div class="info-item">
        <span class="info-label">状态：</span>
        <el-tag :type="task.status === 0 ? 'warning' : 'success'">
          {{ task.status === 0 ? '未完成' : '已完成' }}
        </el-tag>
      </div>
    </div>

    <div v-if="task.attachments && task.attachments.length" class="attachments-section">
      <h4>附件</h4>
      <div class="attachments-list">
        <div
          v-for="attachment in task.attachments"
          :key="attachment.id"
          class="attachment-item"
          @click="handleFileClick(attachment)"
        >
          <el-image
            class="attachment-icon"
            :src="getFileUrl(attachment.name)"
            fit="contain"
          />
          <span class="attachment-name">{{ attachment.name }}</span>
        </div>
      </div>
      <div class="attachment-tip">
        <el-icon><InfoFilled /></el-icon>
        <span>温馨提示：点击可预览PDF文件，其它格式的文件点击仅为下载</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { InfoFilled } from '@element-plus/icons-vue'

// 定义props
const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})

// 方法
const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

const getFileUrl = (fileName) => {
  const baseUrl = import.meta.env.VITE_PIC_BASE_URL
  const extensionName = fileName.substring(fileName.lastIndexOf("."))
  if (extensionName === '.pdf'){
    return baseUrl + "/images/pdf.png"
  } else {
    return baseUrl + "/images/word.png"
  }
}

const handleFileClick = (file) => {
  const baseUrl = import.meta.env.VITE_PIC_BASE_URL
  const fileUrl = baseUrl + file.url
  window.open(fileUrl, '_blank')
}
</script>

<style scoped>
.detail-content {
  padding: 0 10px;
}

.description-section {
  margin-bottom: 20px;
}

.description-section h4,
.attachments-section h4 {
  color: #409eff;
  margin-bottom: 10px;
  font-size: 16px;
}

.description-text {
  white-space: pre-wrap;
  line-height: 1.6;
  color: #555;
  padding: 10px;
  background-color: #f8f8f8;
  border-radius: 4px;
}

.info-section {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 20px;
}

.info-item {
  display: flex;
  align-items: center;
}

.info-label {
  font-weight: bold;
  color: #666;
  margin-right: 8px;
}

.attachments-section {
  margin-top: 20px;
}

.attachments-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.attachment-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 4px;
  background-color: #f5f7fa;
  cursor: pointer;
  transition: all 0.2s;
}

.attachment-item:hover {
  background-color: #e1f5fe;
  transform: translateX(5px);
}

.attachment-icon {
  width: 35px;
  height: 35px;
  margin-right: 10px;
}

.attachment-name {
  color: #333;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 400px;
}

.attachment-tip {
  margin-top: 10px;
  padding: 8px 12px;
  background-color: #f0f7ff;
  border-radius: 4px;
  color: #606266;
  font-size: 13px;
  display: flex;
  align-items: center;
}

.attachment-tip .el-icon {
  color: #409eff;
  margin-right: 8px;
  font-size: 16px;
}
</style>