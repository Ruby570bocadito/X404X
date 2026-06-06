<template>
  <div class="glass-panel p-3">
    <div class="grid grid-cols-4 gap-2">
      <div class="text-center">
        <div class="text-2xl font-mono text-neon">{{ stats.hosts }}</div>
        <div class="text-xs text-gray-500">Hosts</div>
      </div>
      <div class="text-center">
        <div class="text-2xl font-mono text-purple">{{ stats.vulns }}</div>
        <div class="text-xs text-gray-500">Vulns</div>
      </div>
      <div class="text-center">
        <div class="text-2xl font-mono text-yellow-400">{{ stats.creds }}</div>
        <div class="text-xs text-gray-500">Creds</div>
      </div>
      <div class="text-center">
        <div class="text-2xl font-mono text-neon">{{ stats.sessions }}</div>
        <div class="text-xs text-gray-500">Sessions</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAgentStore, useReconStore } from '../stores/index.js'

const agentStore = useAgentStore()
const reconStore = useReconStore()

const stats = computed(() => ({
  hosts: reconStore.hosts.length || 0,
  vulns: reconStore.vulnerabilities.length || 0,
  creds: 0,
  sessions: agentStore.activeAgents.length || 0,
}))
</script>
