import { create } from "@bufbuild/protobuf"
import {
  fireEvent,
  render,
  screen,
  waitFor,
  within,
} from "@testing-library/react"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { AdminPlayersPage } from "./AdminPlayersPage"
import { PlayerSchema, type Player } from "@/gen/cribbly/v1/players_pb"

const mocks = vi.hoisted(() => ({
  listPlayers: vi.fn(),
  createPlayer: vi.fn(),
  updatePlayer: vi.fn(),
  deletePlayer: vi.fn(),
  deleteAllPlayers: vi.fn(),
  generateRandomPlayers: vi.fn(),
}))

vi.mock("@/api/playersClient", () => mocks)

vi.mock("@/lib/devAdminOverride", () => ({
  useDevAdminOverride: () => [true, vi.fn()],
}))

function player(overrides: Partial<Player> = {}): Player {
  return create(PlayerSchema, {
    id: "p1",
    firstName: "Ada",
    lastName: "Lovelace",
    teamId: "",
    ...overrides,
  })
}

function addPlayerForm() {
  const title = screen.getByText("Add a player")
  const card = title.closest("[data-slot=card]")
  if (!(card instanceof HTMLElement)) {
    throw new Error("expected Add a player card")
  }
  return within(card)
}

describe("AdminPlayersPage", () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mocks.listPlayers.mockResolvedValue({ players: [] })
    mocks.createPlayer.mockResolvedValue({ players: [player()] })
    mocks.updatePlayer.mockResolvedValue({
      players: [player({ firstName: "Augusta" })],
    })
    mocks.deletePlayer.mockResolvedValue({ players: [] })
    mocks.deleteAllPlayers.mockResolvedValue({})
    mocks.generateRandomPlayers.mockResolvedValue({ players: [player()] })
  })

  it("loads players on mount", async () => {
    mocks.listPlayers.mockResolvedValue({
      players: [
        player({ id: "a", firstName: "Ada", lastName: "Lovelace" }),
        player({ id: "b", firstName: "Alan", lastName: "Turing" }),
      ],
    })

    render(<AdminPlayersPage />)

    await waitFor(() => {
      expect(mocks.listPlayers).toHaveBeenCalledTimes(1)
    })

    expect(
      await screen.findByRole("heading", { name: /^2 players$/ }),
    ).toBeInTheDocument()
    expect(screen.getByText("Ada Lovelace")).toBeInTheDocument()
    expect(screen.getByText("Alan Turing")).toBeInTheDocument()
  })

  it("shows an error when listing players fails", async () => {
    mocks.listPlayers.mockRejectedValue(new Error("network failed"))

    render(<AdminPlayersPage />)

    expect(await screen.findByText("network failed")).toBeInTheDocument()
  })

  it("submits trimmed names when adding a player", async () => {
    render(<AdminPlayersPage />)

    await screen.findByText(/no players registered yet/i)

    const form = addPlayerForm()
    fireEvent.change(form.getByLabelText(/^first name$/i), {
      target: { value: "  Grace  " },
    })
    fireEvent.change(form.getByLabelText(/^last name$/i), {
      target: { value: " Hopper " },
    })
    fireEvent.click(form.getByRole("button", { name: /^add player$/i }))

    await waitFor(() => {
      expect(mocks.createPlayer).toHaveBeenCalledWith("Grace", "Hopper")
    })
    expect(await screen.findByText("Ada Lovelace")).toBeInTheDocument()
  })

  it("updates a player name from the row editor", async () => {
    mocks.listPlayers.mockResolvedValue({
      players: [player({ id: "p1", firstName: "Ada", lastName: "Lovelace" })],
    })

    render(<AdminPlayersPage />)

    await screen.findByText("Ada Lovelace")

    fireEvent.click(screen.getByRole("button", { name: /edit ada lovelace/i }))

    const first = screen.getByDisplayValue("Ada")
    fireEvent.change(first, { target: { value: "Augusta" } })
    fireEvent.click(screen.getByRole("button", { name: /save name/i }))

    await waitFor(() => {
      expect(mocks.updatePlayer).toHaveBeenCalledWith(
        "p1",
        "Augusta",
        "Lovelace",
      )
    })
    expect(await screen.findByText("Augusta Lovelace")).toBeInTheDocument()
  })

  it("cancels inline edit without saving", async () => {
    mocks.listPlayers.mockResolvedValue({
      players: [player({ id: "p1", firstName: "Ada", lastName: "Lovelace" })],
    })

    render(<AdminPlayersPage />)

    await screen.findByText("Ada Lovelace")

    fireEvent.click(screen.getByRole("button", { name: /edit ada lovelace/i }))
    fireEvent.change(screen.getByDisplayValue("Ada"), {
      target: { value: "X" },
    })
    fireEvent.click(screen.getByRole("button", { name: /cancel editing/i }))

    expect(mocks.updatePlayer).not.toHaveBeenCalled()
    expect(screen.getByText("Ada Lovelace")).toBeInTheDocument()
  })

  it("deletes a player when trash is clicked", async () => {
    mocks.listPlayers.mockResolvedValue({
      players: [player({ id: "p1", firstName: "Ada", lastName: "Lovelace" })],
    })

    render(<AdminPlayersPage />)

    await screen.findByText("Ada Lovelace")

    fireEvent.click(
      screen.getByRole("button", { name: /delete ada lovelace/i }),
    )

    await waitFor(() => {
      expect(mocks.deletePlayer).toHaveBeenCalledWith("p1")
    })
    expect(
      await screen.findByText(/no players registered yet/i),
    ).toBeInTheDocument()
  })

  it("deletes all after confirm", async () => {
    const confirmSpy = vi.spyOn(window, "confirm").mockReturnValue(true)
    mocks.listPlayers.mockResolvedValue({
      players: [player()],
    })

    render(<AdminPlayersPage />)

    await screen.findByText("Ada Lovelace")

    fireEvent.click(screen.getByText("Development tools"))
    fireEvent.click(screen.getByRole("button", { name: /delete all players/i }))

    await waitFor(() => {
      expect(mocks.deleteAllPlayers).toHaveBeenCalledTimes(1)
    })
    expect(confirmSpy).toHaveBeenCalledWith(
      "Delete all players? Players on teams will be unassigned first.",
    )
    confirmSpy.mockRestore()
  })
})
