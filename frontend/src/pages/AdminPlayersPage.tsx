import { useCallback, useEffect, useState } from "react"
import { Check, Pencil, Trash2, X } from "lucide-react"

import {
  createPlayer,
  deleteAllPlayers,
  deletePlayer,
  generateRandomPlayers,
  listPlayers,
  updatePlayer,
} from "@/api/playersClient"
import { useDevAdminOverride } from "@/lib/devAdminOverride"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import type { Player } from "@/gen/cribbly/v1/players_pb"

function displayName(p: Player) {
  return `${p.firstName} ${p.lastName}`.trim()
}

export function AdminPlayersPage() {
  const [players, setPlayers] = useState<Player[]>([])
  const [firstName, setFirstName] = useState("")
  const [lastName, setLastName] = useState("")
  const [randomCount, setRandomCount] = useState("10")
  const [loading, setLoading] = useState(true)
  const [pending, setPending] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editFirstName, setEditFirstName] = useState("")
  const [editLastName, setEditLastName] = useState("")

  const [devAdminPretend] = useDevAdminOverride()
  const devSecretConfigured =
    import.meta.env.DEV &&
    typeof import.meta.env.VITE_DEV_ADMIN_SECRET === "string" &&
    import.meta.env.VITE_DEV_ADMIN_SECRET !== ""
  const devAdminApiReady = devSecretConfigured && devAdminPretend

  const apiBlocked =
    pending || (import.meta.env.DEV && devSecretConfigured && !devAdminApiReady)

  const refresh = useCallback(async () => {
    setError(null)
    const res = await listPlayers()
    setPlayers(res.players)
  }, [])

  useEffect(() => {
    let cancelled = false
    setLoading(true)
    void (async () => {
      try {
        await refresh()
      } catch (e) {
        if (!cancelled) {
          setError(e instanceof Error ? e.message : "Failed to load players")
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    })()
    return () => {
      cancelled = true
    }
  }, [refresh])

  async function run<T>(fn: () => Promise<T>) {
    setError(null)
    setPending(true)
    try {
      await fn()
    } catch (e) {
      setError(e instanceof Error ? e.message : "Request failed")
    } finally {
      setPending(false)
    }
  }

  async function onAdd(e: React.FormEvent) {
    e.preventDefault()
    await run(async () => {
      const res = await createPlayer(firstName.trim(), lastName.trim())
      setPlayers(res.players)
      setFirstName("")
      setLastName("")
    })
  }

  async function onDelete(id: string) {
    await run(async () => {
      const res = await deletePlayer(id)
      setPlayers(res.players)
    })
  }

  async function onDeleteAll() {
    if (
      !window.confirm(
        "Delete all players? Players on teams will be unassigned first.",
      )
    ) {
      return
    }
    await run(async () => {
      await deleteAllPlayers()
      setPlayers([])
    })
  }

  async function onGenerateRandom() {
    const n = Number.parseInt(randomCount, 10)
    await run(async () => {
      const res = await generateRandomPlayers(Number.isFinite(n) ? n : 10)
      setPlayers(res.players)
    })
  }

  function beginEdit(p: Player) {
    setEditingId(p.id)
    setEditFirstName(p.firstName)
    setEditLastName(p.lastName)
  }

  function cancelEdit() {
    setEditingId(null)
    setEditFirstName("")
    setEditLastName("")
  }

  async function onSaveEdit(e: React.FormEvent) {
    e.preventDefault()
    if (!editingId) {
      return
    }
    await run(async () => {
      const res = await updatePlayer(
        editingId,
        editFirstName.trim(),
        editLastName.trim(),
      )
      setPlayers(res.players)
      cancelEdit()
    })
  }

  useEffect(() => {
    if (!editingId) {
      return
    }
    function onKeyDown(ev: KeyboardEvent) {
      if (ev.key === "Escape") {
        setEditingId(null)
        setEditFirstName("")
        setEditLastName("")
      }
    }
    window.addEventListener("keydown", onKeyDown)
    return () => window.removeEventListener("keydown", onKeyDown)
  }, [editingId])

  const countLabel =
    players.length === 0
      ? "No players yet"
      : players.length === 1
        ? "1 player"
        : `${players.length} players`

  return (
    <div className="mx-auto flex w-full max-w-3xl flex-col gap-6 p-4">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">
          Player registration
        </h1>
        <p className="text-muted-foreground mt-1 text-sm">
          Add people who can be paired into teams. Names appear as &quot;First
          Last&quot; everywhere else in the app.
        </p>
      </div>

      {import.meta.env.DEV && !devSecretConfigured ? (
        <p className="text-amber-800 dark:text-amber-200 bg-amber-50 dark:bg-amber-950/40 rounded-md border border-amber-200/80 p-3 text-sm dark:border-amber-800/60">
          Set <span className="font-mono">CRIBBLY_DEV_ADMIN_SECRET</span> on the
          Go server and <span className="font-mono">VITE_DEV_ADMIN_SECRET</span>{" "}
          in <span className="font-mono">frontend/.env.development.local</span>{" "}
          to the same value so API calls can bypass admin login in local dev.
        </p>
      ) : null}
      {import.meta.env.DEV && devSecretConfigured && !devAdminPretend ? (
        <p className="text-amber-800 dark:text-amber-200 bg-amber-50 dark:bg-amber-950/40 rounded-md border border-amber-200/80 p-3 text-sm dark:border-amber-800/60">
          Turn on <span className="font-medium">Pretend logged-in admin</span>{" "}
          in the menu so dev API calls include the admin bypass header.
        </p>
      ) : null}

      <Card>
        <CardHeader>
          <CardTitle>Add a player</CardTitle>
          <CardDescription>
            Enter first and last name, then submit. The list below updates
            immediately.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form
            onSubmit={(e) => void onAdd(e)}
            className="flex flex-col gap-4 sm:flex-row sm:flex-wrap sm:items-end"
          >
            <div className="grid min-w-0 flex-1 gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="player-first">First name</Label>
                <Input
                  id="player-first"
                  name="firstName"
                  autoComplete="given-name"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  disabled={pending}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="player-last">Last name</Label>
                <Input
                  id="player-last"
                  name="lastName"
                  autoComplete="family-name"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  disabled={pending}
                />
              </div>
            </div>
            <Button
              type="submit"
              className="w-full shrink-0 sm:w-auto"
              disabled={apiBlocked}
            >
              {pending ? "Saving…" : "Add player"}
            </Button>
          </form>
        </CardContent>
      </Card>

      <div className="space-y-2">
        <div className="flex items-baseline justify-between gap-2">
          <h2 className="text-base font-medium">{countLabel}</h2>
        </div>
        <div className="bg-card max-h-[min(60vh,28rem)] overflow-auto rounded-xl ring-1 ring-foreground/10">
          {loading ? (
            <p className="text-muted-foreground p-4 text-sm">Loading…</p>
          ) : players.length === 0 ? (
            <p className="text-muted-foreground p-4 text-sm">
              No players registered yet.
            </p>
          ) : (
            <table className="w-full border-collapse text-sm">
              <thead>
                <tr className="bg-muted/50 border-b text-left">
                  <th className="text-muted-foreground px-4 py-2 font-medium">
                    Name
                  </th>
                  <th className="text-muted-foreground w-24 px-2 py-2 text-right font-medium">
                    <span className="sr-only">Edit or remove</span>
                  </th>
                </tr>
              </thead>
              <tbody>
                {players.map((p) => {
                  const onTeam = Boolean(p.teamId)
                  const isEditing = editingId === p.id
                  return (
                    <tr
                      key={p.id}
                      className="border-b border-foreground/5 last:border-0"
                    >
                      <td className="px-4 py-2.5">
                        {isEditing ? (
                          <form
                            className="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-end"
                            onSubmit={(e) => void onSaveEdit(e)}
                          >
                            <div className="grid min-w-0 flex-1 gap-3 sm:grid-cols-2">
                              <div className="space-y-1.5">
                                <Label
                                  htmlFor={`edit-first-${p.id}`}
                                  className="text-xs"
                                >
                                  First name
                                </Label>
                                <Input
                                  id={`edit-first-${p.id}`}
                                  value={editFirstName}
                                  onChange={(e) =>
                                    setEditFirstName(e.target.value)
                                  }
                                  disabled={pending}
                                  autoComplete="off"
                                />
                              </div>
                              <div className="space-y-1.5">
                                <Label
                                  htmlFor={`edit-last-${p.id}`}
                                  className="text-xs"
                                >
                                  Last name
                                </Label>
                                <Input
                                  id={`edit-last-${p.id}`}
                                  value={editLastName}
                                  onChange={(e) =>
                                    setEditLastName(e.target.value)
                                  }
                                  disabled={pending}
                                  autoComplete="off"
                                />
                              </div>
                            </div>
                            <div className="flex shrink-0 gap-1">
                              <Button
                                type="submit"
                                size="icon"
                                disabled={apiBlocked}
                                aria-label="Save name"
                              >
                                <Check className="size-4" />
                              </Button>
                              <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                disabled={pending}
                                onClick={() => cancelEdit()}
                                aria-label="Cancel editing"
                              >
                                <X className="size-4" />
                              </Button>
                            </div>
                          </form>
                        ) : (
                          <span>{displayName(p)}</span>
                        )}
                      </td>
                      <td className="px-2 py-1 text-right whitespace-nowrap">
                        {isEditing ? null : (
                          <span className="inline-flex items-center justify-end gap-0.5">
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              className="text-muted-foreground hover:text-foreground"
                              disabled={apiBlocked}
                              title="Edit name"
                              onClick={() => beginEdit(p)}
                              aria-label={`Edit ${displayName(p)}`}
                            >
                              <Pencil className="size-4" />
                            </Button>
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              className="text-muted-foreground hover:text-destructive"
                              disabled={apiBlocked || onTeam}
                              title={
                                onTeam
                                  ? "Remove from team before deleting"
                                  : "Delete player"
                              }
                              onClick={() => void onDelete(p.id)}
                              aria-label={`Delete ${displayName(p)}`}
                            >
                              <Trash2 className="size-4" />
                            </Button>
                          </span>
                        )}
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          )}
        </div>
      </div>

      {error ? (
        <p className="text-destructive text-sm" role="alert">
          {error}
        </p>
      ) : null}

      <details className="group rounded-xl ring-1 ring-foreground/10">
        <summary className="text-muted-foreground hover:bg-muted/40 cursor-pointer list-none rounded-xl px-4 py-3 text-sm font-medium transition-colors [&::-webkit-details-marker]:hidden">
          <span className="inline-flex items-center gap-2">
            <span aria-hidden>▸</span>
            Development tools
          </span>
        </summary>
        <div className="border-t px-4 py-4">
          <p className="text-muted-foreground mb-4 text-sm">
            Random names and bulk delete are for local testing only.
          </p>
          <div className="flex flex-col gap-6 sm:flex-row sm:items-end">
            <div className="space-y-2">
              <Label htmlFor="random-count">Number of random players</Label>
              <Input
                id="random-count"
                type="number"
                min={1}
                max={500}
                value={randomCount}
                onChange={(e) => setRandomCount(e.target.value)}
                disabled={pending}
                className="max-w-[8rem]"
              />
            </div>
            <Button
              type="button"
              variant="secondary"
              disabled={apiBlocked}
              onClick={() => void onGenerateRandom()}
            >
              Generate
            </Button>
          </div>
          <Separator className="my-6" />
          <Button
            type="button"
            variant="destructive"
            disabled={apiBlocked || players.length === 0}
            onClick={() => void onDeleteAll()}
          >
            Delete all players
          </Button>
        </div>
      </details>
    </div>
  )
}
