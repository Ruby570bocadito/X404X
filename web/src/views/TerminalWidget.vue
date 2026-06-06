<template>
  <div class="glass-panel p-3 h-full flex flex-col">
    <h3 class="text-sm font-mono text-purple mb-2">EMBEDDED TERMINAL</h3>
    <div class="flex-1 bg-black rounded p-2 font-mono text-xs overflow-y-auto" ref="terminalEl">
      <div class="text-neon mb-2">X404X Console v1.0 — Connected via API</div>
      <div v-for="(line, i) in terminalLines" :key="i" class="mb-1">
        <span class="text-neon">x404x > </span>
        <span :class="line.startsWith('[-]') ? 'text-alert' : line.startsWith('[+]') ? 'text-neon' : 'text-gray-300'">{{ line }}</span>
      </div>
      <div class="flex items-center">
        <span class="text-neon">x404x > </span>
        <input v-model="cmdInput" @keyup.enter="execCmd"
          class="flex-1 bg-transparent border-none outline-none text-gray-300 font-mono text-xs ml-1"
          ref="cmdInputEl" autofocus />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { useAgentStore, useReconStore, useAIStore, useCampaignStore } from '../stores/index.js'

const cmdInput = ref('')
const terminalLines = ref([
  'X404X Terminal v1.0 — Type "help" for commands.',
])

const agentStore = useAgentStore()
const reconStore = useReconStore()
const aiStore = useAIStore()
const campaignStore = useCampaignStore()

const execCmd = async () => {
  if (!cmdInput.value.trim()) return
  const cmd = cmdInput.value.trim()
  terminalLines.value.push(cmd)
  cmdInput.value = ''

  const parts = cmd.split(/\s+/)
  const action = parts[0]
  const args = parts.slice(1)

  try {
    switch (action) {
    case 'help':
      terminalLines.value.push('  sessions | hosts | creds | vulns | suggest | ai <text> | clear')
      break
    case 'sessions':
      await agentStore.fetchAgents()
      terminalLines.value.push(`  ${agentStore.activeAgents.length} active sessions:`)
      agentStore.activeAgents.forEach(a => {
        terminalLines.value.push(`  ${a.session_id || a.id}  ${a.local_ip || '?'}  ${a.os || '?'}  ${a.username || '?'}  ${a.status}`)
      })
      break
    case 'hosts':
      await reconStore.fetchHosts()
      terminalLines.value.push(`  ${reconStore.hosts.length} hosts:`)
      reconStore.hosts.forEach(h => {
        terminalLines.value.push(`  ${h.ip}  ${h.hostname || '?'}  ${h.os || '?'}  [${(h.open_ports || []).join(',')}]`)
      })
      break
    case 'creds':
      terminalLines.value.push('  Credentials from API (feature pending)')
      break
    case 'vulns':
      await reconStore.fetchVulnerabilities()
      terminalLines.value.push(`  ${reconStore.vulnerabilities.length} vulns:`)
      reconStore.vulnerabilities.forEach(v => {
        terminalLines.value.push(`  ${v.cve || '?'}  ${v.severity || '?'}  ${v.service}:${v.port}  ${v.target_ip || '?'}`)
      })
      break
    case 'suggest':
      const cid = campaignStore.activeCampaign?.id
      await aiStore.getSuggestions(cid)
      terminalLines.value.push(`  ${aiStore.suggestions.length} suggestions from AI:`)
      aiStore.suggestions.slice(0, 5).forEach(s => {
        terminalLines.value.push(`  [${s.source}] ${s.technique} → ${s.target || 'any'} (conf=${(s.confidence * 100).toFixed(0)}%)`)
      })
      break
    case 'clear':
      terminalLines.value = []
      break
    case 'ai':
      if (args.length > 0) {
        const response = await aiStore.chat(args.join(' '))
        terminalLines.value.push(`[AI] ${response}`)
      } else {
        terminalLines.value.push('[-] Usage: ai <prompt>')
      }
      break
    default:
      terminalLines.value.push(`[-] Unknown: ${cmd}. Try: help, sessions, hosts, creds, vulns, suggest, ai, clear`)
    }
  } catch (e) {
    terminalLines.value.push(`[-] Error: ${e.message}`)
  }

  nextTick(() => {
    const el = document.querySelector('.overflow-y-auto')
    if (el) el.scrollTop = el.scrollHeight
  })
}
</script>
