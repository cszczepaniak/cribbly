/** Query param the Go server checks to serve the React shell instead of legacy HTML. */
export const REACT_QUERY = "react" as const

/**
 * Appends `react=true` to a path, preserving existing query params.
 * Use with {@link ReactLink} / {@link useReactNavigate} so in-app navigation
 * keeps the React shell when running behind the Go server.
 */
export function withReactQuery(path: string): string {
  const q = path.indexOf("?")
  const pathname = q >= 0 ? path.slice(0, q) : path
  const search = q >= 0 ? path.slice(q + 1) : ""
  const params = new URLSearchParams(search)
  params.set(REACT_QUERY, "true")
  const qs = params.toString()
  return `${pathname}?${qs}`
}
