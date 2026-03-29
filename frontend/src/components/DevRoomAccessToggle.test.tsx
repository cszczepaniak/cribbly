import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it } from 'vitest'
import { DevRoomAccessToggle } from './DevRoomAccessToggle'
import { DEV_ROOM_ACCESS_OVERRIDE_KEY } from '@/lib/devRoomAccessOverride'

describe.skipIf(!import.meta.env.DEV)('DevRoomAccessToggle', () => {
  beforeEach(() => {
    localStorage.removeItem(DEV_ROOM_ACCESS_OVERRIDE_KEY)
  })

  it('renders the label and switch', () => {
    render(<DevRoomAccessToggle />)
    expect(screen.getByText('Pretend room access')).toBeInTheDocument()
    expect(
      screen.getByRole('switch', { name: /pretend room access \(dev only\)/i }),
    ).toBeInTheDocument()
  })

  it('persists the override in localStorage when turned on', () => {
    render(<DevRoomAccessToggle />)
    fireEvent.click(screen.getByRole('switch', { name: /pretend room access \(dev only\)/i }))
    expect(localStorage.getItem(DEV_ROOM_ACCESS_OVERRIDE_KEY)).toBe('true')
    expect(
      screen.getByRole('switch', { name: /pretend room access \(dev only\)/i }),
    ).toBeChecked()
  })

  it('clears localStorage when turned off', () => {
    localStorage.setItem(DEV_ROOM_ACCESS_OVERRIDE_KEY, 'true')
    render(<DevRoomAccessToggle />)
    expect(
      screen.getByRole('switch', { name: /pretend room access \(dev only\)/i }),
    ).toBeChecked()
    fireEvent.click(screen.getByRole('switch', { name: /pretend room access \(dev only\)/i }))
    expect(localStorage.getItem(DEV_ROOM_ACCESS_OVERRIDE_KEY)).toBeNull()
    expect(
      screen.getByRole('switch', { name: /pretend room access \(dev only\)/i }),
    ).not.toBeChecked()
  })
})
