import { useCallback, useSyncExternalStore } from "react"

/** localStorage key — only honored when `import.meta.env.DEV` is true. */
export const DEV_ADMIN_OVERRIDE_KEY = "cribbly:devAdminOverride"

function getSnapshot(): boolean {
  if (!import.meta.env.DEV) {
    return false
  }
  try {
    return localStorage.getItem(DEV_ADMIN_OVERRIDE_KEY) === "true"
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
    if (e.key === DEV_ADMIN_OVERRIDE_KEY || e.key === null) {
      onStoreChange()
    }
  }
  const onCustom = () => onStoreChange()
  window.addEventListener("storage", onStorage)
  window.addEventListener("cribbly-dev-admin-override", onCustom)
  return () => {
    window.removeEventListener("storage", onStorage)
    window.removeEventListener("cribbly-dev-admin-override", onCustom)
  }
}

/** Read override (e.g. from console). No-op in production builds. */
export function getDevAdminOverride(): boolean {
  return getSnapshot()
}

/** Persist override and notify subscribers (same tab + other tabs). No-op in production. */
export function setDevAdminOverride(value: boolean): void {
  if (!import.meta.env.DEV) {
    return
  }
  try {
    if (value) {
      localStorage.setItem(DEV_ADMIN_OVERRIDE_KEY, "true")
    } else {
      localStorage.removeItem(DEV_ADMIN_OVERRIDE_KEY)
    }
    window.dispatchEvent(new Event("cribbly-dev-admin-override"))
  } catch {
    // ignore quota / private mode
  }
}

/**
 * Dev-only: send X-Cribbly-Dev-Admin on API calls (needs matching server secret).
 * In production always `[false, no-op]`.
 */
export function useDevAdminOverride(): [boolean, (value: boolean) => void] {
  const enabled = import.meta.env.DEV
  const value = useSyncExternalStore(
    enabled ? subscribe : emptySubscribe,
    enabled ? getSnapshot : () => false,
    getServerSnapshot,
  )
  const set = useCallback((v: boolean) => {
    setDevAdminOverride(v)
  }, [])
  return [value, set]
}
