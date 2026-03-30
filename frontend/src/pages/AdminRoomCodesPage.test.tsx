import { fireEvent, render, screen, waitFor } from "@testing-library/react"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { AdminRoomCodesPage } from "./AdminRoomCodesPage"

const mocks = vi.hoisted(() => ({
  generateAdminRoomCode: vi.fn(),
}))

vi.mock("@/api/roomCodeClient", () => ({
  generateAdminRoomCode: mocks.generateAdminRoomCode,
}))

vi.mock("@/lib/devAdminOverride", () => ({
  useDevAdminOverride: () => [true, vi.fn()],
}))

describe("AdminRoomCodesPage", () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mocks.generateAdminRoomCode.mockResolvedValue({
      code: "XYZ789",
      expiresAt: "2026-01-15T12:00:00.000Z",
    })
  })

  it("generates and displays a room code", async () => {
    render(<AdminRoomCodesPage />)

    fireEvent.click(screen.getByRole("button", { name: /generate room code/i }))

    await waitFor(() => {
      expect(mocks.generateAdminRoomCode).toHaveBeenCalledTimes(1)
    })

    expect(screen.getByText("XYZ789")).toBeInTheDocument()
    expect(screen.getByText(/expires/i)).toBeInTheDocument()
  })

  it("shows an alert when generation fails", async () => {
    mocks.generateAdminRoomCode.mockRejectedValue(new Error("rpc failed"))

    render(<AdminRoomCodesPage />)

    fireEvent.click(screen.getByRole("button", { name: /generate room code/i }))

    expect(await screen.findByRole("alert")).toHaveTextContent("rpc failed")
  })
})
