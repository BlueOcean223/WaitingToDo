<template>
  <el-card 
    class="task-card" 
    @click="handleCardClick"
    shadow="hover"
  >
    <div class="card-content">
      <h3>{{ task.title }}</h3>
      <p class="content">{{ task.description }}</p>
      <el-tag :type="task.status === 0 ? 'info' : 'success'">
        {{ formatDate(task.ddl) }}
      </el-tag>
      
      <div class="card-footer">
        <div class="user-avatars">
          <el-avatar 
            v-for="user in task.users" 
            :key="user.id"
            :size="30"
            :src="BaseUrl + user.pic"
          />
        </div>
        <div class="actions">
          <el-button 
            type="primary" 
            size="small" 
            @click.stop="handleInvite"
          >邀请</el-button>
          <el-button
            type="warning"
            size="small" 
            @click.stop="handleExit"
          >退出</el-button>
        </div>
      </div>
    </div>
  </el-card>

  <!-- 任务详情对话框 -->
  <el-dialog v-model="showDetail" :title="task.title" width="600px" :append-to-body="true">
    <div class="task-detail">
      <p class="detail-description">{{ task.description }}</p>
      <el-tag :type="task.status === 0 ? 'info' : 'success'">
        截止时间: {{ formatDate(task.ddl) }}
      </el-tag>
      
      <div class="detail-users">
        <h4>小组成员</h4>
        <div class="user-list">
          <div v-for="user in task.users" :key="user.id" class="user-item">
            <el-avatar 
              :size="40"
              :src="BaseUrl + user.pic"
            />
            <span class="user-name">{{ user.name }}</span>
            <el-button 
              :disabled="user.id !== userStore.userInfo.id"
              :type="user.status === 1 ? 'success' : 'primary'"
              size="small"
              @click="handleComplete"
            >
              {{ user.status === 1 ? '已完成' : '完成' }}
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'

const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['invite', 'complete','exitTeam'])

const BaseUrl = import.meta.env.VITE_PIC_BASE_URL
const userStore = useUserStore()
const showDetail = ref(false)

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

const handleCardClick = () => {
  showDetail.value = true
}

// 退出小组
const handleExit = () => {
  emit('exitTeam',props.task.id)
}


const handleInvite = () => {
  emit('invite', props.task.id)
}

const handleComplete = () => {
  if (props.task.status === 0) {
    emit('complete', props.task.id)
  }
}
</script>

<style scoped>
.task-card {
  margin: 10px auto;
  max-width: 1000px;
  transition: transform 0.3s ease;
  cursor: pointer;
}

.task-card:hover {
  transform: scale(1.02);
}

.content {
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin: 10px 0;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 15px;
}

.user-avatars {
  display: flex;
  gap: 5px;
}

.actions {
  display: flex;
  gap: 5px;
}

.task-detail {
  padding: 10px;
}

.detail-description {
  color: #666;
  margin-bottom: 20px;
  white-space: pre-wrap;
}

.detail-users {
  margin-top: 20px;
}

.user-list {
  margin-top: 10px;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 10px 0;
  padding: 5px;
  border-radius: 4px;
}

.user-item:hover {
  background-color: #f5f5f5;
}

.user-name {
  flex: 1;
}
</style>