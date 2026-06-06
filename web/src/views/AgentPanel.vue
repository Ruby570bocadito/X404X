<template>
  <div class="glass-panel p-4">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-sm font-mono text-purple">AGENT FLEET</h3>
      <span class="text-xs text-gray-600 font-mono">{{ store.agentCount }} agents</span>
    </div>
    <div v-if="store.loading" class="text-center text-gray-500 text-xs font-mono py-8">
      Loading agents...
    </div>
    <div v-else-if="store.error" class="text-center text-alert text-xs font-mono py-8">
      {{ store.error }}
    </div>
    <table v-else class="w-full text-xs font-mono">
      <thead>
        <tr class="text-gray-600 border-b border-gray-800">
          <th class="text-left py-2">ID</th>
          <th class="text-left">Hostname</th>
          <th class="text-left">OS</th>
          <th class="text-left">User</th>
          <th class="text-left">IP</th>
          <th class="text-left">Status</th>
          <th class="text-left">Uptime</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="agent in store.agents" :key="agent.id" class="border-b border-gray-800/50 hover:bg-panel cursor-pointer">
          <td class="py-2 text-neon">{{ agent.id }}</td>
          <td class="text-gray-300">{{ agent.hostname }}</td>
          <td class="text-gray-400">{{ agent.os }}</td>
          <td :class="agent.user === 'root' || (agent.user && agent.user.includes('SYSTEM')) ? 'text-alert' : 'text-gray-300'">{{ agent.user }}</td>
          <td class="text-gray-400">{{ agent.local_ip }}</td>
          <td>
            <span class="inline-block w-2 h-2 rounded-full mr-1" :class="statusColor(agent.status)"></span>
            <span :class="statusTextColor(agent.status)">{{ agent.status }}</span>
          </td>
          <td class="text-gray-500">{{ agent.uptime }}h</td>
          <td>
            <button v-if="agent.status !== 'dead'" @click="kill(agent)" class="btn btn-reject text-xs px-2 py-0.5">Kill</button>
          </td>
        </tr>
        <tr v-if="store.agents.length === 0">
          <td colspan="8" class="text-center text-gray-600 py-8">No agents connected. Deploy an agent to get started.</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAgentStore } from '../stores/index.js'

const store = useAgentStore()

onMounted(() => store.fetchAgents().catch(() => {}))

const kill = (agent) => {
  if (confirm(`Kill agent ${agent.id}?`)) {
    store.killAgent(agent.id, 'user requested')
  }
}

const statusColor = (s) => ({ online: 'bg-neon', active: 'bg-yellow-400', idle: 'bg-gray-500', dead: 'bg-alert' }[s] || 'bg-gray-500')
const statusTextColor = (s) => ({ online: 'text-neon', active: 'text-yellow-400', idle: 'text-gray-500', dead: 'text-alert' }[s] || 'text-gray-500')
</script>
