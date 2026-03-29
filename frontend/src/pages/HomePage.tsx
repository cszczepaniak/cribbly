import { ChevronRight, LayoutGrid, ListOrdered, Trophy } from 'lucide-react'

import { RoomCodePanel } from '@/components/RoomCodePanel'
import { ReactLink } from '@/components/ReactLink'
import { useRoomAccess } from '@/contexts/roomAccessContext'

const destinationCards = [
  {
    to: '/divisions',
    title: 'Divisions',
    description: 'View your division and prelim games',
    Icon: LayoutGrid,
  },
  {
    to: '/standings',
    title: 'Standings',
    description: 'Scores and rankings',
    Icon: ListOrdered,
  },
  {
    to: '/tournament',
    title: 'Tournament',
    description: 'Tournament bracket',
    Icon: Trophy,
  },
] as const

function HomeLanding() {
  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:py-16">
      <header className="mb-14 text-center">
        <h1 className="text-4xl font-semibold tracking-tight text-foreground sm:text-5xl">
          Welcome to Cribbly
        </h1>
      </header>
      <section className="mb-4">
        <h2 className="text-xs font-medium tracking-wider text-muted-foreground uppercase">Go to</h2>
      </section>
      <div className="grid gap-4 sm:grid-cols-3">
        {destinationCards.map(({ to, title, description, Icon }) => (
          <ReactLink
            key={to}
            to={to}
            className="group bg-card text-card-foreground hover:border-border flex flex-col rounded-lg border shadow-sm transition-all hover:shadow-md focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          >
            <div className="flex flex-1 flex-col p-6">
              <div className="flex items-center justify-between gap-3">
                <span className="bg-primary/10 text-primary flex size-11 shrink-0 items-center justify-center rounded-lg">
                  <Icon
                    className="size-[22px] shrink-0"
                    strokeWidth={2}
                  />
                </span>
                <ChevronRight className="text-muted-foreground group-hover:text-foreground size-[18px] shrink-0 transition-colors" />
              </div>
              <h3 className="mt-4 text-lg leading-tight font-semibold">{title}</h3>
              <p className="text-muted-foreground mt-1.5 text-sm">{description}</p>
            </div>
          </ReactLink>
        ))}
      </div>
    </div>
  )
}

export function HomePage() {
  const { hasAccess, isLoading } = useRoomAccess()

  if (isLoading) {
    return (
      <main className="flex min-h-[calc(100vh-4.5rem)] flex-col items-center justify-center">
        <p className="text-muted-foreground text-sm">Loading...</p>
      </main>
    )
  }

  if (!hasAccess) {
    return (
      <main className="min-h-[calc(100vh-4.5rem)]">
        {hasAccess ? <HomeLanding/> : <RoomCodePanel />}
      </main>
    )
  }

  return (
    <main className="min-h-[calc(100vh-4.5rem)]">
      <HomeLanding />
    </main>
  )
}
