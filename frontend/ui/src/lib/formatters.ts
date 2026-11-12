export const COUNTRY_CODES = [
  { code: '+91', label: 'IN (+91)' },
  { code: '+1', label: 'US (+1)' },
  { code: '+44', label: 'UK (+44)' },
  { code: '+971', label: 'AE (+971)' },
] as const

export function formatINR(amount: number): string {
  return `₹${amount.toLocaleString('en-IN')}`
}

export function formatLakhs(amount: number): string {
  return `₹${(amount / 100000).toFixed(2)}L`
}

export function parseSafeDate(dateStr: string): Date | null {
  if (!dateStr) return null
  const parsed = new Date(dateStr)
  return isNaN(parsed.getTime()) ? null : parsed
}
