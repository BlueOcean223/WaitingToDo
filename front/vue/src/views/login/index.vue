<template>
  <div class="login-container">
    <el-card class="login-box">
      <!-- 登录表单 -->
      <template v-if="currentForm === 'login'">
        <h2 class="login-title">WaitingToDo</h2>
        <el-form 
          ref="loginForm" 
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
          ref="registerForm" 
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
          ref="forgetForm" 
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

<script>
import { ElMessage } from 'element-plus';
import { login, forget, captcha, register } from '@/api/auth'

export default{
  name: 'LoginPage',
  
  data(){
    // 确认密码验证规则
    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== this[this.currentForm + 'Form'].password) {
        callback(new Error('两次输入密码不一致!'));
      } else {
        callback();
      }
    };

    return {
      currentForm: 'login', // login, register, forget
      loading: false,
      
      // 登录表单
      loginForm: {
        email: '',
        password: ''
      },
      loginRules: {
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 16, message: '密码长度在 6 到 16 个字符', trigger: 'blur' }
        ]
      },
      
      // 注册表单
      registerForm: {
        nickname: '',
        email: '',
        captcha: '',
        password: '',
        confirmPassword: ''
      },
      registerRules: {
        nickname: [
          { required: true, message: '请输入昵称', trigger: 'blur' },
          { min: 1, max: 12, message: '昵称长度在 1 到 12 个字符', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        captcha: [
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
      },
      registerCodeDisabled: false,
      registerCodeText: '发送验证码',
      registerCountdown: 0,
      registerTimer: null,
      
      // 忘记密码表单
      forgetForm: {
        email: '',
        captcha: '',
        password: '',
        confirmPassword: ''
      },
      forgetRules: {
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        captcha: [
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
      },
      forgetCodeDisabled: false,
      forgetCodeText: '发送验证码',
      forgetCountdown: 0,
      forgetTimer: null
    }
  },

  methods: {
    // 切换表单
    switchForm(formType) {
      this.currentForm = formType;
    },
    
    // 登录
    async handleLogin(){
      this.$refs.loginForm.validate(async valid => {
        if (valid){
          try{
            this.loading = true
            const res = await login(this.loginForm)
            if (res.data.status === 1){
              ElMessage.success(res.data.message)
              // 保存token和用户信息
              localStorage.setItem('token',res.data.token)
              localStorage.setItem('user',JSON.stringify(res.data.data))// 存放前先将对象转换为字符串
              // 跳转到首页
              this.$router.push('/home')
            }else {
              ElMessage.error(res.data.message)
            }
            this.loading = false
          }catch (error) {
            this.loading = false
            ElMessage.error(error.message)
          }
        }else {
          ElMessage.error('请填写正确的用户名和密码')
          return
        }
      })
    },
    
    // 注册
    async handleRegister() {
      this.$refs.registerForm.validate(async valid => {
        if (valid) {
          try{
            // 注册
            this.loading = true;
            // 封装数据
            const data = {
              name: this.registerForm.nickname,
              email: this.registerForm.email,
              captcha: this.registerForm.captcha,
              password: this.registerForm.password,
            }
            // 发送注册请求
            const res = await register(data)
            // 判断注册成功
            if (res.data.status === 1) { 
              ElMessage.success(res.data.message);
              // 清空表单数据
              this.registerForm = {
                nickname: '',
                email: '',
                captcha: '',
                password: '',
                confirmPassword: ''
              };
              // 重置表单验证状态
              this.$refs.registerForm.resetFields();
              // 切换到登录表单
              this.switchForm('login');
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
        this.loading = false;
      });
    },
    
    // 发送注册验证码
    async sendRegisterCode() {
      try {
        // 先检验邮箱字段，不符合规则时直接抛出异常
        await this.$refs.registerForm.validateField('email');

        // 发送验证码
        await captcha(this.registerForm.email)

        // 六十秒内不能再次发送
        this.registerCodeDisabled = true;
        this.registerCountdown = 60;
        this.registerCodeText = `${this.registerCountdown}秒后重新发送`;
        
        // 设置定时器
        this.registerTimer = setInterval(() => {
          this.registerCountdown--;
          this.registerCodeText = `${this.registerCountdown}秒后重新发送`;
            
          if (this.registerCountdown <= 0) {
            clearInterval(this.registerTimer);
            this.registerCodeDisabled = false;
            this.registerCodeText = '发送验证码';
          }
        }, 1000);
          
        ElMessage.success('验证码已发送');
      } catch (error){
        return;
      }
    },
    
    // 重置密码
    async handleForget() {
      this.$refs.forgetForm.validate(async valid => {
        if (valid) {
          try{
            this.loading = true;
            // 封装数据
            const data = {
              email: this.forgetForm.email,
              captcha: this.forgetForm.captcha,
              password: this.forgetForm.password,
            }
            // 发送请求
            const res = await forget(data)
            // 重置密码成功
            if (res.data.status === 1) { 
              ElMessage.success(res.data.message);

              // 清空表单数据
              this.forgetForm = {
                email: '',
                captcha: '',
                password: '',
                confirmPassword: ''
              }
              // 重置表单验证状态
              this.$refs.forgetForm.resetFields();
              // 切换到登录表单
              this.switchForm('login');
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
        this.loading = false;
      });
    },
    
    // 发送忘记密码验证码
    async sendForgetCode() {
      try{
        // 先校验邮箱字段，不符合规则直接抛出异常
        await this.$refs.forgetForm.validateField('email');

        // 发送验证码
        await captcha(this.forgetForm.email)

        // 六十秒内不能重复触发
        this.forgetCodeDisabled = true;
        this.forgetCountdown = 60;
        this.forgetCodeText = `${this.forgetCountdown}秒后重新发送`;
          
        // 设置定时器
        this.forgetTimer = setInterval(() => {
          this.forgetCountdown--;
          this.forgetCodeText = `${this.forgetCountdown}秒后重新发送`;
            
          if (this.forgetCountdown <= 0) {
            clearInterval(this.forgetTimer);
            this.forgetCodeDisabled = false;
            this.forgetCodeText = '发送验证码';
          }
        }, 1000);
        ElMessage.success('验证码已发送');
      } catch (error) {
        return;
      }
    },

  },

  beforeUnmount() {
    // 清除定时器
    if (this.registerTimer) clearInterval(this.registerTimer);
    if (this.forgetTimer) clearInterval(this.forgetTimer);
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-image: url('@/assets/sysuBackground.png');
  background-size: cover;
  background-position: center;
}

.login-box {
  width: 400px;
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