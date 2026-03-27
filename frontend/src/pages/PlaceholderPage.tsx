import { useLocation } from 'react-router-dom'

type Props = {
  title?: string
}

/** Placeholder for routes you migrate from the Go/Templ app. */
export function PlaceholderPage({ title }: Props) {
  const { pathname, search } = useLocation()
  return (
    <div className="flex flex-col gap-2">
      <h1 className="text-xl font-semibold tracking-tight">
        {title ?? 'Page (not migrated yet)'}
      </h1>
      <p className="break-all font-mono text-sm text-zinc-500 dark:text-zinc-400">
        {pathname}
        {search}
      </p>
    </div>
  )
}
