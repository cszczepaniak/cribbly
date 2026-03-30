import { render, screen } from "@testing-library/react"
import { MemoryRouter } from "react-router-dom"
import { describe, expect, it } from "vitest"
import { ReactLink } from "./ReactLink"

describe("ReactLink", () => {
  it("adds react=true to the link href", () => {
    render(
      <MemoryRouter>
        <ReactLink to="/divisions">Divisions</ReactLink>
      </MemoryRouter>,
    )
    expect(screen.getByRole("link", { name: /divisions/i })).toHaveAttribute(
      "href",
      "/divisions?react=true",
    )
  })

  it("preserves existing query params and sets react=true", () => {
    render(
      <MemoryRouter>
        <ReactLink to="/path?foo=bar">Go</ReactLink>
      </MemoryRouter>,
    )
    const href = screen
      .getByRole("link", { name: /^go$/i })
      .getAttribute("href")
    expect(href).toContain("foo=bar")
    expect(href).toContain("react=true")
  })
})
