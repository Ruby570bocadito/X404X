<template>
  <div class="flex gap-3 h-full">
    <div class="glass-panel flex-1 p-4 relative" ref="mapContainer">
      <h3 class="text-sm font-mono text-purple mb-3">NETWORK TOPOLOGY</h3>
      <svg width="100%" height="100%" viewBox="0 0 500 350">
        <g v-for="(edge, ei) in edges" :key="'e'+ei">
          <line :x1="graphNodes[edge.from]?.x || 0" :y1="graphNodes[edge.from]?.y || 0"
            :x2="graphNodes[edge.to]?.x || 0" :y2="graphNodes[edge.to]?.y || 0"
            stroke="rgba(108,99,255,0.2)" stroke-width="1.5" stroke-dasharray="4" />
        </g>
        <g v-for="(node, ip) in graphNodes" :key="ip">
          <circle :cx="node.x" :cy="node.y" r="22"
            :fill="nodeColor(node.status)"
            stroke="rgba(255,255,255,0.15)" stroke-width="1.5" />
          <text :x="node.x" :y="node.y + 4" text-anchor="middle" fill="white"
            font-size="8" font-family="monospace" font-weight="bold">{{ node.label }}</text>
          <text :x="node.x" :y="node.y + 32" text-anchor="middle"
            :fill="nodeTextColor(node.status)" font-size="9" font-family="monospace">
            {{ node.os?.slice(0, 11) || '?' }}
          </text>
        </g>
      </svg>
    </div>
    <div class="glass-panel w-80 p-3 space-y-2 text-xs font-mono overflow-y-auto">
      <h3 class="text-purple">HOST DETAILS</h3>
      <div v-if="reconStore.loading" class="text-gray-600 py-4 text-center">Loading...</div>
      <div v-else-if="reconStore.hosts.length === 0" class="text-gray-600 py-4 text-center">
        No hosts discovered. Run recon scan.
      </div>
      <div v-for="h in reconStore.hosts" :key="h.ip"
        class="flex justify-between items-center py-1 border-b border-gray-800/50 hover:bg-panel rounded px-1">
        <span class="text-gray-300 w-28 truncate">{{ h.ip }}</span>
        <span class="text-gray-500 w-16 truncate">{{ h.hostname || '?' }}</span>
        <span :class="h.asset_value >= 70 ? 'text-alert' : h.asset_value >= 30 ? 'text-purple' : 'text-gray-500'">
          {{ h.asset_value || 0 }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useReconStore, useAgentStore } from '../stores/index.js'

const reconStore = useReconStore()
const agentStore = useAgentStore()

onMounted(() => reconStore.fetchHosts().catch(() => {}))

const positions = [{ x: 250, y: 60 }, { x: 120, y: 200 }, { x: 380, y: 180 }, { x: 250, y: 280 }, { x: 420, y: 100 }]

const graphNodes = computed(() => {
  const nodes = {}
  reconStore.hosts.forEach((h, i) => {
    const isCompromised = agentStore.agents.some(a => a.local_ip === h.ip && (a.status === 'online' || a.status === 'active'))
    nodes[h.ip] = {
      x: positions[i % positions.length]?.x || 100 + i * 80,
      y: positions[i % positions.length]?.y || 100 + i * 60,
      label: h.hostname || h.ip?.slice(-4),
      os: h.os,
      status: isCompromised ? 'compromised' : 'discovered'
    }
  })
  return nodes
})

const edges = computed(() => {
  const ips = Object.keys(graphNodes.value)
  const result = []
  for (let i = 0; i < ips.length - 1; i++) {
    result.push({ from: ips[i], to: ips[i + 1] })
  }
  return result
})

const nodeColor = (s) => ({ compromised: '#ff4444', discovered: '#6c63ff' }[s] || '#333')
const nodeTextColor = (s) => ({ compromised: '#ff4444', discovered: '#6c63ff' }[s] || '#666')
</script>
