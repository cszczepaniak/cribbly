export function HomePage() {
  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-2xl font-semibold tracking-tight">Cribbly</h1>
      <p className="text-sm text-zinc-600 dark:text-zinc-400">
        React shell with <code className="rounded bg-zinc-200 px-1.5 py-0.5 font-mono text-xs dark:bg-zinc-800">?react=true</code>. The Go
        server serves this app when that query is present; otherwise you get the legacy UI for the same path.
      </p>
      <p className="text-sm text-zinc-600 dark:text-zinc-400">
        With the Go server, bundled assets are under{' '}
        <code className="rounded bg-zinc-200 px-1.5 py-0.5 font-mono text-xs dark:bg-zinc-800">/app/</code>
        ; the Vite dev server serves them from <code className="rounded bg-zinc-200 px-1.5 py-0.5 font-mono text-xs dark:bg-zinc-800">/assets/</code> at the same route paths.
      </p>
    </div>
  )
}
