import { ReactLink } from '@/components/ReactLink'

export function HomePage() {
  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-2xl font-semibold tracking-tight">Cribbly</h1>
      <p className="text-sm text-muted-foreground">
        React shell with{' '}
        <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">?react=true</code>. The Go server serves this
        app when that query is present; otherwise you get the legacy UI for the same path.
      </p>
      <p className="text-sm text-muted-foreground">
        With the Go server, bundled assets are under{' '}
        <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">/app/</code>; the Vite dev server serves
        them from <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">/assets/</code> at the same route
        paths.
      </p>
      <p className="text-sm text-muted-foreground">
        Try the Connect room-code flow on{' '}
        <ReactLink to="/room" className="font-medium text-foreground underline-offset-4 hover:underline">
          room code
        </ReactLink>{' '}
        page (sets the same HttpOnly cookie as legacy POST <code className="font-mono text-xs">/room-code</code>).
      </p>
    </div>
  )
}
