<template>
  <div class="flex gap-3 h-full">
    <div class="flex-1 space-y-3">
      <!-- Browser Mesh Map -->
      <div class="glass-panel p-4" ref="meshContainer">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-sm font-mono text-purple">BROWSER MESH NETWORK</h3>
          <span class="text-xs text-gray-600 font-mono">{{ meshNodes.length }} nodes · WebRTC P2P</span>
        </div>
        <div class="grid grid-cols-4 gap-3">
          <div v-for="node in meshNodes" :key="node.id"
            class="glass-panel p-3 text-center hover:border-purple transition-all cursor-pointer">
            <div class="text-neon text-lg mb-1">{{ node.browser?.split(' ')[0] || '?' }}</div>
            <div class="text-xs text-gray-400 font-mono">{{ node.os }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ node.url }}</div>
            <div class="flex justify-center gap-2 mt-2">
              <span v-if="node.sw_persisted" class="text-neon text-xs">⚡ SW</span>
              <span class="inline-block w-2 h-2 rounded-full bg-neon animate-pulse-slow"></span>
              <span class="text-neon text-xs">{{ node.status }}</span>
            </div>
          </div>
          <div v-if="meshNodes.length === 0"
            class="col-span-4 text-center text-gray-600 py-8 text-xs font-mono">
            No browser implants deployed. Use "exploit/phantom_xss" to infect a target.
          </div>
        </div>
      </div>

      <!-- Stats Row -->
      <div class="grid grid-cols-5 gap-2">
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-neon">{{ stats.totalNodes }}</div>
          <div class="text-xs text-gray-500">Implants</div>
        </div>
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-purple">{{ stats.activeNodes }}</div>
          <div class="text-xs text-gray-500">Active</div>
        </div>
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-yellow-400">{{ stats.cookiesTotal }}</div>
          <div class="text-xs text-gray-500">Cookies</div>
        </div>
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-neon">{{ stats.sessionsTotal }}</div>
          <div class="text-xs text-gray-500">Sessions</div>
        </div>
        <div class="glass-panel p-2 text-center">
          <div class="text-lg font-mono text-purple">{{ stats.swPersisted }}</div>
          <div class="text-xs text-gray-500">SW Active</div>
        </div>
      </div>
    </div>

    <!-- Actions Panel -->
    <div class="glass-panel w-72 p-3 space-y-2 text-xs font-mono overflow-y-auto">
      <h3 class="text-purple mb-2">PHANTOM ACTIONS</h3>
      <button @click="action('inject_xss')" class="btn w-full text-left">🎯 Inject XSS Implant</button>
      <button @click="action('watering_hole')" class="btn w-full text-left">💧 Watering Hole Deploy</button>
      <button @click="action('sw_persist')" class="btn w-full text-left">⚡ Service Worker Persist</button>
      <button @click="action('steal_cookies')" class="btn w-full text-left">🍪 Steal Cookies</button>
      <button @click="action('screenshot')" class="btn w-full text-left">📸 Capture Screenshot</button>
      <button @click="action('keylogger')" class="btn w-full text-left">⌨️ Start Keylogger</button>
      <button @click="action('socks5')" class="btn w-full text-left">🔌 Enable SOCKS5</button>
      <div class="border-t border-gray-800 pt-2 mt-2">
        <button @click="refresh" class="btn w-full text-center text-gray-500">↻ Refresh Status</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const meshNodes = ref([
  { id: 'pw-a1b2c3d4', url: 'https://corp.local', browser: 'Chrome 125', os: 'Windows 11', status: 'active', sw_persisted: true },
  { id: 'pw-e5f6g7h8', url: 'https://mail.corp.local', browser: 'Firefox 127', os: 'Ubuntu 24.04', status: 'active', sw_persisted: false },
])

const stats = computed(() => ({
  totalNodes: meshNodes.value.length,
  activeNodes: meshNodes.value.filter(n => n.status === 'active').length,
  cookiesTotal: 6,
  sessionsTotal: 3,
  swPersisted: meshNodes.value.filter(n => n.sw_persisted).length,
}))

const action = (type) => {
  // In production: calls PhantomWeb bridge via API
  console.log(`PhantomWeb action: ${type}`)
}

const refresh = () => {
  // In production: fetch from API
}

onMounted(refresh)
</script>
