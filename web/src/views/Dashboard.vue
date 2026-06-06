<template>
  <div class="grid grid-cols-3 gap-3 h-full">
    <div class="col-span-2 space-y-3">
      <KillChain />
      <div class="grid grid-cols-2 gap-3">
        <ActivityFeed />
        <SessionList />
      </div>
    </div>
    <div class="space-y-3">
      <AIChatMini />
      <StatCards />
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import KillChain from '../components/KillChain.vue'
import ActivityFeed from '../components/ActivityFeed.vue'
import SessionList from '../components/SessionList.vue'
import AIChatMini from '../components/AIChatMini.vue'
import StatCards from '../components/StatCards.vue'
import { useAgentStore, useReconStore } from '../stores/index.js'

onMounted(() => {
  useAgentStore().fetchAgents().catch(() => {})
  useReconStore().fetchHosts().catch(() => {})
})
</script>
