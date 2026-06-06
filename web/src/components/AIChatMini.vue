<template>
  <div class="glass-panel p-3 h-full flex flex-col">
    <div class="flex items-center justify-between mb-2">
      <h3 class="text-sm font-mono text-purple">AI CONSOLE</h3>
      <button @click="refresh" class="text-xs text-gray-600 hover:text-purple" :disabled="aiStore.loading">
        {{ aiStore.loading ? '...' : '↻ suggest' }}
      </button>
    </div>
    <div class="flex-1 overflow-y-auto space-y-2 text-xs font-mono mb-3">
      <div v-if="aiStore.suggestions.length === 0 && !aiStore.loading" class="text-gray-600">
        Press ↻ for AI attack suggestions
      </div>
      <div v-for="(s, i) in aiStore.suggestions.slice(0, 6)" :key="i"
        class="bg-panel rounded p-2 space-y-1">
        <span :class="s.source === 'ai' ? 'text-purple' : 'text-neon'">{{ i + 1 }}. {{ s.technique }}</span>
        <span class="text-gray-500 text-xs">[{{ s.tactic }}] target: {{ s.target || 'any' }} | conf: {{ (s.confidence * 100).toFixed(0) }}%</span>
      </div>
      <div v-if="aiStore.loading" class="text-gray-500">[AI] Analyzing...</div>
    </div>
    <div class="flex gap-2">
      <button @click="approve(0)" class="btn btn-accept flex-1 text-xs" :disabled="aiStore.suggestions.length === 0">
        Accept #1
      </button>
      <button @click="reject(0)" class="btn btn-reject flex-1 text-xs" :disabled="aiStore.suggestions.length === 0">
        Reject
      </button>
    </div>
  </div>
</template>

<script setup>
import { useAIStore, useCampaignStore } from '../stores/index.js'

const aiStore = useAIStore()
const campaignStore = useCampaignStore()

const refresh = () => {
  const cid = campaignStore.activeCampaign?.id || null
  aiStore.getSuggestions(cid).catch(() => {})
}

const approve = (idx) => {
  const d = aiStore.suggestions[idx]
  if (d) aiStore.approveDecision(d.id).then(() => refresh())
}
const reject = (idx) => {
  const d = aiStore.suggestions[idx]
  if (d) aiStore.rejectDecision(d.id)
}
</script>
