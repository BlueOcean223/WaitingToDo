<template>
  <el-dialog
    :model-value="modelValue"
    :title="title"
    width="30%"
    @update:model-value="$emit('update:modelValue', $event)"
    @close="handleClose"
    class="confirm-dialog"
  >
    <div class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleConfirm">确 定</el-button>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  title: String,
  modelValue: Boolean
})

const emit = defineEmits(['update:modelValue'])

const resolvePromise = ref(null)

const handleConfirm = () => {
  if (resolvePromise.value) {
    resolvePromise.value(true)
  }
  emit('update:modelValue', false)
}

const handleCancel = () => {
  if (resolvePromise.value) {
    resolvePromise.value(false)
  }
  emit('update:modelValue', false)
}

const handleClose = () => {
  if (resolvePromise.value) {
    resolvePromise.value(false)
  }
  emit('update:modelValue', false)
}

// 暴露一个可以调用的方法
const confirm = () => {
  return new Promise((resolve) => {
    resolvePromise.value = resolve
    emit('update:modelValue', true)
  })
}

defineExpose({
  confirm
})
</script>

<style scoped>
.confirm-dialog{
  text-align: center;
}
.dialog-footer {
  text-align: right;
}
</style>