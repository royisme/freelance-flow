export interface DateOnlyParts {
  year: number
  month: number
  day: number
}

export function parseISODateOnly(dateStr: string): DateOnlyParts | null {
  const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(dateStr)
  if (!match) return null

  const year = Number(match[1])
  const month = Number(match[2])
  const day = Number(match[3])

  if (!Number.isInteger(year) || !Number.isInteger(month) || !Number.isInteger(day)) return null
  if (month < 1 || month > 12) return null
  if (day < 1 || day > 31) return null

  return { year, month, day }
}

export function dateOnlySortKey(dateStr: string): number {
  const dateOnly = parseISODateOnly(dateStr)
  if (dateOnly) return dateOnly.year * 10000 + dateOnly.month * 100 + dateOnly.day

  const ts = new Date(dateStr).getTime()
  return Number.isFinite(ts) ? ts : 0
}

export function formatISODateOnlyForLocale(dateStr: string, locale: string): string {
  const dateOnly = parseISODateOnly(dateStr)
  const date = dateOnly ? new Date(dateOnly.year, dateOnly.month - 1, dateOnly.day) : new Date(dateStr)

  if (Number.isNaN(date.getTime())) return dateStr

  return new Intl.DateTimeFormat(locale, {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  }).format(date)
}

