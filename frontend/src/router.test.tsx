import { render, screen } from "@testing-library/react"
import { createMemoryRouter, RouterProvider } from "react-router-dom"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { DEV_ROOM_ACCESS_OVERRIDE_KEY } from "@/lib/devRoomAccessOverride"
import { routeObjects } from "./router"

const roomCodeMocks = vi.hoisted(() => ({
  checkRoomAccess: vi.fn().mockResolvedValue({ hasAccess: true }),
  setRoomCode: vi.fn(),
  doSomething: vi.fn(),
}))

vi.mock("@/api/roomCodeClient", () => ({
  checkRoomAccess: roomCodeMocks.checkRoomAccess,
  setRoomCode: roomCodeMocks.setRoomCode,
  doSomething: roomCodeMocks.doSomething,
}))

describe("router", () => {
  beforeEach(() => {
    globalThis.localStorage?.removeItem?.(DEV_ROOM_ACCESS_OVERRIDE_KEY)
    roomCodeMocks.checkRoomAccess.mockResolvedValue({ hasAccess: true })
  })

  it("renders home", async () => {
    const router = createMemoryRouter(routeObjects, {
      initialEntries: ["/?react=true"],
      basename: "/",
    })
    render(<RouterProvider router={router} />)
    expect(
      await screen.findByRole("heading", { name: /welcome to cribbly/i }),
    ).toBeInTheDocument()
  })

  it("shows room code entry when the user has no room access", async () => {
    roomCodeMocks.checkRoomAccess.mockResolvedValue({ hasAccess: false })
    const router = createMemoryRouter(routeObjects, {
      initialEntries: ["/?react=true"],
      basename: "/",
    })
    render(<RouterProvider router={router} />)
    expect(
      await screen.findByRole("heading", { name: /enter room code/i }),
    ).toBeInTheDocument()
  })
})
