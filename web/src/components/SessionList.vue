<template>
  <div class="glass-panel p-3 h-48 flex flex-col">
    <h3 class="text-sm font-mono text-purple mb-2">ACTIVE SESSIONS</h3>
    <div class="flex-1 overflow-y-auto text-xs font-mono">
      <div class="grid grid-cols-4 text-gray-600 border-b border-gray-800 pb-1 mb-1">
        <span>ID</span><span>Target</span><span>OS</span><span>User</span>
      </div>
      <div v-for="s in sessions" :key="s.id" class="grid grid-cols-4 py-1 hover:bg-panel rounded cursor-pointer">
        <span class="text-neon">{{ s.id }}</span>
        <span class="text-gray-300">{{ s.target }}</span>
        <span class="text-gray-400">{{ s.os }}</span>
        <span :class="s.user.includes('SYSTEM') || s.user === 'root' ? 'text-alert' : 'text-gray-300'">{{ s.user }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
const sessions = [
  { id: 'S1', target: '10.0.0.10', os: 'Win2019', user: 'NT\\SYSTEM' },
  { id: 'S2', target: '10.0.0.20', os: 'Ubuntu24', user: 'root' },
  { id: 'S3', target: '10.0.0.50', os: 'Win11', user: 'user' },
]
</script>
