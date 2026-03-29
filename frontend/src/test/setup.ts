import "@testing-library/jest-dom/vitest"

/**
 * Some Node/Vitest setups expose a partial `localStorage` (e.g. missing `removeItem`).
 * Use an in-memory Storage so components and tests can rely on the full API.
 */
const memoryStore: Record<string, string> = {}
function createMemoryStorage(): Storage {
  return {
    get length() {
      return Object.keys(memoryStore).length
    },
    clear() {
      for (const k of Object.keys(memoryStore)) {
        delete memoryStore[k]
      }
    },
    getItem(key: string) {
      return Object.prototype.hasOwnProperty.call(memoryStore, key)
        ? memoryStore[key]
        : null
    },
    key(index: number) {
      return Object.keys(memoryStore)[index] ?? null
    },
    removeItem(key: string) {
      delete memoryStore[key]
    },
    setItem(key: string, value: string) {
      memoryStore[key] = String(value)
    },
  }
}

const ls = globalThis.localStorage
if (
  typeof ls === "undefined" ||
  typeof ls.getItem !== "function" ||
  typeof ls.setItem !== "function" ||
  typeof ls.removeItem !== "function"
) {
  globalThis.localStorage = createMemoryStorage()
}
