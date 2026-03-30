import { fireEvent, render, screen, within } from "@testing-library/react"
import { createMemoryRouter, RouterProvider } from "react-router-dom"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { DEV_ADMIN_OVERRIDE_KEY } from "@/lib/devAdminOverride"
import { DEV_ROOM_ACCESS_OVERRIDE_KEY } from "@/lib/devRoomAccessOverride"
import { RootLayout } from "./RootLayout"

const roomCodeMocks = vi.hoisted(() => ({
  checkRoomAccess: vi.fn().mockResolvedValue({ hasAccess: true }),
}))

vi.mock("@/api/roomCodeClient", () => ({
  checkRoomAccess: roomCodeMocks.checkRoomAccess,
  setRoomCode: vi.fn(),
}))

describe("RootLayout", () => {
  beforeEach(() => {
    globalThis.localStorage?.removeItem?.(DEV_ROOM_ACCESS_OVERRIDE_KEY)
    globalThis.localStorage?.removeItem?.(DEV_ADMIN_OVERRIDE_KEY)
    roomCodeMocks.checkRoomAccess.mockResolvedValue({ hasAccess: true })
  })

  it("renders the logo and opens the menu with player and admin links", async () => {
    const router = createMemoryRouter(
      [
        {
          path: "/",
          element: <RootLayout />,
          children: [{ index: true, element: <div>Outlet child</div> }],
        },
      ],
      { initialEntries: ["/?react=true"] },
    )
    render(<RouterProvider router={router} />)

    expect(
      await screen.findByRole("link", { name: /^cribbly$/i }),
    ).toBeInTheDocument()
    expect(await screen.findByText("Outlet child")).toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: /open menu/i }))

    const playerNav = screen.getByRole("navigation", { name: /player pages/i })
    expect(playerNav).toBeInTheDocument()
    expect(
      within(playerNav).getByRole("link", { name: /^divisions$/i }),
    ).toHaveAttribute("href", "/divisions?react=true")
    expect(
      screen.getByRole("navigation", { name: /admin pages/i }),
    ).toBeInTheDocument()
    expect(screen.getByRole("link", { name: /^room codes$/i })).toHaveAttribute(
      "href",
      "/admin/room-codes?react=true",
    )
  })
})
