<template>
    <Navbar :activeIndex="navbarIndex"/>
    <div class="user-profile-container">
        <div class="profile-content">
        <!-- 头像区域 -->
        <div class="avatar-section">
            <el-avatar :size="150" :src="picBaseUrl + userInfo.pic" class="user-avatar">
                <img src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png"/>
            </el-avatar>
        </div>

        <!-- 用户信息展示 -->
        <div class="info-card">
            <el-descriptions :column="1" border size="large">
            <el-descriptions-item label="昵称" label-class-name="info-label">
                <span class="info-value">{{ userInfo.name }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="邮箱" label-class-name="info-label">
                <span class="info-value">{{ userInfo.email }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="简介" label-class-name="info-label">
                <span class="info-value">{{ userInfo.description || '暂无简介' }}</span>
            </el-descriptions-item>
            </el-descriptions>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons">
            <el-button type="danger" size="large" @click="handleDelete">删除好友</el-button>
        </div>
        </div>
    </div>
</template>

<script setup>
    import Navbar from '@/components/Navbar.vue'
    import { onMounted,ref } from 'vue';
    import { useRoute } from 'vue-router';
    import { getFiendInfo } from '@/api/friend';

    const route = useRoute();
    const userId = ref(null);
    const userInfo = ref({});
    const navbarIndex = '/friend'
    const picBaseUrl = import.meta.env.VITE_PIC_BASE_URL

    onMounted(() => {
        userId.value = route.params.id
        console.log(userId.value)
        fetchUserInfo()
    })

    const fetchUserInfo = async () => {
        const res = await getFiendInfo(userId.value)
        if(res.data.status === 1){
            userInfo.value = res.data.data
        }else{
            ElMessage.error(res.data.message)
        }
    }

</script>

<style scoped>
.user-profile-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 40px 20px;
}

.profile-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 36px;
  background-color: #fff;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.avatar-section {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.user-avatar {
  border: 4px solid #f0f0f0;
  transition: all 0.3s;
}

.user-avatar:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.info-card {
  width: 100%;
  max-width: 600px;
}

:deep(.info-label) {
  font-size: 16px;
  font-weight: 500;
  color: #606266;
  width: 100px;
}

.info-value {
  font-size: 16px;
  color: #303133;
  white-space: normal; 
  word-wrap: break-word;
}

.action-buttons {
  display: flex;
  gap: 20px;
  margin-top: 20px;
}


:deep(.el-descriptions__body) {
  background-color: #f9f9f9;
  width: 100%;
}

:deep(.el-descriptions__header) {
  margin-bottom: 16px;
}

:deep(.el-dialog__body) {
  padding: 20px 30px;
}

:deep(.el-form-item__label) {
  font-size: 15px;
}
</style>