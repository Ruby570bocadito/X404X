<template>
  <div class="glass-panel p-3 h-full flex flex-col">
    <div class="flex items-center justify-between mb-2">
      <h3 class="text-sm font-mono text-purple">AI CONSOLE</h3>
      <span class="text-xs text-gray-600 font-mono">Specter</span>
    </div>

    <!-- Message log -->
    <div class="flex-1 overflow-y-auto space-y-1 text-xs font-mono mb-2" ref="chatEl">
      <div v-if="messages.length === 0" class="text-gray-600 py-2">
        Press ↻ for tactical suggestions, or type a prompt.
      </div>
      <div v-for="(msg, i) in messages" :key="i" class="leading-relaxed">
        <span :class="msg.role === 'user' ? 'text-purple' : msg.role === 'system' ? 'text-gray-600' : 'text-neon'">
          [{{ msg.role === 'user' ? 'You' : msg.role === 'system' ? 'Specter' : 'Specter' }}]
        </span>
        <span class="text-gray-300 ml-1 whitespace-pre-wrap">{{ msg.content }}</span>
      </div>
      <div v-if="loading" class="text-gray-500 italic">[AI] Thinking...</div>
    </div>

    <!-- Suggest quick action -->
    <div class="flex gap-1 mb-2">
      <button @click="suggest" class="btn text-xs flex-1" :disabled="loading">↻ suggest</button>
    </div>

    <!-- Chat input -->
    <div class="flex gap-1">
      <input v-model="input" @keyup.enter="send"
        class="flex-1 bg-dark border border-gray-800 rounded px-2 py-1 text-xs font-mono text-gray-300 focus:border-purple focus:outline-none"
        placeholder="Ask AI..." :disabled="loading" />
      <button @click="send" class="btn text-xs px-3" :disabled="!input.trim() || loading">↑</button>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { useAIStore, useCampaignStore } from '../stores/index.js'

const aiStore = useAIStore()
const campaignStore = useCampaignStore()
const input = ref('')
const loading = ref(false)
const chatEl = ref(null)
const messages = ref([
  { role: 'system', content: 'Specter online. Type a prompt or press ↻ for tactical recommendations.' }
])

const scrollBottom = async () => {
  await nextTick()
  if (chatEl.value) chatEl.value.scrollTop = chatEl.value.scrollHeight
}

const send = async () => {
  const text = input.value.trim()
  if (!text || loading.value) return
  input.value = ''
  messages.value.push({ role: 'user', content: text })
  await scrollBottom()
  loading.value = true
  try {
    const res = await fetch('/api/ai/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ prompt: text })
    })
    const data = await res.json()
    messages.value.push({ role: 'ai', content: data.response || '[no response]' })
  } catch (e) {
    messages.value.push({ role: 'ai', content: `[Error] ${e.message}` })
  } finally {
    loading.value = false
    await scrollBottom()
  }
}

const suggest = async () => {
  loading.value = true
  const cid = campaignStore.activeCampaign?.id || null
  try {
    // Try decisions endpoint first
    const params = cid ? `?campaign_id=${cid}` : ''
    const res = await fetch(`/api/decisions${params}`)
    if (res.ok) {
      const data = await res.json()
      if (Array.isArray(data) && data.length > 0) {
        const top = data.slice(0, 3)
        const text = top.map((d, i) =>
          `${i + 1}. [${(d.confidence * 100 || 0).toFixed(0)}%] ${d.technique || d.action || 'Unknown'} → ${d.target || 'any'}`
        ).join('\n')
        messages.value.push({ role: 'ai', content: text })
        await scrollBottom()
        return
      }
    }
    // Fallback: ask the AI for suggestions
    const aiRes = await fetch('/api/ai/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ prompt: 'suggest tactical recommendations for the current campaign' })
    })
    const aiData = await aiRes.json()
    messages.value.push({ role: 'ai', content: aiData.response || 'No suggestions available.' })
  } catch (e) {
    messages.value.push({ role: 'ai', content: `[Error] ${e.message}` })
  } finally {
    loading.value = false
    await scrollBottom()
  }
}
</script>
