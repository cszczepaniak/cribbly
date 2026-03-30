import { getDevAdminOverride } from "@/lib/devAdminOverride"

/** Same name as internal/server/middleware.DevAdminHeader */
export const DEV_ADMIN_HEADER = "X-Cribbly-Dev-Admin"

/**
 * Applies the dev-admin bypass header when dev mode, VITE_DEV_ADMIN_SECRET is set, and the user
 * turned on "Pretend logged-in admin" (see DevAdminAccessToggle). Must match CRIBBLY_DEV_ADMIN_SECRET
 * on the server; ignored in production.
 */
export function applyDevAdminHeaders(headers: Headers): void {
  const secret = import.meta.env.VITE_DEV_ADMIN_SECRET
  if (
    !import.meta.env.DEV ||
    typeof secret !== "string" ||
    secret === "" ||
    !getDevAdminOverride()
  ) {
    return
  }
  headers.set(DEV_ADMIN_HEADER, secret)
}
