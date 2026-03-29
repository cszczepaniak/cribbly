import {
  createContext,
  type ReactNode,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'

import { checkRoomAccess } from '@/api/roomCodeClient'
import { useDevRoomAccessOverride } from '@/lib/devRoomAccessOverride'

type RoomAccessContextValue = {
  /** `null` while the probe is still in flight (never when dev override is on). */
  hasAccess: boolean | null
  isLoading: boolean
  refreshRoomAccess: () => Promise<void>
}

const RoomAccessContext = createContext<RoomAccessContextValue | null>(null)

export function RoomAccessProvider({ children }: { children: ReactNode }) {
  const location = useLocation()
  const navigate = useNavigate()
  const [devOverride] = useDevRoomAccessOverride()
  const [probeAccess, setProbeAccess] = useState<boolean | null>(null)

  const hasAccess = devOverride ? true : probeAccess
  const isLoading = devOverride ? false : probeAccess === null

  const refreshRoomAccess = useCallback(async () => {
    if (devOverride) {
      return
    }
    try {
      const res = await checkRoomAccess()
      setProbeAccess(res.hasAccess)
    } catch {
      setProbeAccess(false)
    }
  }, [devOverride])

  useEffect(() => {
    if (devOverride) {
      return
    }
    const ac = new AbortController()
    ;(async () => {
      try {
        const res = await checkRoomAccess()
        if (ac.signal.aborted) {
          return
        }
        setProbeAccess(res.hasAccess)
      } catch {
        if (ac.signal.aborted) {
          return
        }
        setProbeAccess(false)
      }
    })()
    return () => ac.abort()
  }, [devOverride, location.pathname, location.search])

  useEffect(() => {
    if (isLoading) {
      return
    }
    if (location.pathname !== '/' && !hasAccess) {
      navigate('/', { replace: true })
    }
  }, [hasAccess, isLoading, location.pathname, navigate])

  const value = useMemo<RoomAccessContextValue>(
    () => ({
      hasAccess,
      isLoading,
      refreshRoomAccess,
    }),
    [hasAccess, isLoading, refreshRoomAccess],
  )

  return <RoomAccessContext.Provider value={value}>{children}</RoomAccessContext.Provider>
}

export function useRoomAccess() {
  const ctx = useContext(RoomAccessContext)
  if (!ctx) {
    throw new Error('useRoomAccess must be used within RoomAccessProvider')
  }
  return ctx
}

/** Renders child routes; blocks non-home routes until the first room-access probe completes. */
export function RoomAccessOutlet() {
  const location = useLocation()
  const { hasAccess, isLoading } = useRoomAccess()

  if (location.pathname !== '/' && isLoading) {
    return (
      <div className="text-muted-foreground flex flex-1 flex-col items-center justify-center p-8 text-sm">
        Loading…
      </div>
    )
  }

  if (location.pathname !== '/' && hasAccess === false) {
    return null
  }

  return <Outlet />
}
