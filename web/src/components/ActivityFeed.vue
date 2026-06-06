<template>
  <div class="glass-panel p-3 h-48 flex flex-col">
    <div class="flex items-center justify-between mb-2">
      <h3 class="text-sm font-mono text-purple">LIVE FEED</h3>
      <span :class="eventStore.connected ? 'text-neon' : 'text-gray-600'" class="text-xs">● {{ eventStore.connected ? 'WS connected' : 'offline' }}</span>
    </div>
    <div class="flex-1 overflow-y-auto space-y-1 text-xs font-mono" ref="feedEl">
      <div v-for="(event, i) in displayEvents" :key="i" class="flex gap-2">
        <span class="text-gray-600 shrink-0">{{ formatTime(event.timestamp) }}</span>
        <span :class="eventTypeClass(event.type)">[{{ event.typeShort }}]</span>
        <span class="text-gray-300">{{ event.message || event.data?.detail || JSON.stringify(event.data) }}</span>
      </div>
      <div v-if="displayEvents.length === 0" class="text-gray-600 py-4 text-center">
        Awaiting events... Deploy agents to see live activity.
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useEventStore } from '../stores/index.js'

const eventStore = useEventStore()

const displayEvents = computed(() => {
  return eventStore.events.slice(0, 20).map(e => ({
    ...e,
    typeShort: eventShortType(e.type),
    message: formatEventMessage(e),
  }))
})

const formatTime = (ts) => {
  if (!ts) return '--:--:--'
  return new Date(ts).toLocaleTimeString()
}

const eventShortType = (type) => {
  const map = {
    'agent.checkin': '+', 'agent.dead': '!',
    'vuln.found': '*', 'exploit.success': '✓',
    'exploit.failure': '✗', 'credential.captured': '$',
    'phase.changed': '→', 'blue.alert': '!',
    'campaign.started': '▶', 'decision.made': 'AI',
  }
  return map[type] || '*'
}

const eventTypeClass = (type) => ({
  'agent.checkin': 'text-neon', 'agent.dead': 'text-alert',
  'vuln.found': 'text-purple', 'exploit.success': 'text-neon',
  'exploit.failure': 'text-alert', 'credential.captured': 'text-yellow-400',
  'phase.changed': 'text-purple', 'blue.alert': 'text-alert',
  'campaign.started': 'text-neon', 'decision.made': 'text-blue-400',
}[type] || 'text-gray-500')

const formatEventMessage = (e) => {
  if (!e) return ''
  const d = e.data || {}
  switch (e.type) {
    case 'agent.checkin': return `Agent ${e.agent_id || d.agent_id || '?'} checked in`
    case 'agent.dead': return `Agent ${e.agent_id || ''} lost`
    case 'vuln.found': return `${d.cve || d.description || 'Vulnerability'} on ${d.target_ip || d.service || 'host'}`
    case 'exploit.success': return `${d.vector || d.technique || 'Exploit'} succeeded on ${d.target || ''}`
    case 'credential.captured': return `${d.username || ''}@${d.domain || d.source || 'host'}`
    case 'phase.changed': return `Campaign advanced to: ${d}`
    default: return e.type
  }
}
</script>
