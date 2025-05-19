<template>
  <Navbar :activeIndex="$route.path" />
  <div class="user-profile-container">
    <!-- 主体内容 -->
    <div class="profile-content">
      <!-- 头像区域 -->
      <div class="avatar-section">
        <el-tooltip content="更新头像" placement="bottom">
          <el-upload
            class="avatar-uploader"
            :show-file-list="false"
            :before-upload="beforeAvatarUpload"
            :http-request="uploadImgToServer"
          >
            <el-avatar :size="150" :src="avatarUrl" class="user-avatar">
              <img src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png"/>
            </el-avatar>
          </el-upload>
        </el-tooltip>
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
        <el-button type="primary" size="large" @click="showEditDialog">修改信息</el-button>
        <el-button type="warning" size="large" @click="showPasswordDialog">修改密码</el-button>
        <el-button type="danger" size="large" @click="handleLogout">退出登录</el-button>
      </div>
    </div>

    <!-- 修改信息对话框 -->
    <el-dialog v-model="editDialogVisible" title="修改个人信息" width="500px">
      <el-form 
        :model="editForm" 
        label-width="80px" 
        label-position="left"
        :rules="editRules"
        ref="editFormRef"
      >
        <el-form-item label="邮箱">
          <el-input v-model="userInfo.email" disabled prefix-icon="Message"/>
        </el-form-item>
        <el-form-item label="昵称" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入昵称" size="large" prefix-icon="User"/>
        </el-form-item>
        <el-form-item label="简介" prop="description">
          <el-input 
            v-model="editForm.description" 
            type="textarea" 
            :rows="4"
            placeholder="请输入个人简介"
            maxlength="200"
            show-word-limit
            prefix-icon="Info"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false" size="large">取消</el-button>
        <el-button type="primary" @click="submitEdit" size="large">保存修改</el-button>
      </template>
    </el-dialog>

    <!-- 修改密码对话框 - 第一步 -->
    <el-dialog v-model="passwordDialogVisible" title="修改密码" width="500px" v-if="!verifySuccess">
      <el-form 
        :model="passwordForm" 
        label-width="100px" 
        label-position="left"
        :rules="captchaRules"
        ref="captchaFormRef"
      >
        <el-form-item label="邮箱">
          <el-input v-model="userInfo.email" disabled size="large"  prefix-icon="Message"/>
        </el-form-item>
        <el-form-item label="验证码" prop="captcha">
          <div class="captcha-input">
            <el-input 
              v-model="passwordForm.captcha" 
              placeholder="请输入验证码" 
              size="large"
              prefix-icon="Key"
            />
            <el-button 
              :disabled="countdown > 0"
              @click="sendCaptcha"
              size="large"
            >
              {{ countdown > 0 ? `${countdown}秒后重试` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false" size="large">取消</el-button>
        <el-button type="primary" @click="verifyCaptcha" size="large" :loading="verifying">
          下一步
        </el-button>
      </template>
    </el-dialog>

    <!-- 修改密码对话框 - 第二步 -->
    <el-dialog v-model="passwordDialogVisible" title="修改密码" width="500px" v-else>
      <el-form 
        :model="passwordForm" 
        label-width="100px" 
        label-position="left"
        :rules="passwordRules"
        ref="passwordFormRef"
      >
        <el-form-item label="新密码" prop="newPassword">
          <el-input 
            v-model="passwordForm.newPassword" 
            type="password" 
            placeholder="请输入新密码"
            show-password
            size="large"
            prefix-icon="Lock"
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="passwordForm.confirmPassword" 
            type="password" 
            placeholder="请再次输入新密码"
            show-password
            size="large"
            prefix-icon="Lock"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="verifySuccess = false" size="large">上一步</el-button>
        <el-button type="primary" @click="submitPassword" size="large" :loading="submitting">
          确认修改
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup name="PersonalPage">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import Navbar from '../../components/Navbar.vue'
import { useUserStore } from '@/stores/user'
import { captcha } from '@/api/auth'
import { checkCaptcha, reset, update } from '@/api/user'
import { uploadImg } from '@/api/upload'


// 全局用户信息
const userStore = useUserStore()

// 用户信息
const userInfo = ref({
  name: '',
  email: '',
  description: '',
  pic: ''
})

// 获取用户信息
onMounted(() => {
  const user = JSON.parse(localStorage.getItem('user'))
  if(user) {
    userInfo.value = user
  }
})

// 头像计算属性
const avatarUrl = computed(() => {
  return `${import.meta.env.VITE_PIC_BASE_URL}${userInfo.value.pic}`
})

// 修改信息相关
const editDialogVisible = ref(false) // 修改信息对话框可见性
const editForm = ref({ // 修改信息表单数据
  name: '',
  description: ''
})
const editFormRef = ref(null) // 修改信息表单引用
const editRules = { // 修改信息表单验证规则
  name: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '长度不能超过 200 个字符', trigger: 'blur' }
  ]
}

