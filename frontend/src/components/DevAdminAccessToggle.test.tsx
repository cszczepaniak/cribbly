import { fireEvent, render, screen } from "@testing-library/react"
import { beforeEach, describe, expect, it } from "vitest"
import { DevAdminAccessToggle } from "./DevAdminAccessToggle"
import { DEV_ADMIN_OVERRIDE_KEY } from "@/lib/devAdminOverride"

describe.skipIf(!import.meta.env.DEV)("DevAdminAccessToggle", () => {
  beforeEach(() => {
    localStorage.removeItem(DEV_ADMIN_OVERRIDE_KEY)
  })

  it("renders the label and switch", () => {
    render(<DevAdminAccessToggle />)
    expect(screen.getByText("Pretend logged-in admin")).toBeInTheDocument()
    expect(
      screen.getByRole("switch", {
        name: /pretend logged-in admin \(dev only\)/i,
      }),
    ).toBeInTheDocument()
  })

  it("persists the override in localStorage when turned on", () => {
    render(<DevAdminAccessToggle />)
    fireEvent.click(
      screen.getByRole("switch", {
        name: /pretend logged-in admin \(dev only\)/i,
      }),
    )
    expect(localStorage.getItem(DEV_ADMIN_OVERRIDE_KEY)).toBe("true")
    expect(
      screen.getByRole("switch", {
        name: /pretend logged-in admin \(dev only\)/i,
      }),
    ).toBeChecked()
  })

  it("clears localStorage when turned off", () => {
    localStorage.setItem(DEV_ADMIN_OVERRIDE_KEY, "true")
    render(<DevAdminAccessToggle />)
    expect(
      screen.getByRole("switch", {
        name: /pretend logged-in admin \(dev only\)/i,
      }),
    ).toBeChecked()
    fireEvent.click(
      screen.getByRole("switch", {
        name: /pretend logged-in admin \(dev only\)/i,
      }),
    )
    expect(localStorage.getItem(DEV_ADMIN_OVERRIDE_KEY)).toBeNull()
    expect(
      screen.getByRole("switch", {
        name: /pretend logged-in admin \(dev only\)/i,
      }),
    ).not.toBeChecked()
  })
})
