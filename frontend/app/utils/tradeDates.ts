const DAY_MS = 86400000

const parseDate = (value: string): Date | null => {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return null
  const date = new Date(value + 'T00:00:00Z')
  if (!Number.isFinite(date.getTime()) || date.toISOString().slice(0, 10) !== value) return null
  return date
}

export const calculateTradeEndDate = (start: string, duration: string | number, unit: string): string => {
  const date = parseDate(start)
  const amount = Number(duration)
  if (!date || !Number.isInteger(amount) || amount < 1 || amount > 10000) return ''
  if (unit === 'days') {
    date.setUTCDate(date.getUTCDate() + amount)
  } else if (unit === 'months' || unit === 'years') {
    const day = date.getUTCDate()
    const targetMonth = date.getUTCMonth() + (unit === 'years' ? amount * 12 : amount)
    date.setUTCDate(1)
    date.setUTCMonth(targetMonth)
    const lastDay = new Date(date.getTime())
    lastDay.setUTCMonth(lastDay.getUTCMonth() + 1, 0)
    date.setUTCDate(Math.min(day, lastDay.getUTCDate()))
  } else return ''
  date.setUTCDate(date.getUTCDate() - 1)
  const result = date.toISOString().slice(0, 10)
  return /^\d{4}-\d{2}-\d{2}$/.test(result) ? result : ''
}

export const tradeDaysRemaining = (end: string, today: string): number | null => {
  const endDate = parseDate(end)
  const todayDate = parseDate(today)
  return endDate && todayDate ? Math.round((endDate.getTime() - todayDate.getTime()) / DAY_MS) : null
}

export const bangkokToday = (): string => {
  const parts = new Intl.DateTimeFormat('en-GB', { timeZone: 'Asia/Bangkok', year: 'numeric', month: '2-digit', day: '2-digit' }).formatToParts(new Date())
  const part = (name: string) => parts.find(item => item.type === name)?.value || ''
  return part('year') + '-' + part('month') + '-' + part('day')
}
