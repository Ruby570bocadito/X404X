<template>
  <div class="h-full flex flex-col">
    <div class="flex items-center justify-between mb-3">
      <h2 class="text-lg font-mono text-purple">CAMPAIGN MANAGER</h2>
      <button @click="showWizard = true" class="btn text-xs px-3 py-1">+ New Campaign</button>
    </div>

    <!-- Campaign Wizard Modal -->
    <div v-if="showWizard" class="fixed inset-0 z-50 flex items-center justify-center bg-black/70" @click.self="showWizard = false">
      <div class="glass-panel p-6 w-[500px] max-h-[80vh] overflow-y-auto">
        <h3 class="text-lg font-mono text-purple mb-4">NEW CAMPAIGN</h3>
        <div class="space-y-3">
          <div>
            <label class="text-xs font-mono text-gray-400 block mb-1">Name</label>
            <input v-model="form.name" class="w-full bg-dark border border-gray-800 rounded px-3 py-2 text-sm font-mono text-gray-200" placeholder="Operation Name">
          </div>
          <div>
            <label class="text-xs font-mono text-gray-400 block mb-1">Target Scope</label>
            <input v-model="form.targetScope" class="w-full bg-dark border border-gray-800 rounded px-3 py-2 text-sm font-mono text-gray-200" placeholder="10.0.0.0/24">
          </div>
          <div>
            <label class="text-xs font-mono text-gray-400 block mb-1">Goal</label>
            <input v-model="form.goal" class="w-full bg-dark border border-gray-800 rounded px-3 py-2 text-sm font-mono text-gray-200" placeholder="Data exfiltration">
          </div>
          <div>
            <label class="text-xs font-mono text-gray-400 block mb-1">Profile</label>
            <select v-model="form.profile" class="w-full bg-dark border border-gray-800 rounded px-3 py-2 text-sm font-mono text-gray-200">
              <option value="stealth">Stealth (low noise, slow)</option>
              <option value="aggressive">Aggressive (fast, loud)</option>
              <option value="audit">Audit (requires human approval)</option>
            </select>
          </div>
          <div class="flex items-center gap-2">
            <input type="checkbox" v-model="form.autoApprove" id="auto" class="accent-purple">
            <label for="auto" class="text-xs font-mono text-gray-400">Auto-approve safe decisions</label>
          </div>
          <div class="flex gap-2 pt-2">
            <button @click="createCampaign" :disabled="creating" class="btn flex-1 text-sm py-2">
              {{ creating ? 'Creating...' : 'Create Campaign' }}
            </button>
            <button @click="showWizard = false" class="btn text-sm py-2 px-4">Cancel</button>
          </div>
        </div>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto space-y-2">
      <div v-if="campaigns.length === 0" class="text-gray-600 text-center py-12 font-mono text-sm">
        <div class="text-4xl mb-3">🎯</div>
        No campaigns. Create one to start operations.
      </div>

      <div v-for="c in campaigns" :key="c.id" class="glass-panel p-3 cursor-pointer"
        :class="c.id === activeId ? 'border border-purple' : 'border border-transparent'"
        @click="selectCampaign(c.id)">
        <div class="flex items-center justify-between mb-1">
          <span class="text-sm font-mono text-gray-200">{{ c.name }}</span>
          <div class="flex gap-2">
            <span class="text-[10px] px-2 py-0.5 rounded font-mono" :class="statusBadge(c.status)">
              {{ c.status?.toUpperCase() || 'DRAFT' }}
            </span>
            <span class="text-[10px] font-mono text-gray-600">{{ c.phase || 'recon' }}</span>
          </div>
        </div>
        <div class="text-xs font-mono text-gray-600 mb-2">{{ c.target_scope || c.targetScope }} · {{ c.goal }}</div>
        <div class="h-1 bg-dark rounded-full overflow-hidden mb-1">
          <div class="h-full bg-gradient-to-r from-purple to-neon transition-all"
            :style="{ width: ((c.progress || 0) * 100) + '%' }"></div>
        </div>
        <div class="flex items-center justify-between text-[10px] font-mono text-gray-600">
          <span>{{ Math.round((c.progress || 0) * 100) }}% · {{ c.agent_count || 0 }} agents</span>
          <div class="flex gap-1">
            <button v-if="c.status === 'running'" @click.stop="pauseCampaign(c.id)" class="px-2 py-0.5 rounded bg-yellow-900/20 text-yellow-400 border border-yellow-800/30 text-[10px]">⏸</button>
            <button v-if="c.status === 'paused'" @click.stop="resumeCampaign(c.id)" class="px-2 py-0.5 rounded bg-neon/10 text-neon border border-neon/30 text-[10px]">▶</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useCampaignStore } from '../stores/index.js'

const campaignStore = useCampaignStore()
const showWizard = ref(false)
const creating = ref(false)
const activeId = computed(() => campaignStore.activeCampaign?.id)

const form = ref({ name: '', targetScope: '', goal: '', profile: 'stealth', autoApprove: false })

onMounted(() => campaignStore.fetchCampaigns().catch(() => {}))

const campaigns = computed(() => campaignStore.campaigns || [])

async function createCampaign() {
  creating.value = true
  try {
    await campaignStore.createCampaign(
      form.value.name || 'AutoCampaign',
      form.value.targetScope || '10.0.0.0/24',
      form.value.goal || 'Data exfiltration',
      form.value.profile || 'stealth',
      form.value.autoApprove,
    )
    showWizard.value = false
    form.value = { name: '', targetScope: '', goal: '', profile: 'stealth', autoApprove: false }
  } finally { creating.value = false }
}

function selectCampaign(id) { campaignStore.getCampaign(id).catch(() => {}) }
function pauseCampaign(id) { campaignStore.pauseCampaign(id).catch(() => {}) }
function resumeCampaign(id) { campaignStore.resumeCampaign(id).catch(() => {}) }

function statusBadge(s) {
  const m = {
    'running': 'bg-neon/10 text-neon border border-neon/30',
    'paused': 'bg-yellow-900/20 text-yellow-400 border border-yellow-800/30',
    'completed': 'bg-green-900/20 text-green-400 border border-green-800/30',
    'draft': 'bg-gray-800 text-gray-400 border border-gray-700',
    'failed': 'bg-red-900/20 text-red-400 border border-red-800/30',
  }
  return m[s] || m['draft']
}
</script>
