import { Outlet } from 'react-router-dom'
import { ReactLink } from '../components/ReactLink'

export function RootLayout() {
  return (
    <div className="min-h-screen bg-zinc-50 text-zinc-900 dark:bg-zinc-950 dark:text-zinc-50">
      <header className="border-b border-zinc-200 bg-white px-4 py-3 dark:border-zinc-800 dark:bg-zinc-900">
        <nav className="mx-auto flex max-w-3xl flex-wrap items-center gap-4 text-sm">
          <ReactLink
            to="/"
            className="font-medium text-zinc-700 hover:text-zinc-950 dark:text-zinc-300 dark:hover:text-white"
          >
            Home
          </ReactLink>
          <ReactLink
            to="/admin/games"
            className="text-zinc-600 hover:text-zinc-950 dark:text-zinc-400 dark:hover:text-white"
          >
            Admin games
          </ReactLink>
          <span className="text-zinc-400 dark:text-zinc-600">(sample links)</span>
        </nav>
      </header>
      <main className="mx-auto max-w-3xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