const showEditDialog = () => {
  editForm.value = {
    name: userInfo.value.name,
    description: userInfo.value.description || ''
  }
  editDialogVisible.value = true
}

const submitEdit = async () => {
  try {
    await editFormRef.value.validate()
    
    // 调用API提交修改
    const res = await update(editForm.value)
    if (res.data.status === 1){
      ElMessage.success(res.data.message)
      // 更新本地存储用户信息
      userInfo.value.name = editForm.value.name
      userInfo.value.description = editForm.value.description

      localStorage.setItem('user', JSON.stringify(userInfo.value))

      // 更新全局用户信息
      userStore.updateUserInfo(userInfo.value)
    }else{
      ElMessage.error(res.data.message)
    }
    editDialogVisible.value = false
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error('修改失败: ' + error.message)
    }
  }
}

// 修改密码相关
const passwordDialogVisible = ref(false) // 修改密码对话框可见性
const verifySuccess = ref(false) //  验证成功
const countdown = ref(0) // 验证码倒计时
const verifying = ref(false) // 验证码验证中
const submitting = ref(false) // 提交中
const passwordForm = ref({ // 修改密码表单数据
  captcha: '',
  newPassword: '',
  confirmPassword: ''
})
const captchaFormRef = ref(null) // 验证码表单引用
const passwordFormRef = ref(null) // 修改密码表单引用

const captchaRules = {
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度为6位', trigger: 'blur' }
  ]
}

const passwordRules = {
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validatePassword, trigger: 'blur' }
  ]
}

function validatePassword(rule, value, callback) {
  if (value !== passwordForm.value.newPassword) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const showPasswordDialog = () => {
  passwordDialogVisible.value = true
  verifySuccess.value = false
  passwordForm.value = {
    captcha: '',
    newPassword: '',
    confirmPassword: ''
  }
}

const sendCaptcha = async () => {
  try {
    // 调用API发送验证码
    await captcha(userInfo.value.email)
    
    ElMessage.success('验证码已发送至您的邮箱')
    startCountdown()
  } catch (error) {
    ElMessage.error('验证码发送失败')
  }
}

const startCountdown = () => {
  countdown.value = 60
  const timer = setInterval(() => {
    countdown.value--
    if(countdown.value <= 0) clearInterval(timer)
  }, 1000)
}

// 验证码验证
const verifyCaptcha = async () => {
  try {
    verifying.value = true
    await captchaFormRef.value.validate()
    
    // 调用API验证验证码
    const data = {
      captcha: passwordForm.value.captcha
    }

    const res = await checkCaptcha(data)
    if(res.data.status === 1){
      verifySuccess.value = true
      ElMessage.success('验证码验证成功')
    }else{
      ElMessage.error(res.data.message)
    }
    
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error('验证失败: ' + error.message)
    }
  } finally {
    verifying.value = false
  }
}


// 修改密码提交
const submitPassword = async () => {
  try {
    submitting.value = true
    await passwordFormRef.value.validate()
    
    // 调用API修改密码
    const data = {
      password: passwordForm.value.newPassword
    }

    const res = await reset(data)
    if(res.data.status === 1){
      ElMessage.success('密码修改成功')
      passwordDialogVisible.value = false
    }else{
      ElMessage.error('密码修改失败: ' + res.data.message)
    }
    
    
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error('修改失败: ' + error.message)
    }
  } finally {
    submitting.value = false
  }
}


// 头像上传前的校验
const beforeAvatarUpload = (file) => {
  // 检查文件格式和文件大小，限制格式为 .png/.jpg/.jpeg
  const allowedExtensions = ['.png', '.jpg', '.jpeg']
  const fileExtension = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
  const isFormatValid = allowedExtensions.includes(fileExtension)
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isFormatValid) {
    ElMessage.error('只能上传 .png/.jpg/.jpeg 格式的图片')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过5MB')
    return false
  }
  return true
}

// 向后端上传图片
const uploadImgToServer = async (options) => {
  // 构造formdata
  const { file } = options
  const formData = new FormData()
  formData.append('file', file)

  const res = await uploadImg(formData)
  if(res.data.status === 1){
    // 上传成功
    // 更新本地存储的用户信息
    userInfo.value.pic = res.data.data
    localStorage.setItem('user', JSON.stringify(userInfo.value))
    // 更新全局用户信息
    userStore.updateUserInfo(userInfo.value)
    ElMessage.success(res.data.message)
  }else{
    // 上传失败
    ElMessage.error(res.data.message)
  }
}

// 退出登录
const handleLogout = () => {
    //  清除本地存储的token和用户信息
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    // 跳转到登录页
    window.location.href = '/login'
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

.captcha-input {
  display: flex;
  gap: 10px;
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