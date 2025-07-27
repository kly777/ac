<template>
  <div class="ai-response">
    <h2>Chat</h2>

    <div class="chat-container">
      <div v-for="(message, index) in messages" :key="index" :class="message.role">
        <strong>{{ message.role === 'user' ? '我' : 'AI' }}:</strong> {{ message.content }}
      </div>
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
import { ref, onMounted } from 'vue'
import { Q } from '@/services/api.service'
import { onAIResponse } from '@/services/ws-events.service'



interface Message {
  role: 'user' | 'ai';
  content: string;
}

const messages = ref<Message[]>([]);
const currentQuestion = ref('');
const isSending = ref(false);
const error = ref('');


onMounted(() => {
  onAIResponse((content) => {
    console.log('AI Response:', content);
    if(content.status==="done"){
    }else if(content.status==="ing"){
      messages.value[messages.value.length-1].content+=content.content;
    }

  });
});

const sendQuestion = async () => {
  if (!currentQuestion.value.trim()) return;

  isSending.value = true;
  error.value = '';

  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: currentQuestion.value
  });

  try {
    // 使用封装的Q方法发送问题
    await Q(currentQuestion.value);
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

.chat-container {
  height: 300px;
  overflow-y: auto;
  margin-bottom: 15px;
  padding: 10px;
  background-color: white;
  border-radius: 4px;
}

.user {
  text-align: right;
  color: #2c6bed;
  margin-bottom: 8px;
}

.ai {
  text-align: left;
  color: #333;
  margin-bottom: 8px;
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