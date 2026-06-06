<template>
  <div class="glass-panel p-4">
    <h3 class="text-sm font-mono text-purple mb-3">KILL CHAIN PROGRESS</h3>
    <div class="flex items-center gap-2 text-xs font-mono">
      <div v-for="(phase, i) in phases" :key="i" class="flex items-center gap-2">
        <span :class="phaseClass(phase.status)">{{ phase.icon }}</span>
        <span :class="phaseTextClass(phase.status)">{{ phase.name }}</span>
        <span v-if="i < phases.length - 1" class="text-gray-700 mx-1">─</span>
      </div>
    </div>
    <div class="mt-3 bg-dark rounded-full h-2 overflow-hidden">
      <div class="h-full bg-gradient-to-r from-purple to-neon rounded-full transition-all duration-500" :style="{ width: progress + '%' }"></div>
    </div>
    <span class="text-xs text-gray-600 mt-1 block">{{ progress }}% complete</span>
  </div>
</template>

<script setup>
const phases = [
  { name: 'Recon', status: 'done', icon: '✓' },
  { name: 'Weaponize', status: 'done', icon: '✓' },
  { name: 'Deliver', status: 'done', icon: '✓' },
  { name: 'Exploit', status: 'active', icon: '◉' },
  { name: 'Install', status: 'pending', icon: '□' },
  { name: 'C2', status: 'pending', icon: '□' },
  { name: 'Exfil', status: 'pending', icon: '□' },
]

const progress = 67

const phaseClass = (status) => ({
  done: 'text-neon',
  active: 'text-yellow-400 animate-pulse-slow',
  pending: 'text-gray-700',
}[status])

const phaseTextClass = (status) => ({
  done: 'text-neon',
  active: 'text-yellow-400',
  pending: 'text-gray-600',
}[status])
</script>
