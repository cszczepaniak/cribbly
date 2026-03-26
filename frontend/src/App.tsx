function App() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-zinc-50 p-8 text-zinc-900 dark:bg-zinc-950 dark:text-zinc-50">
      <h1 className="text-2xl font-semibold tracking-tight">Cribbly</h1>
      <p className="max-w-md text-center text-sm text-zinc-600 dark:text-zinc-400">
        React frontend (Vite + Tailwind). Served by the Go app at{' '}
        <code className="rounded bg-zinc-200 px-1.5 py-0.5 font-mono text-xs dark:bg-zinc-800">
          /app/
        </code>
        .
      </p>
    </div>
  )
}

export default App
