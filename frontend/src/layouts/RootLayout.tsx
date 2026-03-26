import { Outlet } from 'react-router-dom'
import { ReactLink } from '../components/ReactLink'

export function RootLayout() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="border-b border-border bg-card px-4 py-3">
        <nav className="mx-auto flex max-w-3xl flex-wrap items-center gap-4 text-sm">
          <ReactLink
            to="/"
            className="font-medium text-foreground hover:underline"
          >
            Home
          </ReactLink>
          <ReactLink to="/room" className="text-muted-foreground hover:text-foreground">
            Room code
          </ReactLink>
          <ReactLink
            to="/admin/games"
            className="text-muted-foreground hover:text-foreground"
          >
            Admin games
          </ReactLink>
          <span className="text-muted-foreground">(sample links)</span>
        </nav>
      </header>
      <main className="mx-auto max-w-3xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
