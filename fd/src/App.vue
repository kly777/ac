<script setup lang="ts">
import { onMounted } from 'vue'
import RealTimeData from './components/RealTimeData.vue'
import ConnectionStatus from './components/ConnectionStatus.vue'
import AIResponse from './components/AIResponse.vue'
import wsService from './services/websocket.service'

onMounted(() => {
  wsService.connect()
})
</script>

<template>
  <div class="app-container">
    <AIResponse class="ai-response" />
    <div class="main-container">
      <ConnectionStatus />
      <RealTimeData />
    </div>
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: Arial, sans-serif;
  background-color: #f5f5f5;
}

.app-container {
  box-sizing: border-box;
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: row;
  padding: 20px;
  gap: 20px;

  background-color:#333;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.ai-response {
  flex: 1;
  height: 100%;
  background-color: #f5f5f5;
}

.main-container {
  flex: 1;
}

h1 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}
</style>
