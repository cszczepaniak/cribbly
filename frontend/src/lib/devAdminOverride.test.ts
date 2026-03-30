import { beforeEach, describe, expect, it } from "vitest"
import {
  DEV_ADMIN_OVERRIDE_KEY,
  getDevAdminOverride,
  setDevAdminOverride,
} from "./devAdminOverride"

describe.skipIf(!import.meta.env.DEV)("devAdminOverride", () => {
  beforeEach(() => {
    localStorage.removeItem(DEV_ADMIN_OVERRIDE_KEY)
  })

  it("returns false when the key is unset", () => {
    expect(getDevAdminOverride()).toBe(false)
  })

  it("round-trips through localStorage", () => {
    setDevAdminOverride(true)
    expect(localStorage.getItem(DEV_ADMIN_OVERRIDE_KEY)).toBe("true")
    expect(getDevAdminOverride()).toBe(true)
    setDevAdminOverride(false)
    expect(localStorage.getItem(DEV_ADMIN_OVERRIDE_KEY)).toBeNull()
    expect(getDevAdminOverride()).toBe(false)
  })
})
