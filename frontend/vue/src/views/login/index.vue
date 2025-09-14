<template>
  <div class="login-container">
    <el-card class="login-box">
      <!-- 登录表单 -->
      <template v-if="currentForm === 'login'">
        <h2 class="login-title">WaitingToDo</h2>
        <el-form 
          ref="loginFormRef" 
          :model="loginForm" 
          :rules="loginRules" 
          label-position="top"
          @keyup.enter="handleLogin"
        >
          <el-form-item label="邮箱" prop="email">
            <el-input 
              v-model="loginForm.email" 
              type="text"
              placeholder="请输入邮箱"
              prefix-icon="Message"
            />
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="loginForm.password" 
              placeholder="请输入密码"
              type="password"
              show-password
              prefix-icon="Lock"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              class="login-button" 
              :loading="loading"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
          
          <div class="form-footer">
            <span class="link-button" @click="switchForm('register')">注册账号</span>
            <span class="link-button forget-pwd" @click="switchForm('forget')">忘记密码</span>
          </div>
        </el-form>
      </template>

      <!-- 注册表单 -->
      <template v-else-if="currentForm === 'register'">
        <h2 class="login-title">注册账号</h2>
        <el-form 
          ref="registerFormRef" 
          :model="registerForm" 
          :rules="registerRules" 
          label-position="top"
          @keyup.enter="handleRegister"
        >
          <el-form-item label="昵称" prop="nickname">
            <el-input 
              v-model="registerForm.nickname" 
              placeholder="请输入昵称"
              prefix-icon="User"
            />
          </el-form-item>
          
          <el-form-item label="邮箱" prop="email">
            <el-input 
              v-model="registerForm.email" 
              placeholder="请输入邮箱"
              prefix-icon="Message"
            />
          </el-form-item>
          
          <el-form-item label="验证码" prop="code">
            <div class="code-input">
              <el-input 
                v-model="registerForm.captcha" 
                placeholder="请输入验证码"
                prefix-icon="Key"
              />
              <el-button 
                class="send-code-btn" 
                :disabled="registerCodeDisabled"
                @click="sendRegisterCode"
              >
                {{ registerCodeText }}
              </el-button>
            </div>
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="registerForm.password" 
              placeholder="请输入密码"
              type="password"
              show-password
              prefix-icon="Lock"
            />
          </el-form-item>
          
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="registerForm.confirmPassword" 
              placeholder="请再次输入密码"
              type="password"
              show-password
              prefix-icon="Lock"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              class="login-button" 
              :loading="loading"
              @click="handleRegister"
            >
              注册
            </el-button>
          </el-form-item>
          
          <div class="form-footer">
            <span class="link-button" @click="switchForm('login')">已有账号</span>
          </div>
        </el-form>
      </template>

      <!-- 忘记密码表单 -->
      <template v-else-if="currentForm === 'forget'">
        <h2 class="login-title">重置密码</h2>
        <el-form 
          ref="forgetFormRef" 
          :model="forgetForm" 
          :rules="forgetRules" 
          label-position="top"
          @keyup.enter="handleForget"
        >
          <el-form-item label="邮箱" prop="email">
            <el-input 
              v-model="forgetForm.email" 
              placeholder="请输入邮箱"
              prefix-icon="Message"
            />
          </el-form-item>
          
          <el-form-item label="验证码" prop="code">
            <div class="code-input">
              <el-input 
                v-model="forgetForm.captcha" 
                placeholder="请输入验证码"
                prefix-icon="Key"
              />
              <el-button 
                class="send-code-btn" 
                :disabled="forgetCodeDisabled"
                @click="sendForgetCode"
              >
                {{ forgetCodeText }}
              </el-button>
            </div>
          </el-form-item>
          
          <el-form-item label="新密码" prop="password">
            <el-input 
              v-model="forgetForm.password" 
              placeholder="请输入新密码"
              type="password"
              show-password
              prefix-icon="Lock"
            />
          </el-form-item>
          
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="forgetForm.confirmPassword" 
              placeholder="请再次输入新密码"
              type="password"
              show-password
              prefix-icon="Lock"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              class="login-button" 
              :loading="loading"
              @click="handleForget"
            >
              重置密码
            </el-button>
          </el-form-item>
          
          <div class="form-footer">
            <span class="link-button" @click="switchForm('login')">想起来了</span>
          </div>
        </el-form>
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import { login, forget, captcha, register } from '@/api/auth'
import { useUserStore } from '@/stores/user.js'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const currentForm = ref('login') // login, register, forget
const loading = ref(false)

