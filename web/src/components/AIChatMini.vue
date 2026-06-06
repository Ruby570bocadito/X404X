<template>
  <div class="glass-panel p-3 h-full flex flex-col">
    <h3 class="text-sm font-mono text-purple mb-2">AI CONSOLE</h3>
    <div class="flex-1 overflow-y-auto space-y-2 text-xs font-mono mb-3">
      <div class="text-gray-500">> analyze target 10.0.0.10</div>
      <div class="bg-panel rounded p-2 space-y-1">
        <div class="text-neon">[Specter] Analysis:</div>
        <div class="text-gray-400">• SMB 445 → EternalBlue (MS17-010)</div>
        <div class="text-gray-400">• RDP 3389 → BlueKeep check</div>
        <div class="text-gray-400">• OS: Windows Server 2019</div>
      </div>
      <div class="bg-panel rounded p-2 space-y-1">
        <div class="text-purple">[Apex] Recommendations:</div>
        <div class="text-neon">1. EternalBlue [conf: 0.95]</div>
        <div class="text-gray-400">2. Kerberoast [conf: 0.82]</div>
        <div class="text-gray-400">3. PrintNightmare [conf: 0.71]</div>
      </div>
    </div>
    <div class="flex gap-2">
      <button class="btn btn-accept flex-1">Accept #1</button>
      <button class="btn btn-reject flex-1">Reject</button>
    </div>
  </div>
</template>
