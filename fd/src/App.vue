<script setup lang="ts">
import RealTimeData from './components/RealTimeData.vue'
</script>

<template>
  <div class="app-container">
    <h1>实时监控面板</h1>
    <RealTimeData />
  </div>
</template>

<style>
body {
  margin: 0;
  padding: 20px;
  font-family: Arial, sans-serif;
  background-color: #f5f5f5;
}

.app-container {
  max-width: 1200px;
  margin: 0 auto;
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

h1 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}
</style>
