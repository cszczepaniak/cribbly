import { Code, ConnectError } from '@connectrpc/connect'
import { useId, useState, type FormEvent } from 'react'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { setRoomCode } from '@/api/roomCodeClient'
import { useReactNavigate } from '@/hooks/useReactNavigate'

export function RoomCodePage() {
  const id = useId()
  const navigate = useReactNavigate()

  const [value, setValue] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function onSubmit(e: FormEvent) {
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
    <Card className="mx-auto w-full max-w-md">
      <CardHeader>
        <CardTitle>Enter room code</CardTitle>
        <CardDescription>
          Submit your room code to set an HttpOnly cookie in the Go server (durable across future requests).
        </CardDescription>
      </CardHeader>

      <form onSubmit={onSubmit}>
        <CardContent className="space-y-4">
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
            <p className="text-sm text-destructive" role="alert">
              {error}
            </p>
          ) : null}
        </CardContent>

        <CardFooter>
          <Button type="submit" className="w-full" disabled={loading}>
            {loading ? 'Checking…' : 'Continue'}
          </Button>
        </CardFooter>
      </form>
    </Card>
  )
}

