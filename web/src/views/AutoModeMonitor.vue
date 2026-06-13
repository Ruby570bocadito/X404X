<template>
  <div class="h-full flex flex-col">
    <div class="flex items-center justify-between mb-3">
      <h2 class="text-lg font-mono text-purple">AUTOMODE MONITOR</h2>
      <div class="flex items-center gap-2">
        <span class="text-xs font-mono text-gray-500">Auto-approve</span>
        <button @click="toggleAuto" class="w-8 h-4 rounded-full border transition-colors"
          :class="autoApprove ? 'bg-neon/30 border-neon' : 'bg-gray-800 border-gray-700'">
          <span class="block w-3 h-3 rounded-full transition-transform"
            :class="autoApprove ? 'bg-neon translate-x-4' : 'bg-gray-500 translate-x-0.5'"></span>
        </button>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto space-y-1">
      <div v-if="activities.length === 0" class="text-gray-600 text-center py-12 font-mono text-sm">
        <div class="text-4xl mb-3">⚡</div>
        No autonomous activity yet. Enable auto-mode or deploy agents.
      </div>
      <div v-for="(a, i) in activities" :key="i" class="text-xs font-mono leading-relaxed flex gap-2 py-0.5">
        <span class="text-gray-600 w-16 shrink-0">{{ formatTime(a.timestamp) }}</span>
        <span class="px-1 rounded text-[10px]" :class="activityBadge(a.type)">{{ a.type }}</span>
        <span class="text-gray-300">{{ a.description }}</span>
        <span v-if="a.agent_id" class="text-purple ml-auto text-[10px]">{{ a.agent_id }}</span>
      </div>
    </div>

    <div class="text-[10px] font-mono text-gray-600 mt-2 flex items-center justify-between">
      <span>{{ activities.length }} events</span>
      <button @click="clearActivities" class="hover:text-gray-400">Clear</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const activities = ref([])
const autoApprove = ref(false)

onMounted(async () => {
  try {
    const r = await fetch('/api/activities')
    if (r.ok) activities.value = await r.json()
  } catch {}
  try {
    const r = await fetch('/api/config')
    if (r.ok) {
      const cfg = await r.json()
      autoApprove.value = cfg.auto_approve ?? false
    }
  } catch {}
})

function toggleAuto() {
  autoApprove.value = !autoApprove.value
  fetch('/api/config', {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ auto_approve: autoApprove.value })
  }).catch(() => {})
}

function formatTime(ts) {
  if (!ts) return '--:--:--'
  const d = new Date(ts)
  return d.toLocaleTimeString('en-US', { hour12: false })
}

function clearActivities() { activities.value = [] }

function activityBadge(t) {
  const m = {
    'decision': 'bg-purple/20 text-purple border border-purple/30',
    'scan': 'bg-cyan-900/20 text-cyan-400 border border-cyan-800/30',
    'exploit': 'bg-red-900/20 text-red-400 border border-red-800/30',
    'lateral': 'bg-pink-900/20 text-pink-400 border border-pink-800/30',
    'privesc': 'bg-orange-900/20 text-orange-400 border border-orange-800/30',
    'credential': 'bg-yellow-900/20 text-yellow-400 border border-yellow-800/30',
    'error': 'bg-red-900/30 text-red-300 border border-red-700/30',
    'info': 'bg-gray-800 text-gray-400 border border-gray-700',
  }
  return m[t] || m['info']
}
</script>
