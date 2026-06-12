<template>
  <div class="h-full flex flex-col">
    <div class="flex items-center justify-between mb-3">
      <h2 class="text-lg font-mono text-purple">CAMPAIGN TIMELINE</h2>
      <div class="flex gap-2">
        <select v-model="filter" class="bg-dark border border-gray-800 rounded px-2 py-1 text-xs font-mono text-gray-300">
          <option value="all">All Events</option>
          <option value="phase.changed">Phase Changes</option>
          <option value="agent.checkin">Agent Checkins</option>
          <option value="exploit.success">Exploits</option>
          <option value="decision.made">AI Decisions</option>
          <option value="credential.captured">Credentials</option>
          <option value="vuln.found">Vulnerabilities</option>
        </select>
        <span class="text-xs font-mono text-gray-600 self-center">
          {{ filteredEvents.length }} events
        </span>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto relative" ref="timelineEl">
      <div v-if="filteredEvents.length === 0" class="text-gray-600 text-center py-12 font-mono text-sm">
        <div class="text-4xl mb-3">⏳</div>
        No events yet. Start a campaign or deploy agents to see the timeline.
      </div>

      <div v-for="(event, i) in filteredEvents" :key="i" class="relative pl-8 pb-4">
        <!-- Vertical line -->
        <div v-if="i < filteredEvents.length - 1" class="absolute left-[11px] top-6 bottom-0 w-px"
          :class="isPhaseChange(event) ? 'bg-purple' : 'bg-gray-800'"></div>

        <!-- Event dot -->
        <div class="absolute left-1 top-1 w-5 h-5 rounded-full border-2 flex items-center justify-center text-[10px]"
          :class="dotClass(event)">

          <span v-if="isPhaseChange(event)">⬡</span>
          <span v-else-if="event.type === 'agent.checkin'">⬤</span>
          <span v-else-if="event.type === 'exploit.success'">✓</span>
          <span v-else-if="event.type === 'decision.made'">AI</span>
          <span v-else>●</span>
        </div>

        <!-- Event content -->
        <div class="glass-panel p-3 cursor-pointer"
          :class="{ 'border-l-2 border-purple': isPhaseChange(event) }"
          @click="toggleExpand(i)">
          <div class="flex items-center gap-2 mb-1">
            <span class="text-xs text-gray-600 font-mono">{{ formatTime(event.timestamp) }}</span>
            <span class="text-[10px] px-1.5 py-0.5 rounded font-mono" :class="badgeClass(event.type)">
              {{ eventTypeLabel(event.type) }}
            </span>
            <span class="text-xs text-gray-400 font-mono truncate flex-1">{{ eventTitle(event) }}</span>
          </div>

          <div v-if="expanded[i]" class="mt-2 pt-2 border-t border-gray-800 text-xs font-mono space-y-1">
            <div v-for="(value, key) in eventDetail(event)" :key="key" class="flex gap-2">
              <span class="text-gray-600 shrink-0">{{ key }}:</span>
              <span class="text-gray-300 break-all">{{ value }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useEventStore } from '../stores/index.js'

const eventStore = useEventStore()
const filter = ref('all')
const expanded = ref({})
const timelineEl = ref(null)

const toggleExpand = (i) => {
  expanded.value[i] = !expanded.value[i]
}

const filteredEvents = computed(() => {
  let events = eventStore.events || []
  if (filter.value !== 'all') {
    events = events.filter(e => e.type === filter.value)
  }
  return [...events].reverse()
})

watch(() => filteredEvents.value.length, async () => {
  await nextTick()
  if (timelineEl.value) {
    timelineEl.value.scrollTop = 0
  }
})

const formatTime = (ts) => {
  if (!ts) return '--:--:--'
  return new Date(ts).toLocaleString()
}

const isPhaseChange = (e) => e.type === 'phase.changed' || e.type === 'campaign.started'

const eventTypeLabel = (type) => {
  const map = {
    'agent.checkin': 'CHECKIN',
    'agent.dead': 'DEAD',
    'campaign.started': 'CAMPAIGN',
    'campaign.paused': 'PAUSED',
    'phase.changed': 'PHASE',
    'vuln.found': 'VULN',
    'exploit.success': 'EXPLOIT',
    'exploit.failure': 'FAIL',
    'credential.captured': 'CRED',
    'decision.made': 'AI',
    'blue.alert': 'BLUE',
    'recon.scan_complete': 'SCAN',
  }
  return map[type] || type
}

