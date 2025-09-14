<template>
  <div class="home-container">
    <NavBar :activeIndex="$route.path" />

    <div class="main-content">
      <!-- 任务列表区域 -->
      <div class="task-list">
        <!-- 添加筛选按钮 -->
        <div class="filter-controls">
          <el-radio-group v-model="filterStatus" @change="handleFilterChange">
            <el-radio-button label="all">全部任务</el-radio-button>
            <el-radio-button label="uncompleted">未完成</el-radio-button>
          </el-radio-group>
        </div>

        <TaskCard 
          v-for="task in tasks" 
          :key="task.id" 
          :task="task"
          @complete="handleComplete"
          @change="handleChange"
          @delete="handleDelete"
        />
        <div v-if="loading" class="loading">加载中...</div>
        <!-- 没有更多数据时，提示用户已加载完 -->
        <div v-if="!hasMore && !loading" class="no-more">全部任务都在这里了</div>
      </div>
      
      <!-- 截止提醒区 -->
      <div class="deadline-reminder">
        <h3>赶不上ddl啦！</h3>
        <el-scrollbar>
          <div 
            v-for="task in urgentTasks"
            :key="task.id"
            class="urgent-task"
            @click="openUrgentTask(task)"
          >
            {{ task.title }} - {{ formatDate(task.ddl) }}
          </div>
          <div v-if="urgentTasks==null" class="no-more">暂时不用慌</div>
        </el-scrollbar>
      </div>
    </div>

    <!-- 发布按钮 -->
    <el-button 
      type="primary" 
      circle 
      class="publish-btn"
      @click="handlePublish"
    >
      <el-icon><Plus /></el-icon>
    </el-button>
  </div>

  <!-- 添加任务模态框 -->
  <el-dialog v-model="showAddTask" title="添加任务" width="500px">
    <el-form 
      ref="addTaskFormRef"
      :model="addTaskForm" 
      label-width="80px" 
      label-position="left"
      :rules="addTaskRules"
    >
      <el-form-item label="标题" prop="title">
        <el-input v-model="addTaskForm.title" placeholder="请输入标题"/>
      </el-form-item>
      <el-form-item label="内容" prop="description">
        <el-input 
          v-model="addTaskForm.description" 
          type="textarea" 
          :rows="8"
          placeholder="请输入任务内容"
          maxlength="500"
          show-word-limit
          prefix-icon="Info"
        />
      </el-form-item>
      <el-form-item label="截止时间" prop="ddl">
        <el-date-picker 
          v-model="addTaskForm.ddl" 
          type="datetime" 
          placeholder="请选择截止时间" 
          size="large"
        />
      </el-form-item>
      <el-form-item label="附件">
        <el-upload
          :file-list="attachments"
          class="upload-demo"
          drag
          multiple
          :limit="3"
          :auto-upload="false"
          :on-exceed="handleExceed"
          :on-change="handleFileChange"
          :on-remove="handleFileRemove"
          accept=".doc,.docx,.pdf"
        >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">
            拖拽文件到此处或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持上传 docx/doc/pdf 文件，单个文件不超过30M，最多上传3个文件
            </div>
          </template>
        </el-upload>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="showAddTask = false" size="large">取消</el-button>
      <el-button type="primary" @click="submitAddTask" size="large">完成</el-button>
    </template>
  </el-dialog>

  <!-- 修改任务模态框 -->
    <el-dialog v-model="showEditTask" title="修改任务" width="500px">
      <el-form 
        :model="editTaskForm" 
        label-width="80px" 
        label-position="left"
        :rules="addTaskRules"
        ref="editTaskFormRef"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="editTaskForm.title" placeholder="请输入标题"/>
        </el-form-item>
        <el-form-item label="内容" prop="description">
          <el-input 
            v-model="editTaskForm.description" 
            type="textarea" 
            :rows="8"
            placeholder="请输入任务内容"
            maxlength="500"
            show-word-limit
            prefix-icon="Info"
          />
        </el-form-item>
        <el-form-item label="截止时间" prop="ddl">
          <el-date-picker 
            v-model="editTaskForm.ddl" 
            type="datetime" 
            placeholder="请选择截止时间" 
            size="large"
          />
        </el-form-item>
        <el-form-item label="附件">
          <el-upload
            :file-list="editTaskForm.attachments"
            class="upload-demo"
            drag
            multiple
            :limit="3"
            :auto-upload="false"
            :on-exceed="handleExceed"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            accept=".doc,.docx,.pdf"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              拖拽文件到此处或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持上传 docx/doc/pdf 文件，单个文件不超过20M，最多上传3个文件
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditTask = false" size="large">取消</el-button>
        <el-button type="primary" @click="submitEditTask" size="large">完成</el-button>
      </template>
    </el-dialog>

  <!-- 任务详情（来自截止提醒区，复用 TaskDetail） -->
  <el-dialog v-model="urgentDetailVisible" :title="detailTask ? detailTask.title : ''" width="550px" class="task-detail-dialog">
    <TaskDetail v-if="detailTask" :task="detailTask" />
    <template #footer>
      <el-button type="primary" @click="urgentDetailVisible = false">关闭</el-button>
    </template>
  </el-dialog>

  <!-- 确认窗口 -->
  <ConfirmDialog
    ref="confirmDialogRef"
    v-model="dialogVisible"
    :title="dialogTitle"
  />


