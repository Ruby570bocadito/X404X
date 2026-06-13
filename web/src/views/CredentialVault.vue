<template>
  <div class="h-full flex flex-col">
    <div class="flex items-center justify-between mb-3">
      <h2 class="text-lg font-mono text-purple">CREDENTIAL VAULT</h2>
      <input v-model="filter" class="bg-dark border border-gray-800 rounded px-3 py-1 text-xs font-mono text-gray-300 w-48" placeholder="Filter...">
    </div>

    <div class="flex-1 overflow-auto glass-panel p-0">
      <table class="w-full text-xs font-mono">
        <thead class="sticky top-0 bg-panel">
          <tr class="text-gray-500 text-left">
            <th class="p-2">User</th>
            <th class="p-2">Domain</th>
            <th class="p-2">Hash Type</th>
            <th class="p-2">Hash</th>
            <th class="p-2">Source</th>
            <th class="p-2">Agent</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(c, i) in filteredCreds" :key="i" class="border-t border-gray-800 hover:bg-dark/50">
            <td class="p-2 text-gray-300">{{ c.username || '-' }}</td>
            <td class="p-2 text-gray-300">{{ c.domain || '-' }}</td>
            <td class="p-2">
              <span class="px-1.5 py-0.5 rounded text-[10px]" :class="hashBadge(c.hash_type)">
                {{ c.hash_type || 'N/A' }}
              </span>
            </td>
            <td class="p-2 text-gray-500 font-mono truncate max-w-[200px]" :title="c.hash">{{ c.hash?.substring(0, 32) || '-' }}</td>
            <td class="p-2 text-gray-500">{{ c.source || '-' }}</td>
            <td class="p-2 text-gray-500">{{ c.agent_id || '-' }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="filteredCreds.length === 0" class="text-gray-600 text-center py-12 font-mono text-sm">
        {{ creds.length === 0 ? 'No credentials captured yet.' : 'No matches for filter.' }}
      </div>
    </div>
    <div class="text-[10px] font-mono text-gray-600 mt-2">
      {{ filteredCreds.length }} / {{ creds.length }} credentials
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const creds = ref([])
const filter = ref('')

onMounted(async () => {
  try {
    const r = await fetch('/api/creds')
    if (r.ok) creds.value = await r.json()
  } catch { creds.value = [] }
})

const filteredCreds = computed(() => {
  if (!filter.value) return creds.value
  const f = filter.value.toLowerCase()
  return creds.value.filter(c =>
    (c.username || '').toLowerCase().includes(f) ||
    (c.domain || '').toLowerCase().includes(f) ||
    (c.source || '').toLowerCase().includes(f) ||
    (c.agent_id || '').toLowerCase().includes(f)
  )
})

function hashBadge(t) {
  const m = {
    'NTLM': 'bg-yellow-900/20 text-yellow-400 border border-yellow-800/30',
    'SHA1': 'bg-orange-900/20 text-orange-400 border border-orange-800/30',
    'SHA256': 'bg-red-900/20 text-red-400 border border-red-800/30',
    'MD5': 'bg-gray-800 text-gray-400 border border-gray-700',
    'bcrypt': 'bg-green-900/20 text-green-400 border border-green-800/30',
    'Kerberos': 'bg-purple/20 text-purple border border-purple/30',
  }
  return m[t] || 'bg-gray-800 text-gray-400 border border-gray-700'
}
</script>
