import { Code, ConnectError } from '@connectrpc/connect'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { RoomCodePanel } from './RoomCodePanel'

const mocks = vi.hoisted(() => ({
  setRoomCode: vi.fn(),
  refreshRoomAccess: vi.fn().mockResolvedValue(undefined),
}))

vi.mock('@/api/roomCodeClient', () => ({
  setRoomCode: mocks.setRoomCode,
}))

vi.mock('@/contexts/roomAccessContext', () => ({
  useRoomAccess: () => ({
    hasAccess: false,
    isLoading: false,
    refreshRoomAccess: mocks.refreshRoomAccess,
  }),
}))

function renderPanel(initialPath = '/room?react=true') {
  const router = createMemoryRouter(
    [
      { path: '/', element: <div>Home</div> },
      { path: '/room', element: <RoomCodePanel /> },
    ],
    { initialEntries: [initialPath] },
  )
  return render(<RouterProvider router={router} />)
}

describe('RoomCodePanel', () => {
  beforeEach(() => {
    mocks.setRoomCode.mockReset()
    mocks.refreshRoomAccess.mockReset()
    mocks.refreshRoomAccess.mockResolvedValue(undefined)
    mocks.setRoomCode.mockResolvedValue(undefined)
  })

  it('submits a trimmed room code, refreshes access, and navigates home', async () => {
    renderPanel()
    const input = screen.getByRole('textbox', { name: /room code/i })
    fireEvent.change(input, { target: { value: '  GOOD123  ' } })
    fireEvent.click(screen.getByRole('button', { name: /^continue$/i }))

    await waitFor(() => {
      expect(mocks.setRoomCode).toHaveBeenCalledWith('GOOD123')
    })
    expect(mocks.refreshRoomAccess).toHaveBeenCalledTimes(1)
    await screen.findByText('Home')
  })

  it('shows the invalid/expired message for InvalidArgument', async () => {
    mocks.setRoomCode.mockRejectedValue(new ConnectError('bad', Code.InvalidArgument))
    renderPanel()
    fireEvent.change(screen.getByRole('textbox', { name: /room code/i }), {
      target: { value: 'BAD' },
    })
    fireEvent.click(screen.getByRole('button', { name: /^continue$/i }))

    const alert = await screen.findByRole('alert')
    expect(alert).toHaveTextContent('That code is not valid or has expired.')
  })

  it('shows a ConnectError message for other RPC errors', async () => {
    mocks.setRoomCode.mockRejectedValue(new ConnectError('Server blew up', Code.Unavailable))
    renderPanel()
    fireEvent.change(screen.getByRole('textbox', { name: /room code/i }), {
      target: { value: 'X' },
    })
    fireEvent.click(screen.getByRole('button', { name: /^continue$/i }))

    const alert = await screen.findByRole('alert')
    expect(alert.textContent).toMatch(/server blew up/i)
  })

  it('shows a generic message for unexpected errors', async () => {
    mocks.setRoomCode.mockRejectedValue(new Error('network down'))
    renderPanel()
    fireEvent.change(screen.getByRole('textbox', { name: /room code/i }), {
      target: { value: 'X' },
    })
    fireEvent.click(screen.getByRole('button', { name: /^continue$/i }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Something went wrong.')
  })

  it('shows a loading label while the request is in flight', async () => {
    let release!: () => void
    const pending = new Promise<void>((resolve) => {
      release = resolve
    })
    mocks.setRoomCode.mockImplementation(() => pending)

    renderPanel()
    fireEvent.change(screen.getByRole('textbox', { name: /room code/i }), {
      target: { value: 'WAIT' },
    })
    fireEvent.click(screen.getByRole('button', { name: /^continue$/i }))

    expect(await screen.findByRole('button', { name: /checking/i })).toBeDisabled()
    release()
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /^continue$/i })).toBeInTheDocument()
    })
  })
})
