<template>
  <div class="flex flex-col h-screen scanlines">
    <Header />
    <Tabs :active="activeTab" @switch="activeTab = $event" />
    <main class="flex-1 overflow-auto p-3">
      <Dashboard v-if="activeTab === 'dashboard'" />
      <CampaignManager v-else-if="activeTab === 'campaign'" />
      <DecisionsPanel v-else-if="activeTab === 'decisions'" />
      <AgentPanel v-else-if="activeTab === 'agents'" />
      <NetworkMap v-else-if="activeTab === 'recon'" />
      <VulnerabilityHeatmap v-else-if="activeTab === 'vulns'" />
      <ServicePortTable v-else-if="activeTab === 'services'" />
      <PayloadBuilder v-else-if="activeTab === 'builder'" />
      <AutoModeMonitor v-else-if="activeTab === 'automode'" />
      <AIConsole v-else-if="activeTab === 'ai'" />
      <TerminalWidget v-else-if="activeTab === 'terminal'" />
      <MetricsPanel v-else-if="activeTab === 'metrics'" />
      <BrowserMesh v-else-if="activeTab === 'browser'" />
      <CredentialVault v-else-if="activeTab === 'creds'" />
      <CampaignTimeline v-else-if="activeTab === 'timeline'" />
      <DocsPanel v-else />
    </main>
    <Footer />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import Header from './components/Header.vue'
import Tabs from './components/Tabs.vue'
import Dashboard from './views/Dashboard.vue'
import CampaignManager from './views/CampaignManager.vue'
import DecisionsPanel from './views/DecisionsPanel.vue'
import AgentPanel from './views/AgentPanel.vue'
import NetworkMap from './views/NetworkMap.vue'
import VulnerabilityHeatmap from './views/VulnerabilityHeatmap.vue'
import ServicePortTable from './views/ServicePortTable.vue'
import AutoModeMonitor from './views/AutoModeMonitor.vue'
import AIConsole from './views/AIConsole.vue'
import PayloadBuilder from './views/PayloadBuilder.vue'
import TerminalWidget from './views/TerminalWidget.vue'
import MetricsPanel from './views/MetricsPanel.vue'
import DocsPanel from './views/DocsPanel.vue'
import BrowserMesh from './views/BrowserMesh.vue'
import CredentialVault from './views/CredentialVault.vue'
import CampaignTimeline from './views/CampaignTimeline.vue'
import Footer from './components/Footer.vue'
import {
  useAgentStore, useCampaignStore, useReconStore,
  useMetricsStore, useEventStore, useAIStore
} from './stores/index.js'

const activeTab = ref('dashboard')

let pollTimer = null

onMounted(() => {
  const agentStore = useAgentStore()
  const campaignStore = useCampaignStore()
  const reconStore = useReconStore()
  const metricsStore = useMetricsStore()
  const eventStore = useEventStore()
  const aiStore = useAIStore()

  // Initial fetch — all stores query the API
  agentStore.fetchAgents().catch(() => {})
  campaignStore.fetchCampaigns().catch(() => {})
  reconStore.fetchHosts().catch(() => {})
  reconStore.fetchServices().catch(() => {})
  reconStore.fetchVulnerabilities().catch(() => {})
  metricsStore.fetchMetrics().catch(() => {})
  metricsStore.fetchBlueMetrics().catch(() => {})

  // AI suggestions (async, non-blocking)
  setTimeout(() => {
    const cid = campaignStore.activeCampaign?.id
    if (cid) aiStore.getSuggestions(cid).catch(() => {})
  }, 2000)

  // Poll every 10 seconds for live updates
  pollTimer = setInterval(() => {
    agentStore.fetchAgents().catch(() => {})
    reconStore.fetchHosts().catch(() => {})
  }, 10000)

  // Connect WebSocket for real-time events
  eventStore.connect()
})

onUnmounted(() => {
  clearInterval(pollTimer)
  useEventStore().disconnect()
})
</script>