// 登录表单
const loginForm = ref({
  email: '',
  password: ''
})

// 注册表单
const registerForm = ref({
  nickname: '',
  email: '',
  captcha: '',
  password: '',
  confirmPassword: ''
})

// 忘记密码表单
const forgetForm = ref({
  email: '',
  captcha: '',
  password: '',
  confirmPassword: ''
})

// 验证码相关状态
const registerCodeDisabled = ref(false)
const registerCodeText = ref('发送验证码')
const registerCountdown = ref(0)
const registerTimer = ref(null)

const forgetCodeDisabled = ref(false)
const forgetCodeText = ref('发送验证码')
const forgetCountdown = ref(0)
const forgetTimer = ref(null)

// ref引用
const loginFormRef = ref(null)
const registerFormRef = ref(null)
const forgetFormRef = ref(null)

// 验证规则
// 确认密码验证规则
const validateConfirmPassword = (rule, value, callback) => {
  const formName = currentForm.value + 'Form'
  let password = ''
  if (formName === 'registerForm') {
    password = registerForm.value.password
  } else if (formName === 'forgetForm') {
    password = forgetForm.value.password
  }
  
  if (value !== password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const loginRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 16, message: '密码长度在 6 到 16 个字符', trigger: 'blur' }
  ]
}

const registerRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 1, max: 12, message: '昵称长度在 1 到 12 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度为6位', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 16, message: '密码长度在 6 到 16 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const forgetRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度为6位', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 16, message: '密码长度在 6 到 16 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

// 方法
// 切换表单
const switchForm = (formType) => {
  currentForm.value = formType
}
    
// 登录
const handleLogin = async () => {
  loginFormRef.value.validate(async valid => {
    if (valid){
      try{
        loading.value = true
        const res = await login(loginForm.value)
        if (res.data.status === 1){
          ElMessage.success(res.data.message)
          // 保存token和用户信息
          localStorage.setItem('token',res.data.token)
          localStorage.setItem('user',JSON.stringify(res.data.data))// 存放前先将对象转换为字符串
          // 更新全局用户信息
          useUserStore().updateUserInfo(res.data.data)
          // 跳转到首页
          router.push('/')
        }else {
          ElMessage.error(res.data.message)
        }
        loading.value = false
      }catch (error) {
        loading.value = false
        ElMessage.error(error.message)
      }
    }else {
      ElMessage.error('请填写正确的用户名和密码')
      return
    }
  })
}
    
// 注册
const handleRegister = async () => {
  registerFormRef.value.validate(async valid => {
    if (valid) {
      try{
        // 注册
        loading.value = true;
        // 封装数据
        const data = {
          name: registerForm.value.nickname,
          email: registerForm.value.email,
          captcha: registerForm.value.captcha,
          password: registerForm.value.password,
        }
        // 发送注册请求
        const res = await register(data)
        // 判断注册成功
        if (res.data.status === 1) { 
          ElMessage.success(res.data.message);
          // 清空表单数据
          registerForm.value = {
            nickname: '',
            email: '',
            captcha: '',
            password: '',
            confirmPassword: ''
          };
          // 重置表单验证状态
          registerFormRef.value.resetFields();
          // 切换到登录表单
          switchForm('login');
        }
        // 注册失败
        else{
          ElMessage.error(res.data.message)
        }
        
      }catch (error) {
        ElMessage.error(error.message)
      }
    } else {
      ElMessage.error('请填写正确的注册信息');
    }
    loading.value = false;
  });
}
    
