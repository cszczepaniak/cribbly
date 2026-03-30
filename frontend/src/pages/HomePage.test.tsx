import { render, screen } from "@testing-library/react"
import { MemoryRouter } from "react-router-dom"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { HomePage } from "./HomePage"

const mocks = vi.hoisted(() => {
  const refreshRoomAccess = vi.fn()
  return {
    refreshRoomAccess,
    useRoomAccess: vi.fn(() => ({
      hasAccess: true as boolean | null,
      isLoading: false,
      refreshRoomAccess,
    })),
  }
})

vi.mock("@/contexts/roomAccessContext", () => ({
  useRoomAccess: () => mocks.useRoomAccess(),
}))

function renderHome() {
  return render(
    <MemoryRouter>
      <HomePage />
    </MemoryRouter>,
  )
}

describe("HomePage", () => {
  beforeEach(() => {
    mocks.refreshRoomAccess.mockReset()
    mocks.useRoomAccess.mockImplementation(() => ({
      hasAccess: true,
      isLoading: false,
      refreshRoomAccess: mocks.refreshRoomAccess,
    }))
  })

  it("shows a loading state while room access is loading", () => {
    mocks.useRoomAccess.mockImplementation(() => ({
      hasAccess: null,
      isLoading: true,
      refreshRoomAccess: mocks.refreshRoomAccess,
    }))

    renderHome()

    expect(screen.getByText("Loading...")).toBeInTheDocument()
  })

  it("shows the room code panel when there is no access", () => {
    mocks.useRoomAccess.mockImplementation(() => ({
      hasAccess: false,
      isLoading: false,
      refreshRoomAccess: mocks.refreshRoomAccess,
    }))

    renderHome()

    expect(
      screen.getByRole("heading", { name: /enter room code/i }),
    ).toBeInTheDocument()
  })

  it("shows the landing content when access is granted", () => {
    renderHome()

    expect(
      screen.getByRole("heading", { name: /welcome to cribbly/i }),
    ).toBeInTheDocument()
    expect(screen.getByRole("link", { name: /divisions/i })).toBeInTheDocument()
  })
})