</template>

<script setup name="HomeView">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import NavBar from '@/components/NavBar.vue'
import TaskCard from '@/components/TaskCard.vue'
import TaskDetail from '@/components/TaskDetail.vue'
import ConfirmDialog  from '@/components/ConfirmDialog.vue'
import { Plus,UploadFilled,Delete } from '@element-plus/icons-vue'
import { ElMessage,ElLoading } from 'element-plus'
import { getList, add, remove, update, getUrgent } from '@/api/task'
import { uploadFile,deleteFile } from '@/api/upload'

// 响应式数据
const tasks = ref([])
const urgentTasks = ref([])
const filterStatus = ref('all') // 'all' 或者 'uncompleted' 筛选状态
const loading = ref(false)
const currentPage = ref(1)
const showAddTask = ref(false) // 是否显示添加任务表单
const showEditTask = ref(false) // 是否显示修改任务表单
const hasMore = ref(true) // 是否有更多数据
const addTaskForm = ref({
  title: '',
  description: '',
  ddl: '',
})
const addTaskRules = {
  title: [
    { required: true, message: '请输入任务标题', trigger: 'blur' } // 失去焦点时触发
  ],
  description: [
    { required: true, message: '请输入任务描述', trigger: 'blur' }
  ],
  ddl: [
    { required: true, message: '请选择任务截止时间', trigger: ['blur','change'] } // 失去焦点与改变时均触发
  ]
}
const editTaskForm = ref({
  id: '',
  title: '',
  description: '',
  ddl: '',
  attachments: [],  // 包含所有文件（已有+新增）
  deletedAttachments: [], // 记录被删除的已有文件
})
const dialogVisible = ref(false)
const dialogTitle = ref('')
const attachments = ref([]) // 添加任务附件
const loadingInstance = ref(null) // 加载实例
// 新增：截止提醒区的详情弹窗
const urgentDetailVisible = ref(false)
const detailTask = ref(null)

