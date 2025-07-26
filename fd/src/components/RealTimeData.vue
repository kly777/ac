<template>
  <div class="real-time-data">
    <h2>实时数据</h2>
    
    <div class="add-info-form">
      <h3>添加自定义信息</h3>
      <form @submit.prevent="addCustomInfo">
        <div class="form-group">
          <label for="title">标题:</label>
          <input type="text" id="title" v-model="newInfo.title" required>
        </div>
        <div class="form-group">
          <label for="content">内容:</label>
          <textarea id="content" v-model="newInfo.content" required></textarea>
        </div>
        <button type="submit">提交</button>
      </form>
    </div>

    <div class="data-section">
      <h3>Informer内容</h3>
      <ul>
        <li v-for="(item, index) in informerData" :key="'informer-' + index">
          {{ item.title }}: {{ item.content }}
        </li>
      </ul>
    </div>

    <div class="data-section">
      <h3>ManagerAI输出</h3>
      <div class="output">{{ managerAIOutput }}</div>
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

interface NewInfo {
  title: string
  content: string
}

// 定义数据类型
interface InformerItem {
  title: string
  content: string
}

interface ManagerCommand {
  type: string
  content: string
}

// 响应式数据
const informerData = ref<InformerItem[]>([])
const managerAIOutput = ref('')
const managerCommands = ref<ManagerCommand[]>([])

// 新信息表单
const newInfo = ref<NewInfo>({ title: '', content: '' })

// WebSocket连接
let socket: WebSocket | null = null

onMounted(() => {
  // 建立WebSocket连接
  socket = new WebSocket('ws://localhost:8080/ws')
  if (socket===null) {
    console.log('WebSocket连接失败')
  }

  socket.onopen = () => {
    console.log('WebSocket连接已建立')
  }

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)

      switch (data.type) {
        case 'informer':
          informerData.value.push(data.data)
          break
        case 'managerAI':
          managerAIOutput.value += data.data
          break
        case 'managerCommand':
          managerCommands.value.push({
            type: data.data.Type,
            content: data.data.Content
          })
          break
      }
    } catch (e) {
      console.error('解析WebSocket数据失败:', e)
    }
  }

  socket.onerror = (error) => {
    console.error('WebSocket错误:', error)
  }

  socket.onclose = () => {
    console.log('WebSocket连接已关闭')
  }
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

onBeforeUnmount(() => {
  if (socket) {
    socket.close()
  }
})
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

input, textarea {
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