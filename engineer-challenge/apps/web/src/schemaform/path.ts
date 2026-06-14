/* eslint-disable @typescript-eslint/no-explicit-any */
export function getByPath(obj: any, path: string): any {
  return path.split('.').reduce((o, k) => (o == null ? undefined : o[k]), obj)
}

export function setByPath<T>(obj: T, path: string, value: unknown): T {
  const keys = path.split('.')
  const copy: any = Array.isArray(obj) ? [...(obj as any)] : { ...(obj as any) }
  let cur = copy
  for (let i = 0; i < keys.length - 1; i++) {
    const k = keys[i]
    cur[k] = Array.isArray(cur[k]) ? [...cur[k]] : { ...(cur[k] ?? {}) }
    cur = cur[k]
  }
  cur[keys[keys.length - 1]] = value
  return copy
}