// 模板引用
const addTaskFormRef = ref(null)
const editTaskFormRef = ref(null)
const confirmDialogRef = ref(null)
// 方法定义
const fetchTasks = async () => {
      try {
    // 如果没有更多数据，直接返回
    if (!hasMore.value) return;

    loading.value = true
    // 根据筛选状态传递不同的status给后端
    const status = filterStatus.value === 'uncompleted' ? 0 :undefined
    // 调用API接口
    const res = await getList(currentPage.value,5,status)

    if(res.data.status === 1){
      if(res.data.data == null){
        // 没有更多数据
        hasMore.value = false
        return;
      }
      tasks.value = [...tasks.value, ...res.data.data]
      currentPage.value++
    }else{
      ElMessage.error(res.data.message)
    }
  } finally {
    loading.value = false
  }
}
// 筛选变化处理方法
const handleFilterChange = () => { 
  currentPage.value = 1
  hasMore.value = true
  tasks.value = []
  fetchTasks()
}
// 查询紧急任务
const fetchUrgentTasks = async () => {
  try{
    const res = await getUrgent()
    if(res.data.status === 1){
      urgentTasks.value = res.data.data
    }else{
      ElMessage.error(res.data.message)
    }
  }catch(error){
    ElMessage.error(error.message)
  }
}
const handleScroll = () => {
  const { scrollTop, clientHeight, scrollHeight } = document.documentElement
  if (scrollTop + clientHeight >= scrollHeight - 100 && !loading.value && hasMore.value) {
    fetchTasks()
  }
}
const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}
const handlePublish = () => {
  // 显示添加任务模态框
  showAddTask.value = true
}
// 新增：从截止提醒区打开详情
const openUrgentTask = (task) => {
  // 若已加载完整任务（含附件等），优先展示完整数据
  const full = tasks.value.find(t => t.id === task.id)
  detailTask.value = full || task
  urgentDetailVisible.value = true
}
// 提交新增任务
const submitAddTask = async () => {
  try{
    // 检查表格参数是否合法
    await addTaskFormRef.value.validate()

    loadingInstance.value = ElLoading.service({
      lock: true,
      text: '正在提交...',
      background: 'rgba(0, 0, 0, 0.7)'
    })

    // 将数据发送给后端
    const data = {
      ...addTaskForm.value,
      type: 0
    }


    const res = await add(data)
    if(res.data.status === 1){
      // 上传文件
      if(attachments.value.length > 0){
        // 构造formdata
        const formData = new FormData()
        attachments.value.forEach(file => {
          formData.append('files',file.raw) // filelist获取的是封装过的file对象
        })
        var res0 = await uploadFile(res.data.data.id,formData)
        // 如果上传失败，则重试三次
        var cnt = 3
        if(res0.data.status !== 1 && cnt > 0){
          // 休眠一分钟再尝试
          cnt-=1
          await new Promise(resolve => setTimeout(resolve, 60000))
          res0 = await uploadFile(res.data.data.id,formData)
        }
      }

      // 添加成功
      ElMessage.success(res.data.message)
      showAddTask.value = false
      addTaskForm.value = {
        title: '',
        description: '',
        ddl: '',
      }
      attachments.value = []
      // 刷新任务列表
      currentPage.value = 1
      hasMore.value = true
      tasks.value = []
      fetchTasks()
      // 刷新紧急任务列表
      fetchUrgentTasks()
    }else{
      // 添加失败
      ElMessage.error(res.data.message)
    }
  }catch(error){
    if(error instanceof Error){
      ElMessage.error('任务发布失败: ' + error.message)
    }
  } finally{
    loadingInstance.value.close()
  }
}
// 提交修改任务
const submitEditTask = async () => { 
  // 检查表格参数是否合法
  try{
    await editTaskFormRef.value.validate()

    loadingInstance.value = ElLoading.service({
      lock: true,
      text: '正在提交...',
      background: 'rgba(0, 0, 0, 0.7)'
    })


    // 需要删除的旧文件
    if(editTaskForm.value.deletedAttachments.length > 0){
      const deletedFileIds = editTaskForm.value.deletedAttachments.map(file => file.id).filter(Boolean)
      await deleteFile(deletedFileIds)
      editTaskForm.value.deletedAttachments = [] // 重置删除的旧文件
    }

    // 需要新上传的文件
    const formData = new FormData()
    const needUploadFiles = editTaskForm.value.attachments.filter(file => !file.isExisting)
    if (needUploadFiles.length > 0) { 
      needUploadFiles.forEach(file => {
        formData.append('files', file.raw)
      })
      const res1 = await uploadFile(editTaskForm.value.id, formData)
      if(res1.data.status !== 1) {
        ElMessage.error("上传文件失败")
        return
      }
      // 添加上传文件的url
      editTaskForm.value.attachments.forEach(file => {
        if (!file.isExisting){
          file.url = "/files/" + editTaskForm.value.id + "/" + file.raw.name
        }
      })
    }

    // 将数据发送给后端
    const data = {
      ...editTaskForm.value,
      type: 0,
      // 移除不需要的字段
      attachments: undefined,
      deletedAttachments: undefined
    }
    const res = await update(data)
    if(res.data.status === 1){
      // 修改成功
      ElMessage.success(res.data.message)
      showEditTask.value = false
      // 更新任务列表中的当前任务信息
      for(let i = 0; i < tasks.value.length; i++){
        if(tasks.value[i].id === data.id){
          tasks.value[i].title = data.title
          tasks.value[i].description = data.description
          tasks.value[i].ddl = data.ddl
          tasks.value[i].attachments = editTaskForm.value.attachments
          break
        }
      }
    }else{
      ElMessage.error(res.data.message)
    }
  }catch(error){
    if(error instanceof Error){
      ElMessage.error('任务修改失败: ' + error.message)
    }
  } finally{
    loadingInstance.value.close()
  }
}
// 完成任务
const handleComplete = async (task) => {
  dialogTitle.value = '这个任务搞定了吗？'
  dialogVisible.value = true
  const isConfirmed = await confirmDialogRef.value.confirm()
  if (isConfirmed){
    // 用户确认执行操作
    task.status = 1
    const res = await update(task)
    if(res.data.status === 1){
      // 任务完成
      ElMessage.success(res.data.message)
      // 将任务列表中的该任务更新
      tasks.value = tasks.value.map(item => item.id === task.id ? task : item)
    }else{
      // 任务未完成
      ElMessage.error(res.data.message)
    }
  }
}
// 修改任务
const handleChange = (task) => {
  editTaskForm.value.id = task.id
  editTaskForm.value.title = task.title
  editTaskForm.value.description = task.description
  editTaskForm.value.ddl = task.ddl
  editTaskForm.value.deletedAttachments = []

  // 初始化已有附件，添加标记
  if(task.attachments != null){
    editTaskForm.value.attachments = task.attachments.map(file => ({
      ...file,
      isExisting: true, // 标记为已有附件
      uid: file.id // 使用唯一标识
    }))
  }else {
    editTaskForm.value.attachments = []
  }

  showEditTask.value = true
}
// 删除任务
const handleDelete = async (id) => {
  dialogTitle.value = '确定要删除该任务吗？'
  dialogVisible.value = true
  const isConfirmed = await confirmDialogRef.value.confirm()
  if (isConfirmed){
    // 用户确认执行操作
    const res = await remove(id)
    if(res.data.status === 1){
      // 删除成功
      ElMessage.success(res.data.message)
      // 刷新任务列表,删除该任务
      tasks.value = tasks.value.filter(item => item.id !== id)
    }else{
      // 删除失败
      ElMessage.error(res.data.message)
    }
  }
}
// 文件上传相关方法
const handleExceed = (files, fileList) => {
  ElMessage.warning("最多只能上传3个文件");
}

