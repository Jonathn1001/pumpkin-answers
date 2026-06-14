import type { ReactNode } from 'react'

const tones = {
  green: 'bg-green-100 text-green-800', blue: 'bg-blue-100 text-blue-800',
  amber: 'bg-amber-100 text-amber-800', gray: 'bg-gray-100 text-gray-700', red: 'bg-red-100 text-red-800',
}

export function Badge({ children, tone = 'gray' }: { children: ReactNode; tone?: keyof typeof tones }) {
  return <span className={`rounded px-2 py-0.5 text-xs font-medium ${tones[tone]}`}>{children}</span>
}
