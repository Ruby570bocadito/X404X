<template>
  <div class="glass-panel p-4">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-sm font-mono text-purple">BLUEFORGE METRICS</h3>
      <button @click="refresh" class="text-xs text-gray-600 hover:text-purple" :disabled="loading">↻ refresh</button>
    </div>
    <div class="grid grid-cols-2 gap-3" v-if="!loading">
      <div class="glass-panel p-3 text-center">
        <div class="text-3xl font-mono text-neon">{{ Math.round(metrics.stealth_rating * 100) || 0 }}%</div>
        <div class="text-xs text-gray-500 mt-1">Evasion Rate</div>
      </div>
      <div class="glass-panel p-3 text-center">
        <div class="text-3xl font-mono text-alert">{{ blueMetrics.filter(m => m.detected).length }}</div>
        <div class="text-xs text-gray-500 mt-1">Detections</div>
      </div>
    </div>
    <div v-else class="text-center text-gray-600 py-4 text-xs">Loading metrics...</div>
    <div class="mt-3 text-xs font-mono" v-if="blueMetrics.length > 0">
      <div class="text-gray-500 mb-2">Detection Log:</div>
      <div class="space-y-1">
        <div v-for="(m, i) in blueMetrics" :key="i" class="flex gap-2">
          <span :class="m.detected ? 'text-purple' : 'text-neon'">[{{ m.tool || '?' }}]</span>
          <span :class="m.detected ? 'text-gray-300' : 'text-gray-400'">
            {{ m.detected ? m.alert_type : 'Bypassed: ' + (m.tool || 'evasion') }}
          </span>
          <span class="text-gray-600 ml-auto text-xs">{{ formatTime(m.timestamp) }}</span>
        </div>
      </div>
    </div>
    <div v-else class="text-gray-600 text-xs mt-2">No detection data yet. Deploy agents to collect metrics.</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMetricsStore } from '../stores/index.js'

const metricsStore = useMetricsStore()
const loading = ref(false)
const metrics = ref({})
const blueMetrics = ref([])

const refresh = async () => {
  loading.value = true
  await Promise.all([
    metricsStore.fetchMetrics().then(() => { metrics.value = metricsStore.metrics }),
    metricsStore.fetchBlueMetrics().then(() => { blueMetrics.value = metricsStore.blueMetrics })
  ]).catch(() => {})
  loading.value = false
}

const formatTime = (ts) => ts ? new Date(ts).toLocaleTimeString() : '--:--'

onMounted(refresh)
</script>