const handleFileChange = (file, fileList) => {
  const allowedExt = [".doc", ".docx", ".pdf"]; 
  const maxSize = 20 * 1024 * 1024; // 20MB
  // 过滤出合法文件
  fileList = fileList.filter(f => {
    const ext = f.name.slice(f.name.lastIndexOf(".")).toLowerCase();
    if (!allowedExt.includes(ext)) {
      ElMessage.warning(`文件格式不支持`);
      return false;
    }
    if (f.size > maxSize) {
      ElMessage.warning(`文件超过大小20MB`);
      return false;
    }
    return true;
  });
  
  // 更新附件列表
  if (showAddTask.value){
    attachments.value = fileList;
  }else if (showEditTask.value) {
    const updatedList = fileList.map(file => {
      // 如果是已有文件，保留原有标记
      if(file.isExisting) return file;
      // 新文件添加标记
      return {
        ...file,
        isExisting: false,
        uid: Date.now()
      }
    })

    editTaskForm.value.attachments = updatedList;
  }
}

const handleFileRemove = (file, fileList) => {
  // 更新附件列表
  if(showAddTask.value){
    attachments.value = fileList;
  }else if(showEditTask.value) {
    editTaskForm.value.attachments = fileList;

    // 如果删除的是已有文件，则加入删除列表
    if(file.isExisting) {
      editTaskForm.value.deletedAttachments.push(file)
    }
  }
}

const removeFile = (formType, index) => {
  if (formType === 'add') {
    attachments.value.splice(index, 1);
  } else {
    editTaskForm.value.attachments.splice(index, 1);
  }
}

// 生命周期钩子
onMounted(() => {
  fetchTasks()
  fetchUrgentTasks()
  window.addEventListener('scroll', handleScroll)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})

</script>

<style scoped>
.home-container {
  min-height: 100vh;
}
.main-content {
  display: flex;
  gap: 20px;
  padding: 20px;
}
.task-list {
  flex: 4;
}
.deadline-reminder {
  flex: 1;
  background: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  text-align: center;
  height: 400px;
}
.publish-btn {
  position: fixed;
  bottom: 40px;
  right: 40px;
  width: 56px;
  height: 56px;
}
.loading {
  text-align: center;
  padding: 20px;
  color: #666;
}
.no-more {
  text-align: center;
  padding: 20px;
  color: #999;
  font-size: 14px;
}
.urgent-task {
  color: red;
  cursor: pointer;
  padding: 6px 0;
  text-align: center;
}
.filter-controls {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-end;
}
.upload-demo {
  width: 100%;
}

.delete-btn {
  margin-left: 10px;
}
/* 缩小上传拖拽区域 */
:deep(.el-upload-dragger) {
  padding: 12px 15px; 
  height: 120px; 
}
/* 调整上传图标大小 */
:deep(.el-icon--upload) {
  font-size: 40px; 
  margin: 8px 0;
}

/* 调整提示文字样式 */
:deep(.el-upload__text) {
  font-size: 13px; 
  line-height: 1.4;
  margin: 4px 0;
}

/* 任务详情弹窗样式（复用 TaskCard） */
.task-detail-dialog {
  border-radius: 8px;
}
</style>