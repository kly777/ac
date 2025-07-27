<template>
  <div class="ai-response">
    <h2>Chat</h2>

    <div class="response-container">
      <div class="answer">{{ response }}</div>
    </div>

    <div class="input">
      <input v-model="currentQuestion" placeholder="输入问题..." :disabled="isSending" @keyup.enter="sendQuestion" />
      <button @click="sendQuestion" :disabled="!currentQuestion || isSending">
        {{ isSending ? '发送中...' : '发送' }}
      </button>
    </div>
    <div v-if="error" class="error">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { request } from '@/services/http.service'



const response = ref("")
const currentQuestion = ref('')
const isSending = ref(false)
const error = ref('')


const sendQuestion = async () => {
  if (!currentQuestion.value.trim()) return;

  isSending.value = true;
  error.value = '';

  try {
    await request('http://localhost:8080/q', 'POST', { content: currentQuestion.value })
    currentQuestion.value = '';
  } catch (err) {
    error.value = '发送失败: ' + (err as Error).message;
  } finally {
    isSending.value = false;
  }
};
</script>

<style scoped>
.ai-response {
  position: relative;
  padding: 16px;
  font-family: Arial, sans-serif;
  border: 1px solid #eee;
  border-radius: 8px;
}

.response-container {
  max-height: 300px;
  overflow-y: auto;
  margin-bottom: 15px;
}

.response-item {
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.question {
  font-weight: bold;
  margin-bottom: 5px;
  color: #333;
}

.answer {
  padding-left: 15px;
  color: #666;
}

.input {

  position: absolute;
  bottom: 20px;
  display: flex;
  gap: 10px;
}

textarea {
  flex: 1;
  min-height: 80px;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  resize: vertical;
}

button {
  background-color: #007bff;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  align-self: flex-end;
}

button:hover {
  background-color: #0056b3;
}
</style>