import { defineStore } from 'pinia'
import { ref } from 'vue'

const API_BASE = '/api'

export const usePhantomStore = defineStore('phantom', () => {
  const nodes = ref([])
  const stats = ref({
    totalNodes: 0, activeNodes: 0, cookiesTotal: 0,
    sessionsTotal: 0, swPersisted: 0, meshLatency: 45,
  })
  const loading = ref(false)

  async function fetchStatus() {
    loading.value = true
    try {
      const res = await fetch(`${API_BASE}/phantom/status`)
      if (res.ok) {
        const data = await res.json()
        stats.value = { ...stats.value, ...data }
      }
    } catch {}
    loading.value = false
  }

  async function fetchNodes() {
    try {
      const res = await fetch(`${API_BASE}/phantom/nodes`)
      if (res.ok) nodes.value = await res.json()
    } catch {}
  }

  async function execute(action, params = {}) {
    try {
      const res = await fetch(`${API_BASE}/phantom/${action}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(params),
      })
      if (res.ok) {
        const result = await res.json()
        await Promise.all([fetchNodes(), fetchStatus()])
        return result
      }
    } catch (e) {
      console.error('Phantom action failed:', e)
    }
    return null
  }

  return { nodes, stats, loading, fetchStatus, fetchNodes, execute }
})
