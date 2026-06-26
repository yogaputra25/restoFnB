const apiBase = process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080'

export default defineNuxtConfig({
  compatibilityDate: '2026-06-16',
  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
  css: ['~/assets/css/tokens.css', '~/assets/css/main.css'],
  tailwindcss: {
    configPath: 'tailwind.config.ts',
  },
  runtimeConfig: {
    public: {
      apiBase: '',
    },
  },
  components: {
    dirs: [
      { path: '~/components/ui', prefix: '' },
      { path: '~/components/admin', prefix: '' },
      '~/components',
    ],
  },
  routeRules: {
    '/api/**': { proxy: `${apiBase}/api/**` },
    '/uploads/**': { proxy: `${apiBase}/uploads/**` },
  },
  app: {
    head: {
      title: 'restoFnB',
    },
  },
})
