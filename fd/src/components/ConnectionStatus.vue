<template>
  <div class="connection-status" :class="statusClass">
    WebSocket状态: {{ statusText }}
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import wsService from '../services/websocket.service'

const status = ref(wsService.getStatus())
const statusText = ref('')
const statusClass = ref('')

// 状态映射
const statusMap: Record<string, { text: string; class: string }> = {
  'connecting': { text: '连接中...', class: 'connecting' },
  'connected': { text: '已连接', class: 'connected' },
  'disconnected': { text: '已断开', class: 'disconnected' }
}

// 更新状态显示
function updateStatus() {
  const statusInfo = statusMap[status.value] || { text: '未知状态', class: 'unknown' }
  statusText.value = statusInfo.text
  statusClass.value = statusInfo.class
}

// 监听状态变化
import { watch } from 'vue'

const stopWatch = watch(
  () => wsService.connectionStatus.value,
  (newStatus) => {
    status.value = newStatus
    updateStatus()
  }
)

onMounted(updateStatus)
onBeforeUnmount(stopWatch)
</script>

<style scoped>
.connection-status {
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: bold;
  text-align: center;
  margin: 10px 0;
}

.connecting {
  background-color: #ffcc00;
  color: #333;
}

.connected {
  background-color: #4caf50;
  color: white;
}

.disconnected {
  background-color: #f44336;
  color: white;
}

.unknown {
  background-color: #9e9e9e;
  color: white;
}
</style>