<template>
  <nav class="flex gap-1 px-3 pb-2">
    <button v-for="tab in tabs" :key="tab.id" @click="$emit('switch', tab.id)"
      class="px-4 py-2 text-sm font-mono rounded-t-lg transition-all"
      :class="tab.id === active ? 'text-neon border-b-2 border-neon bg-panel' : 'text-gray-500 hover:text-gray-300'">
      {{ tab.label }}
      <span v-if="tab.badge !== null" class="ml-1 text-xs text-neon">● {{ tab.badge }}</span>
    </button>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useAgentStore, useReconStore } from '../stores/index.js'

const props = defineProps({ active: { type: String, default: 'dashboard' } })
defineEmits(['switch'])

const agentStore = useAgentStore()
const reconStore = useReconStore()

const tabs = computed(() => [
  { id: 'dashboard', label: 'Dashboard', badge: null },
  { id: 'agents', label: 'Agents', badge: agentStore.activeAgents.length },
  { id: 'recon', label: 'Recon', badge: reconStore.hosts.length },
  { id: 'builder', label: 'Payloads', badge: null },
  { id: 'ai', label: 'AI', badge: null },
  { id: 'browser', label: 'Browser', badge: null },
  { id: 'terminal', label: 'Terminal', badge: null },
  { id: 'metrics', label: 'Metrics', badge: null },
  { id: 'docs', label: 'Docs', badge: null },
])
</script>
