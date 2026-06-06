<template>
  <div class="glass-panel p-3 h-full flex flex-col">
    <h3 class="text-sm font-mono text-purple mb-2">EMBEDDED TERMINAL</h3>
    <div class="flex-1 bg-black rounded p-2 font-mono text-xs overflow-y-auto" ref="terminalEl">
      <div class="text-neon mb-2">X404X Console v1.0 — Connected to orchestrator</div>
      <div v-for="(line, i) in terminalLines" :key="i" class="mb-1">
        <span class="text-neon">x404x > </span>
        <span class="text-gray-300">{{ line }}</span>
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

const cmdInput = ref('')
const terminalLines = ref([
  'help',
  '  Core | Campaign | Module | AI | Database | Lab',
  '  Type "help" for full command list.',
  'sessions',
  '  Id  Target        OS            User                    Status',
  '  1   10.0.0.10     Windows 2019  NT AUTHORITY\\SYSTEM     active',
  '  2   10.0.0.20     Ubuntu 24.04  root                    active',
])

const execCmd = () => {
  if (!cmdInput.value.trim()) return
  terminalLines.value.push(cmdInput.value)
  const cmd = cmdInput.value.trim()
  if (cmd === 'sessions') {
    terminalLines.value.push('  Id  Target        OS            User                    Status')
    terminalLines.value.push('  1   10.0.0.10     Windows 2019  NT AUTHORITY\\SYSTEM     active')
    terminalLines.value.push('  2   10.0.0.20     Ubuntu 24.04  root                    active')
  } else if (cmd === 'hosts') {
    terminalLines.value.push('  10.0.0.10  DC    Windows 2019  compromised')
    terminalLines.value.push('  10.0.0.20  DB    Ubuntu 24.04  scanned')
  } else if (cmd === 'help') {
    terminalLines.value.push('  Core | Campaign | Module | AI | Database | Lab')
  } else if (cmd.startsWith('use ')) {
    terminalLines.value.push(`[+] ${cmd.split(' ')[1]} loaded. Use 'show options'.`)
  } else if (cmd === 'ai suggest') {
    terminalLines.value.push('[AI] 1. EternalBlue [0.95]  2. Kerberoast [0.82]')
  } else if (cmd === 'clear') {
    terminalLines.value = []
  } else {
    terminalLines.value.push(`[*] Command executed: ${cmd}`)
  }
  cmdInput.value = ''
  nextTick(() => {
    const el = document.querySelector('.overflow-y-auto')
    if (el) el.scrollTop = el.scrollHeight
  })
}
</script>
