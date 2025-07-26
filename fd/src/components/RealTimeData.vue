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