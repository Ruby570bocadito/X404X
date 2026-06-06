<template>
  <div class="flex gap-3 h-full">
    <div class="glass-panel flex-1 p-4 relative" ref="mapContainer">
      <h3 class="text-sm font-mono text-purple mb-3">NETWORK TOPOLOGY</h3>
      <svg width="100%" height="100%" class="overflow-visible">
        <g v-for="edge in edges" :key="edge.from + edge.to">
          <line :x1="nodes[edge.from].x" :y1="nodes[edge.from].y" :x2="nodes[edge.to].x" :y2="nodes[edge.to].y"
            stroke="rgba(108,99,255,0.2)" stroke-width="1" stroke-dasharray="4" />
        </g>
        <g v-for="(node, id) in nodes" :key="id">
          <circle :cx="node.x" :cy="node.y" r="20" :fill="nodeColor(node.status)" stroke="rgba(255,255,255,0.1)" />
          <text :x="node.x" :y="node.y + 5" text-anchor="middle" fill="white" font-size="8" font-family="monospace">{{ id }}</text>
          <text :x="node.x" :y="node.y + 32" text-anchor="middle" :fill="nodeTextColor(node.status)" font-size="9" font-family="monospace">{{ node.label }}</text>
        </g>
      </svg>
    </div>
    <div class="glass-panel w-80 p-3 space-y-2 text-xs font-mono">
      <h3 class="text-purple">HOST DETAILS</h3>
      <div v-for="(node, id) in nodes" :key="id" class="flex justify-between py-1 border-b border-gray-800/50">
        <span class="text-gray-300">{{ id }}</span>
        <span class="text-gray-500">{{ node.label }}</span>
        <span :class="statusInfo(node.status)">{{ node.status }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
const nodes = {
  'DC': { x: 200, y: 80, label: 'Win2019', status: 'compromised' },
  'DB': { x: 120, y: 220, label: 'Ubuntu24', status: 'scanned' },
  'WS1': { x: 300, y: 200, label: 'Win11', status: 'compromised' },
  'WS2': { x: 380, y: 140, label: 'Win11', status: 'discovered' },
  'WEB': { x: 250, y: 300, label: 'CentOS8', status: 'discovered' },
}

const edges = [
  { from: 'DC', to: 'WS1' },
  { from: 'DC', to: 'WS2' },
  { from: 'DC', to: 'DB' },
  { from: 'WS1', to: 'WEB' },
]

const nodeColor = (s) => ({ compromised: '#ff4444', scanned: '#6c63ff', discovered: '#333' }[s])
const nodeTextColor = (s) => ({ compromised: '#ff4444', scanned: '#6c63ff', discovered: '#666' }[s])
const statusInfo = (s) => ({ compromised: 'text-alert', scanned: 'text-purple', discovered: 'text-gray-600' }[s])
</script>
