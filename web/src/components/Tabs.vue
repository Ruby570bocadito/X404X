<template>
  <nav class="flex gap-1 px-3 pb-2">
    <button
      v-for="tab in tabs"
      :key="tab.id"
      @click="$emit('switch', tab.id)"
      class="px-4 py-2 text-sm font-mono rounded-t-lg transition-all"
      :class="tab.id === active ? 'text-neon border-b-2 border-neon bg-panel' : 'text-gray-500 hover:text-gray-300'"
    >
      {{ tab.label }}
      <span v-if="tab.badge" class="ml-1 text-xs text-neon">● {{ tab.badge }}</span>
    </button>
  </nav>
</template>

<script setup>
import { inject } from 'vue'

const props = defineProps({
  active: { type: String, default: 'dashboard' },
})

const emit = defineEmits(['switch'])

const activeTab = inject('activeTab')
const switchTab = inject('switchTab')

const tabs = [
  { id: 'dashboard', label: 'Dashboard' },
  { id: 'agents', label: 'Agents', badge: '5' },
  { id: 'recon', label: 'Recon' },
  { id: 'ai', label: 'AI' },
  { id: 'terminal', label: 'Terminal' },
  { id: 'metrics', label: 'Metrics' },
  { id: 'docs', label: 'Docs' },
]
</script>
