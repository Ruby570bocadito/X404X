<template>
  <div class="glass-panel p-3 h-48 flex flex-col">
    <h3 class="text-sm font-mono text-purple mb-2">ACTIVE SESSIONS</h3>
    <div v-if="agentStore.loading" class="text-center text-gray-600 text-xs py-4">Loading...</div>
    <div v-else-if="agentStore.activeAgents.length === 0" class="text-center text-gray-600 text-xs py-4">
      No active sessions. Deploy an agent.
    </div>
    <div v-else class="flex-1 overflow-y-auto text-xs font-mono">
      <div class="grid grid-cols-4 text-gray-600 border-b border-gray-800 pb-1 mb-1 sticky top-0 bg-dark bg-opacity-90">
        <span>ID</span><span>Target</span><span>OS</span><span>User</span>
      </div>
      <div v-for="a in agentStore.activeAgents" :key="a.id"
        class="grid grid-cols-4 py-1 hover:bg-panel rounded cursor-pointer transition-colors">
        <span class="text-neon">{{ a.session_id || a.id?.slice(0, 6) }}</span>
        <span class="text-gray-300">{{ a.local_ip || '—' }}</span>
        <span class="text-gray-400">{{ (a.os || '?').slice(0, 14) }}</span>
        <span :class="isPrivileged(a.username) ? 'text-alert' : 'text-gray-300'">{{ a.username || '?' }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAgentStore } from '../stores/index.js'
const agentStore = useAgentStore()
const isPrivileged = (u) => u && (u.includes('SYSTEM') || u === 'root')
</script>
