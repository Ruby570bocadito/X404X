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
      <div class="h-full bg-gradient-to-r from-purple to-neon rounded-full transition-all duration-500"
        :style="{ width: (campaignStore.activeCampaign?.progress || 0) * 100 + '%' }"></div>
    </div>
    <span class="text-xs text-gray-600 mt-1 block">
      {{ Math.round((campaignStore.activeCampaign?.progress || 0) * 100) }}% complete
      <span v-if="campaignStore.activeCampaign" class="text-purple ml-2">{{ campaignStore.activeCampaign.phase }}</span>
    </span>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useCampaignStore } from '../stores/index.js'

const campaignStore = useCampaignStore()

onMounted(() => campaignStore.fetchCampaigns().catch(() => {}))

const phases = computed(() => {
  const currentPhase = campaignStore.activeCampaign?.phase || 'recon'
  const phaseOrder = { recon: 0, weaponization: 1, delivery: 2, exploitation: 3, installation: 4, c2: 5, actions_on_objective: 6, exfiltration: 7 }
  const currentIdx = phaseOrder[currentPhase] ?? 0

  const allPhases = ['Recon', 'Weaponize', 'Deliver', 'Exploit', 'Install', 'C2', 'Exfil']
  return allPhases.map((name, i) => ({
    name,
    icon: i < currentIdx ? '✓' : i === currentIdx ? '◉' : '□',
    status: i < currentIdx ? 'done' : i === currentIdx ? 'active' : 'pending'
  }))
})

const phaseClass = (s) => ({ done: 'text-neon', active: 'text-yellow-400 animate-pulse-slow', pending: 'text-gray-700' }[s])
const phaseTextClass = (s) => ({ done: 'text-neon', active: 'text-yellow-400', pending: 'text-gray-600' }[s])
</script>
