<template>
  <div class="glass-panel p-3 h-full flex flex-col">
    <h3 class="text-sm font-mono text-purple mb-2">AI CONSOLE — Specter + Apex (Ollama)</h3>
    <div class="flex-1 overflow-y-auto font-mono text-xs space-y-2 mb-3" ref="aiChat">
      <div v-for="(msg, i) in messages" :key="i" :class="msg.role === 'user' ? 'text-purple' : 'text-neon'">
        <span class="text-gray-600">[{{ msg.role === 'user' ? 'You' : 'Specter' }}]</span>
        {{ msg.content }}
      </div>
      <div v-if="store.loading" class="text-gray-500 italic">[AI] Thinking...</div>
    </div>
    <div class="flex gap-2">
      <input v-model="prompt" @keyup.enter="send"
        class="flex-1 bg-dark border border-gray-800 rounded px-3 py-2 text-xs font-mono text-gray-300 focus:border-purple focus:outline-none"
        placeholder="Ask the AI assistant..." :disabled="store.loading" />
      <button @click="send" class="btn" :disabled="!prompt.trim() || store.loading">Send</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAIStore } from '../stores/index.js'

const store = useAIStore()
const prompt = ref('')
const messages = store.messages

// Seed with welcome message if empty
if (messages.length === 0) {
  messages.push({ role: 'system', content: 'Connected to Specter-Terminal. Model: llama3.2 | Ollama: localhost:11434' })
  messages.push({ role: 'ai', content: 'X404X AI ready. Type a prompt or use "suggest" for tactical recommendations.' })
}

const send = async () => {
  if (!prompt.value.trim() || store.loading) return
  await store.chat(prompt.value)
  prompt.value = ''
}
</script>