// 发送注册验证码
const sendRegisterCode = async () => {
  try {
    // 先检验邮箱字段，不符合规则时直接抛出异常
    await registerFormRef.value.validateField('email');

    // 发送验证码
    await captcha(registerForm.value.email)

    // 六十秒内不能再次发送
    registerCodeDisabled.value = true;
    registerCountdown.value = 60;
    registerCodeText.value = `${registerCountdown.value}秒后重新发送`;
    
    // 设置定时器
    registerTimer.value = setInterval(() => {
      registerCountdown.value--;
      registerCodeText.value = `${registerCountdown.value}秒后重新发送`;
        
      if (registerCountdown.value <= 0) {
        clearInterval(registerTimer.value);
        registerCodeDisabled.value = false;
        registerCodeText.value = '发送验证码';
      }
    }, 1000);
      
    ElMessage.success('验证码已发送');
  } catch (error){
    return;
  }
}
    
// 重置密码
const handleForget = async () => {
  forgetFormRef.value.validate(async valid => {
    if (valid) {
      try{
        loading.value = true;
        // 封装数据
        const data = {
          email: forgetForm.value.email,
          captcha: forgetForm.value.captcha,
          password: forgetForm.value.password,
        }
        // 发送请求
        const res = await forget(data)
        // 重置密码成功
        if (res.data.status === 1) { 
          ElMessage.success(res.data.message);

          // 清空表单数据
          forgetForm.value = {
            email: '',
            captcha: '',
            password: '',
            confirmPassword: ''
          }
          // 重置表单验证状态
          forgetFormRef.value.resetFields();
          // 切换到登录表单
          switchForm('login');
        }
        // 重置密码失败
        else{
          ElMessage.error(res.data.message)
        }
      }catch (error) {
        ElMessage.error(error.message);
      }
    } else {
      ElMessage.error('请填写正确的信息');
    }
    loading.value = false;
  });
}
    
// 发送忘记密码验证码
const sendForgetCode = async () => {
  try{
    // 先校验邮箱字段，不符合规则直接抛出异常
    await forgetFormRef.value.validateField('email');

    // 发送验证码
    await captcha(forgetForm.value.email)

    // 六十秒内不能重复触发
    forgetCodeDisabled.value = true;
    forgetCountdown.value = 60;
    forgetCodeText.value = `${forgetCountdown.value}秒后重新发送`;
      
    // 设置定时器
    forgetTimer.value = setInterval(() => {
      forgetCountdown.value--;
      forgetCodeText.value = `${forgetCountdown.value}秒后重新发送`;
        
      if (forgetCountdown.value <= 0) {
        clearInterval(forgetTimer.value);
        forgetCodeDisabled.value = false;
        forgetCodeText.value = '发送验证码';
      }
    }, 1000);
    ElMessage.success('验证码已发送');
  } catch (error) {
    return;
  }
}

// 生命周期钩子
onBeforeUnmount(() => {
  // 清除定时器
  if (registerTimer.value) clearInterval(registerTimer.value);
  if (forgetTimer.value) clearInterval(forgetTimer.value);
})
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  position: fixed;
  inset: 0;
  background: url('@/assets/sysuBackground.png')  center/cover no-repeat fixed;
}

.login-box {
  width: 350px;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
  background-color: rgba(255, 255, 255, 0.9);
}

.login-title {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.login-button {
  width: 100%;
  margin-top: 10px;
}

.form-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
}

.link-button {
  color: #409eff;
  cursor: pointer;
  font-size: 14px;
}

.link-button:hover {
  text-decoration: underline;
}

.forget-pwd {
  text-align: right;
}

.code-input {
  display: flex;
  align-items: center;
}

.send-code-btn {
  margin-left: 10px;
  width: 120px;
}

</style>