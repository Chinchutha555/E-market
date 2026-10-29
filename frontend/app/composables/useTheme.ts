// สลับธีมสว่าง/มืดของทั้งระบบ เก็บค่าไว้ใน localStorage และตั้ง data-theme ที่ <html>
export type AppTheme = 'light' | 'dark'

const STORAGE_KEY = 'emarket-theme'

export const useTheme = () => {
  const theme = useState<AppTheme>('app-theme', () => 'light')

  const apply = (value: AppTheme) => {
    theme.value = value
    if (import.meta.client) {
      document.documentElement.dataset.theme = value
      try { localStorage.setItem(STORAGE_KEY, value) } catch {}
    }
  }

  const init = () => {
    if (!import.meta.client) return
    let saved: string | null = null
    try { saved = localStorage.getItem(STORAGE_KEY) } catch {}
    const preferred: AppTheme = saved === 'dark' || saved === 'light'
      ? saved
      : (window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
    theme.value = preferred
    document.documentElement.dataset.theme = preferred
  }

  const toggle = () => apply(theme.value === 'dark' ? 'light' : 'dark')

  return { theme, isDark: computed(() => theme.value === 'dark'), init, toggle, apply }
}
