import { useCallback, useSyncExternalStore } from 'react'

/** localStorage key — only honored when `import.meta.env.DEV` is true. */
export const DEV_ROOM_ACCESS_OVERRIDE_KEY = 'cribbly:devRoomAccessOverride'

function getSnapshot(): boolean {
  if (!import.meta.env.DEV) {
    return false
  }
  try {
    return localStorage.getItem(DEV_ROOM_ACCESS_OVERRIDE_KEY) === 'true'
  } catch {
    return false
  }
}

function getServerSnapshot(): boolean {
  return false
}

const emptySubscribe = () => () => {}

function subscribe(onStoreChange: () => void): () => void {
  if (!import.meta.env.DEV) {
    return () => {}
  }
  const onStorage = (e: StorageEvent) => {
    if (e.key === DEV_ROOM_ACCESS_OVERRIDE_KEY || e.key === null) {
      onStoreChange()
    }
  }
  const onCustom = () => onStoreChange()
  window.addEventListener('storage', onStorage)
  window.addEventListener('cribbly-dev-room-override', onCustom)
  return () => {
    window.removeEventListener('storage', onStorage)
    window.removeEventListener('cribbly-dev-room-override', onCustom)
  }
}

/** Read override (e.g. from console). No-op in production builds. */
export function getDevRoomAccessOverride(): boolean {
  return getSnapshot()
}

/** Persist override and notify subscribers (same tab + other tabs). No-op in production. */
export function setDevRoomAccessOverride(value: boolean): void {
  if (!import.meta.env.DEV) {
    return
  }
  try {
    if (value) {
      localStorage.setItem(DEV_ROOM_ACCESS_OVERRIDE_KEY, 'true')
    } else {
      localStorage.removeItem(DEV_ROOM_ACCESS_OVERRIDE_KEY)
    }
    window.dispatchEvent(new Event('cribbly-dev-room-override'))
  } catch {
    // ignore quota / private mode
  }
}

/**
 * Dev-only: pretend the browser has room access (skips probe + redirect to `/`).
 * In production always `[false, no-op]`.
 */
export function useDevRoomAccessOverride(): [boolean, (value: boolean) => void] {
  const enabled = import.meta.env.DEV
  const value = useSyncExternalStore(
    enabled ? subscribe : emptySubscribe,
    enabled ? getSnapshot : () => false,
    getServerSnapshot,
  )
  const set = useCallback((v: boolean) => {
    setDevRoomAccessOverride(v)
  }, [])
  return [value, set]
}
