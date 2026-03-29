import { Code, ConnectError } from '@connectrpc/connect'
import { useId, useState, type SubmitEvent } from 'react'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { setRoomCode } from '@/api/roomCodeClient'
import { useReactNavigate } from '@/hooks/useReactNavigate'

export function RoomCodePanel() {
  const id = useId()
  const navigate = useReactNavigate()

  const [value, setValue] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function onSubmit(e: SubmitEvent<HTMLFormElement>) {
    e.preventDefault()
    setError(null)
    setLoading(true)

    try {
      await setRoomCode(value.trim())
      navigate('/')
    } catch (err) {
      if (err instanceof ConnectError && err.code === Code.InvalidArgument) {
        setError('That code is not valid or has expired.')
      } else if (err instanceof ConnectError) {
        setError(err.message || 'Something went wrong.')
      } else {
        setError('Something went wrong.')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto max-w-md px-4 py-16 sm:py-24">
      <Card className="border-0 shadow-lg">
        <CardHeader className="p-8 pb-2 text-center">
          <h1 className="text-2xl font-semibold text-foreground">Enter room code</h1>
          <CardDescription className="mt-2">
            Ask your organizer for the code to view divisions, standings, and the tournament.
          </CardDescription>
        </CardHeader>
        <CardContent className="px-8 pt-4 pb-8">
          <form
            onSubmit={onSubmit}
            className="flex flex-col gap-4"
          >
            <div className="space-y-2">
              <Label htmlFor={id}>Room code</Label>
              <Input
                id={id}
                name="room_code"
                autoComplete="off"
                autoCapitalize="characters"
                autoCorrect="off"
                spellCheck={false}
                placeholder="e.g. ABC123"
                value={value}
                onChange={(ev) => setValue(ev.target.value)}
                disabled={loading}
              />
            </div>

            {error ? (
              <p
                className="text-sm text-destructive"
                role="alert"
              >
                {error}
              </p>
            ) : null}

            <Button
              type="submit"
              className="w-full"
              disabled={loading}
            >
              {loading ? 'Checking...' : 'Continue'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
