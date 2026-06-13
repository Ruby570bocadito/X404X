<template>
  <div class="h-full flex flex-col">
    <div class="flex items-center justify-between mb-3">
      <h2 class="text-lg font-mono text-purple">SERVICES & PORTS</h2>
      <span class="text-xs font-mono text-gray-500">{{ services.length }} services</span>
    </div>

    <div class="flex-1 overflow-auto glass-panel p-0">
      <table class="w-full text-xs font-mono">
        <thead class="sticky top-0 bg-panel">
          <tr class="text-gray-500 text-left">
            <th class="p-2">Host</th>
            <th class="p-2">Port</th>
            <th class="p-2">Protocol</th>
            <th class="p-2">Service</th>
            <th class="p-2">Banner</th>
            <th class="p-2">State</th>
            <th class="p-2">Tunnel</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(s, i) in services" :key="i" class="border-t border-gray-800 hover:bg-dark/50">
            <td class="p-2 text-gray-300">{{ s.host || s.ip || '-' }}</td>
            <td class="p-2 text-neon">{{ s.port }}</td>
            <td class="p-2 text-gray-400">{{ s.protocol || 'tcp' }}</td>
            <td class="p-2 text-gray-300">{{ s.service || s.name || '-' }}</td>
            <td class="p-2 text-gray-500 truncate max-w-[250px]" :title="s.banner">{{ s.banner || '-' }}</td>
            <td class="p-2">
              <span class="inline-block w-2 h-2 rounded-full mr-1" :class="s.state === 'open' ? 'bg-neon' : s.state === 'filtered' ? 'bg-yellow-400' : 'bg-gray-600'"></span>
              <span :class="s.state === 'open' ? 'text-neon' : s.state === 'filtered' ? 'text-yellow-400' : 'text-gray-600'">{{ s.state || 'unknown' }}</span>
            </td>
            <td class="p-2">
              <span v-if="s.tunnel_port" class="text-purple text-[10px]">{{ s.tunnel_port }}</span>
              <span v-else class="text-gray-700">-</span>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="services.length === 0" class="text-gray-600 text-center py-12 font-mono text-sm">
        <div class="text-4xl mb-3">🔌</div>
        No services discovered. Run a port scan.
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useReconStore } from '../stores/index.js'

const reconStore = useReconStore()
const services = computed(() => reconStore.services || [])

onMounted(() => reconStore.fetchServices().catch(() => {}))
</script>
