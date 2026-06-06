<template>
  <div class="glass-panel p-4 space-y-4">
    <h3 class="text-sm font-mono text-purple">X404X — Architecture & Documentation</h3>
    <div class="prose prose-invert max-w-none text-sm">
      <h4 class="text-neon font-mono">X404X Framework</h4>
      <p class="text-gray-400">
        Semi-autonomous Red Team platform covering the complete cyber kill chain.
        Integrates 11 specialized offensive security tools into a unified monorepo
        with gRPC protocol, encrypted communication, AI-powered decision making (Ollama),
        kernel-level persistence, and a modern web dashboard.
      </p>
      <h4 class="text-neon font-mono mt-4">Components</h4>
      <table class="w-full text-xs font-mono">
        <thead><tr class="text-gray-600"><th class="text-left py-1">Component</th><th class="text-left">Language</th><th class="text-left">Phase</th></tr></thead>
        <tbody>
          <tr><td class="text-purple py-1">CLI (x404x)</td><td class="text-gray-400">Go</td><td class="text-gray-400">All</td></tr>
          <tr><td class="text-purple">Orchestrator</td><td class="text-gray-400">Go</td><td class="text-gray-400">All</td></tr>
          <tr><td class="text-purple">Pulse-C2</td><td class="text-gray-400">Go + Vue 3</td><td class="text-gray-400">C2</td></tr>
          <tr><td class="text-purple">Rise-Privilege</td><td class="text-gray-400">Go</td><td class="text-gray-400">PrivEsc</td></tr>
          <tr><td class="text-purple">Vault-Kernel</td><td class="text-gray-400">C + Go</td><td class="text-gray-400">Persistence</td></tr>
          <tr><td class="text-purple">Horizon-Intel</td><td class="text-gray-400">Python</td><td class="text-gray-400">Recon</td></tr>
          <tr><td class="text-purple">Wormy-ML</td><td class="text-gray-400">Python</td><td class="text-gray-400">Lateral</td></tr>
          <tr><td class="text-purple">Specter + Apex</td><td class="text-gray-400">Python</td><td class="text-gray-400">AI</td></tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
