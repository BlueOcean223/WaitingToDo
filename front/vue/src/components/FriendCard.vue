<template>
  <el-card 
    class="friend-card" 
    shadow="hover" 
    @click="$emit('click')"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false"
    :style="{ transform: isHovered ? 'scale(1.03)' : 'scale(1)' }"
  >
    <div class="card-content">
      <div class="first-line">
        <el-avatar :size="50" :src="picBaseUrl + friend.pic" />
        <span class="friend-name">{{ friend.name }}</span>
      </div>
      <div class="second-line">
        <p class="friend-description">{{ friend.description === '' ? "用户暂无简介" : friend.description}}</p>
      </div>
      <div class="actions">
        <el-button type="primary" size="small" @click.stop="$emit('send-message', friend.id)">
          捎句话
        </el-button>
        <el-button type="danger" size="small" @click.stop="$emit('delete-friend', friend.id)">
          删除
        </el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  friend: {
    type: Object,
    required: true
  }
})

defineEmits(['click', 'send-message', 'delete-friend'])

const isHovered = ref(false)
const picBaseUrl = import.meta.env.VITE_PIC_BASE_URL
</script>

<style scoped>
.friend-card {
  cursor: pointer;
  transition: transform 0.3s ease;
}

.card-content {
  padding: 10px;
}

.first-line {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.friend-name {
  margin-left: 15px;
  font-size: 16px;
  font-weight: bold;
}

.friend-description {
  color: #666;
  font-size: 14px;
  margin-bottom: 15px;
}

.actions {
  display: flex;
  justify-content: space-between;
}
</style>