import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { DEV_ROOM_ACCESS_OVERRIDE_KEY, useDevRoomAccessOverride } from '@/lib/devRoomAccessOverride'

/** Dev-only menu row: pretend the SPA has room access (skips probe + `/` redirect). */
export function DevRoomAccessToggle() {
  const [on, setOn] = useDevRoomAccessOverride()

  if (!import.meta.env.DEV) {
    return null
  }

  return (
    <div className="border-background/20 border-t px-4 py-3">
      <div className="flex items-center justify-between gap-3">
        <Label
          htmlFor="dev-room-access-override"
          className="text-background/90 text-sm font-normal"
        >
          Pretend room access
        </Label>
        <Switch
          id="dev-room-access-override"
          checked={on}
          onCheckedChange={setOn}
          aria-label="Pretend room access (dev only)"
        />
      </div>
      <p className="text-background/55 mt-1.5 text-xs">
        Skips the access probe and home redirect. Persists in localStorage (
        <code className="font-mono text-[0.65rem]">{DEV_ROOM_ACCESS_OVERRIDE_KEY}</code>
        ).
      </p>
    </div>
  )
}
