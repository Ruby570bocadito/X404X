<template>
  <div class="flex flex-col h-screen scanlines">
    <Header />
    <Tabs :active="activeTab" @switch="activeTab = $event" />
    <main class="flex-1 overflow-auto p-3">
      <Dashboard v-if="activeTab === 'dashboard'" />
      <AgentPanel v-else-if="activeTab === 'agents'" />
      <NetworkMap v-else-if="activeTab === 'recon'" />
      <AIConsole v-else-if="activeTab === 'ai'" />
      <TerminalWidget v-else-if="activeTab === 'terminal'" />
      <MetricsPanel v-else-if="activeTab === 'metrics'" />
      <DocsPanel v-else />
    </main>
    <Footer />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Header from './components/Header.vue'
import Tabs from './components/Tabs.vue'
import Dashboard from './views/Dashboard.vue'
import AgentPanel from './views/AgentPanel.vue'
import NetworkMap from './views/NetworkMap.vue'
import AIConsole from './views/AIConsole.vue'
import TerminalWidget from './views/TerminalWidget.vue'
import MetricsPanel from './views/MetricsPanel.vue'
import DocsPanel from './views/DocsPanel.vue'
import Footer from './components/Footer.vue'
import { useAgentStore, useCampaignStore, useReconStore, useMetricsStore, useEventStore } from './stores/index.js'

const activeTab = ref('dashboard')

onMounted(() => {
  const agentStore = useAgentStore()
  const campaignStore = useCampaignStore()
  const reconStore = useReconStore()
  const metricsStore = useMetricsStore()
  const eventStore = useEventStore()

  // Attempt to fetch real data — falls back gracefully if API not available
  agentStore.fetchAgents().catch(() => {})
  campaignStore.fetchCampaigns().catch(() => {})
  reconStore.fetchHosts().catch(() => {})
  metricsStore.fetchMetrics().catch(() => {})

  // Connect WebSocket for live events
  eventStore.connect()
})
</script>
