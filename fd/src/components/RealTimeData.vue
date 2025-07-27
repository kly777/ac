<template>
  <div class="real-time-data">
    <h2>实时数据</h2>
    <div class="data-section">
      <h3>Informer内容</h3>
      <ul>
        <li v-for="(item, index) in informerData" :key="'informer-' + index">
          {{ item.title }}: {{ item.content }}
        </li>
      </ul>
    </div>
    <div class="data-section">
      <h3>Manager命令</h3>
      <ul>
        <li v-for="(cmd, index) in managerCommands" :key="'command-' + index">
          {{ cmd.type }}: {{ cmd.content }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as wsEvents from '../services/ws-events.service'

// 定义数据类型
interface InformerItem {
  title: string
  content: string
}

interface ManagerCommand {
  type: string
  content: string
}

interface NewInfo {
  title: string
  content: string
}

// 响应式数据
const informerData = ref<InformerItem[]>([])
const managerAIOutput = ref('')
const managerCommands = ref<ManagerCommand[]>([])
const newInfo = ref<NewInfo>({ title: '', content: '' })

// 事件处理函数
const handleInformer = (data: any) => {
  informerData.value.push(data)
}

const handleManagerAI = (data: any) => {
  managerAIOutput.value += data
}

const handleManagerCommand = (data: any) => {
  managerCommands.value.push({
    type: data.Type,
    content: data.Content
  })
}

onMounted(() => {
  // 注册事件监听
  wsEvents.onRealTimeData(handleInformer) // 使用实时数据通道传递informer消息
  // 注意：managerAI和managerCommand需要后端实现对应消息类型
})

onBeforeUnmount(() => {
  // 注销事件监听
  wsEvents.offRealTimeData(handleInformer)
})

// 添加自定义信息
const addCustomInfo = async () => {
  try {
    const response = await fetch('http://localhost:8080/add-info', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        title: newInfo.value.title,
        content: newInfo.value.content
      })
    })

    if (!response.ok) {
      throw new Error(`请求失败: ${response.status}`)
    }

    // 清空表单
    newInfo.value = { title: '', content: '' }
    alert('信息添加成功！')
  } catch (error) {
    console.error('添加信息失败:', error)
    const err = error as Error
    alert(`添加信息失败: ${err.message}`)
  }
}

</script>

<style scoped>
.real-time-data {
  padding: 20px;
  font-family: Arial, sans-serif;
}

.add-info-form {
  margin-bottom: 30px;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #f9f9f9;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

input,
textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-sizing: border-box;
}

textarea {
  min-height: 80px;
  resize: vertical;
}

button {
  background-color: #007bff;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}

.data-section {
  margin-bottom: 20px;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 4px;
}

h3 {
  margin-top: 0;
}

ul {
  padding-left: 20px;
}

.output {
  white-space: pre-wrap;
  background-color: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  max-height: 200px;
  overflow-y: auto;
}
</style>