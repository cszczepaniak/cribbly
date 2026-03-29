import { beforeEach, describe, expect, it } from 'vitest'
import {
  DEV_ROOM_ACCESS_OVERRIDE_KEY,
  getDevRoomAccessOverride,
  setDevRoomAccessOverride,
} from './devRoomAccessOverride'

describe.skipIf(!import.meta.env.DEV)('devRoomAccessOverride', () => {
  beforeEach(() => {
    localStorage.removeItem(DEV_ROOM_ACCESS_OVERRIDE_KEY)
  })

  it('returns false when the key is unset', () => {
    expect(getDevRoomAccessOverride()).toBe(false)
  })

  it('round-trips through localStorage', () => {
    setDevRoomAccessOverride(true)
    expect(localStorage.getItem(DEV_ROOM_ACCESS_OVERRIDE_KEY)).toBe('true')
    expect(getDevRoomAccessOverride()).toBe(true)
    setDevRoomAccessOverride(false)
    expect(localStorage.getItem(DEV_ROOM_ACCESS_OVERRIDE_KEY)).toBeNull()
    expect(getDevRoomAccessOverride()).toBe(false)
  })
})
