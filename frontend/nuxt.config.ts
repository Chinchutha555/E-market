const configuredApiBase = process.env.NUXT_PUBLIC_API_BASE
  || process.env.NUXT_PUBLIC_API_URL
  || '/api/v1'

export default defineNuxtConfig({
  vite: {
    server: {
      // Used only for an explicitly enabled local testing tunnel. Production
      // builds do not run the Vite development server.
      allowedHosts: process.env.NUXT_ALLOW_TUNNEL_HOSTS === 'true' ? true : undefined
    }
  },
  future: {
    compatibilityVersion: 4,
  },
  compatibilityDate: '2024-04-03',
  modules: ['@pinia/nuxt'],
  devtools: { enabled: false },
  telemetry: false,
  ssr: false,
  routeRules: {
    '/api/v1/**': {
      proxy: `${process.env.BACKEND_API_URL || process.env.BACKEND_URL || 'http://127.0.0.1:8080'}/api/v1/**`
    }
  },
  experimental: {
    appManifest: false
  },
  components: [
    { path: '~/components', pathPrefix: false },
    { path: '~/components/buttons', pathPrefix: false },
    { path: '~/components/layouts', pathPrefix: false },
    { path: '~/components/modals', pathPrefix: false },
    { path: '~/components/maps', pathPrefix: false },
    { path: '~/components/extras', pathPrefix: false },
    { path: '~/core/modules/stall/ui/components', pathPrefix: false },
  ],
  css: ['maplibre-gl/dist/maplibre-gl.css', '~/assets/css/main.css'],
  app: {
    head: {
      title: 'e-Market | ระบบจัดการแผงเช่าตลาดและสัญญาการเช่า',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content: 'ระบบบริหารจัดการแผงเช่าตลาด สัญญาการเช่าผู้เช่า และติดตามรายรับอย่างมีประสิทธิภาพ รองรับ REST API และ JWT'
        }
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Bai+Jamjuree:ital,wght@0,200;0,300;0,400;0,500;0,600;0,700;1,200;1,300;1,400;1,500;1,600;1,700&family=Plus+Jakarta+Sans:wght@400;500;600;700;800&display=swap'
        }
      ]
    }
  },
  runtimeConfig: {
    public: {
      // Both HTTP clients must use the same backend. Previously one could use
      // apiUrl while the other silently stayed on the volatile Nitro API.
      apiBase: configuredApiBase,
      apiUrl: configuredApiBase
    }
  }
})