const badgeClass = (type) => ({
  'agent.checkin': 'bg-neon/10 text-neon border border-neon/30',
  'agent.dead': 'bg-red-900/20 text-red-400 border border-red-800/30',
  'campaign.started': 'bg-purple/20 text-purple border border-purple/30',
  'phase.changed': 'bg-purple/20 text-purple border border-purple/30',
  'vuln.found': 'bg-yellow-900/20 text-yellow-400 border border-yellow-800/30',
  'exploit.success': 'bg-neon/10 text-neon border border-neon/30',
  'exploit.failure': 'bg-red-900/20 text-red-400 border border-red-800/30',
  'credential.captured': 'bg-blue-900/20 text-blue-400 border border-blue-800/30',
  'decision.made': 'bg-cyan-900/20 text-cyan-400 border border-cyan-800/30',
  'blue.alert': 'bg-red-900/20 text-red-400 border border-red-800/30',
  'recon.scan_complete': 'bg-neon/10 text-neon border border-neon/30',
}[type] || 'bg-gray-800 text-gray-400 border border-gray-700')

const dotClass = (e) => {
  if (isPhaseChange(e)) return 'border-purple bg-dark text-purple'
  const map = {
    'agent.checkin': 'border-neon bg-dark text-neon',
    'agent.dead': 'border-red-500 bg-dark text-red-400',
    'exploit.success': 'border-neon bg-dark text-neon',
    'exploit.failure': 'border-red-500 bg-dark text-red-400',
    'decision.made': 'border-cyan-400 bg-dark text-cyan-400',
    'credential.captured': 'border-blue-400 bg-dark text-blue-400',
    'vuln.found': 'border-yellow-400 bg-dark text-yellow-400',
    'blue.alert': 'border-red-500 bg-dark text-red-400',
  }
  return map[e.type] || 'border-gray-600 bg-dark text-gray-400'
}

const eventTitle = (e) => {
  const d = e.data || {}
  switch (e.type) {
    case 'agent.checkin':
      return `Agent ${e.agent_id || d.agent_id || '?'} checked in — ${d.hostname || 'unknown'} (${d.os || '?'})`
    case 'agent.dead':
      return `Agent ${e.agent_id || ''} marked dead — ${d.reason || 'connection lost'}`
    case 'campaign.started':
      return `Campaign "${d.name || '?'}" started — target: ${d.target_scope || 'local'}`
    case 'campaign.paused':
      return `Campaign paused`
    case 'phase.changed':
      return `Phase → ${d.to || d.phase || '?'} (${Math.round((d.progress || 0) * 100)}%)`
    case 'vuln.found':
      return `${d.cve || 'CVE'} [${d.severity || '?'}] on ${d.target_ip || 'host'} — ${d.service || ''}`
    case 'exploit.success':
      return `${d.exploit || d.technique || 'Exploit'} succeeded on ${d.target || 'host'}`
    case 'exploit.failure':
      return `Exploit failed on ${d.target || 'host'}: ${d.error || 'unknown'}`
    case 'credential.captured':
      return `${d.username || '?'}@${d.domain || d.source || 'host'} captured`
    case 'decision.made':
      return `AI decision: ${d.tactic || '?'} → ${d.technique || '?'} [${d.mitre_id || 'T????'}] (${Math.round((d.confidence || 0) * 100)}%)`
    case 'blue.alert':
      return `Blue detection: ${d.tool || '?'} — ${d.alert_type || '?'}`
    case 'recon.scan_complete':
      return `Recon complete: ${d.hosts_found || 0} hosts, ${d.vulns_found || 0} vulns on ${d.target || '?'}`
    default:
      return e.type
  }
}

const eventDetail = (e) => {
  const d = e.data || {}
  const detail = {}
  if (e.agent_id) detail.agent_id = e.agent_id
  Object.entries(d).forEach(([k, v]) => {
    if (v !== null && v !== undefined && v !== '' && typeof v !== 'object') {
      detail[k] = String(v)
    }
  })
  if (Object.keys(detail).length === 0) {
    detail.type = e.type
    detail.raw = JSON.stringify(e)
  }
  return detail
}
</script>
