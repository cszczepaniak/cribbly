import { useState } from "react"

import { generateAdminRoomCode } from "@/api/roomCodeClient"
import { useDevAdminOverride } from "@/lib/devAdminOverride"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

export function AdminRoomCodesPage() {
  const [code, setCode] = useState<string | null>(null)
  const [expiresAt, setExpiresAt] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)

  const [devAdminPretend] = useDevAdminOverride()
  const devSecretConfigured =
    import.meta.env.DEV &&
    typeof import.meta.env.VITE_DEV_ADMIN_SECRET === "string" &&
    import.meta.env.VITE_DEV_ADMIN_SECRET !== ""
  const devAdminApiReady = devSecretConfigured && devAdminPretend

  async function onGenerate() {
    setError(null)
    setPending(true)
    try {
      const res = await generateAdminRoomCode()
      setCode(res.code)
      setExpiresAt(
        new Date(res.expiresAt).toLocaleString(undefined, {
          dateStyle: "medium",
          timeStyle: "short",
        }),
      )
    } catch (e) {
      setCode(null)
      setExpiresAt(null)
      setError(e instanceof Error ? e.message : "Something went wrong")
    } finally {
      setPending(false)
    }
  }

  return (
    <div className="mx-auto flex w-full max-w-lg flex-col gap-6 p-4">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Room codes</h1>
        <p className="text-muted-foreground mt-1 text-sm">
          Generate a code visitors can use to unlock the site for 24 hours.
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
          <CardTitle>New code</CardTitle>
          <CardDescription>
            Creates a random 6-character code and stores it in the database.
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <Button
            type="button"
            onClick={() => void onGenerate()}
            disabled={
              pending ||
              (import.meta.env.DEV && devSecretConfigured && !devAdminApiReady)
            }
          >
            {pending ? "Generating…" : "Generate room code"}
          </Button>

          {error ? (
            <p className="text-destructive text-sm" role="alert">
              {error}
            </p>
          ) : null}

          {code ? (
            <div className="space-y-1 rounded-md border bg-card p-4">
              <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
                Code
              </p>
              <p className="font-mono text-2xl font-semibold tracking-widest">
                {code}
              </p>
              {expiresAt ? (
                <p className="text-muted-foreground mt-2 text-sm">
                  Expires {expiresAt}
                </p>
              ) : null}
            </div>
          ) : null}
        </CardContent>
      </Card>
    </div>
  )
}
