import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import {
  DEV_ADMIN_OVERRIDE_KEY,
  useDevAdminOverride,
} from "@/lib/devAdminOverride"

/** Dev-only menu row: send X-Cribbly-Dev-Admin so the Go server treats you as an admin (with matching secret). */
export function DevAdminAccessToggle() {
  const [on, setOn] = useDevAdminOverride()

  if (!import.meta.env.DEV) {
    return null
  }

  return (
    <div className="border-background/20 border-t px-4 py-3">
      <div className="flex items-center justify-between gap-3">
        <Label
          htmlFor="dev-admin-access-override"
          className="text-background/90 text-sm font-normal"
        >
          Pretend logged-in admin
        </Label>
        <Switch
          id="dev-admin-access-override"
          checked={on}
          onCheckedChange={setOn}
          aria-label="Pretend logged-in admin (dev only)"
        />
      </div>
      <p className="text-background/55 mt-1.5 text-xs">
        Sends the dev admin header on API calls when{" "}
        <code className="font-mono text-[0.65rem]">VITE_DEV_ADMIN_SECRET</code>{" "}
        is set. Persists in localStorage (
        <code className="font-mono text-[0.65rem]">
          {DEV_ADMIN_OVERRIDE_KEY}
        </code>
        ).
      </p>
    </div>
  )
}
