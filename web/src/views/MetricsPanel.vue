<template>
  <div class="glass-panel p-4">
    <h3 class="text-sm font-mono text-purple mb-3">BLUEFORGE METRICS</h3>
    <div class="grid grid-cols-2 gap-3">
      <div class="glass-panel p-3 text-center">
        <div class="text-3xl font-mono text-neon">87%</div>
        <div class="text-xs text-gray-500 mt-1">Evasion Rate</div>
      </div>
      <div class="glass-panel p-3 text-center">
        <div class="text-3xl font-mono text-alert">2</div>
        <div class="text-xs text-gray-500 mt-1">Detections</div>
      </div>
    </div>
    <div class="mt-3 text-xs font-mono">
      <div class="text-gray-500 mb-2">Detection Log:</div>
      <div class="space-y-1">
        <div class="flex gap-2"><span class="text-purple">[Suricata]</span><span class="text-gray-300">Scan detected on .20</span><span class="text-gray-600 ml-auto">14:30:00</span></div>
        <div class="flex gap-2"><span class="text-purple">[Suricata]</span><span class="text-gray-300">SMB exploit attempt flagged</span><span class="text-gray-600 ml-auto">14:31:15</span></div>
        <div class="flex gap-2"><span class="text-neon">[Bypassed]</span><span class="text-gray-400">Wormy-ML polymorphic engine</span><span class="text-gray-600 ml-auto">14:31:45</span></div>
      </div>
    </div>
  </div>
</template>
