/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        neon: '#00ff41',
        purple: '#6c63ff',
        alert: '#ff4444',
        dark: '#0a0a0f',
        panel: 'rgba(255,255,255,0.03)',
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'monospace'],
        sans: ['Inter', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
