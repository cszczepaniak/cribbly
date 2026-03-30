import { Menu } from "lucide-react"

import { DevAdminAccessToggle } from "@/components/DevAdminAccessToggle"
import { DevRoomAccessToggle } from "@/components/DevRoomAccessToggle"
import { ReactLink } from "@/components/ReactLink"
import {
  RoomAccessOutlet,
  RoomAccessProvider,
} from "@/contexts/roomAccessContext"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"

function CribblyLogoLink({ className }: { className?: string }) {
  return (
    <ReactLink to="/" className={className}>
      Crib<span className="text-red-500">b</span>
      <span className="text-blue-600">l</span>
      <span className="text-green-500">y</span>
    </ReactLink>
  )
}

const playerLinks = [
  { label: "Divisions", to: "/divisions" },
  { label: "Standings", to: "/standings" },
  { label: "Tournament", to: "/tournament" },
] as const

const adminLinks = [
  { label: "Players", to: "/admin/players" },
  { label: "Teams", to: "/admin/teams" },
  { label: "Divisions", to: "/admin/divisions" },
  { label: "Games", to: "/admin/games" },
  { label: "Users", to: "/admin/users" },
  { label: "My Profile", to: "/admin/profile" },
  { label: "Room Codes", to: "/admin/room-codes" },
] as const

export function RootLayout() {
  return (
    <RoomAccessProvider>
      <Sheet>
        <div className="flex h-dvh max-h-dvh max-w-screen flex-col overflow-hidden">
          <header className="bg-primary text-background sticky top-0 z-50 flex flex-shrink-0 flex-row items-center justify-between px-6 py-4 text-2xl font-semibold tracking-wide">
            <CribblyLogoLink className="text-background hover:opacity-90" />
            <SheetTrigger asChild>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="text-background hover:bg-background/15 hover:text-background [&_svg]:size-6"
                aria-label="Open menu"
              >
                <Menu />
              </Button>
            </SheetTrigger>
          </header>
          <div className="flex min-h-0 min-w-0 flex-1 flex-col overflow-x-hidden overflow-y-auto bg-muted/30">
            <RoomAccessOutlet />
          </div>
        </div>
        <SheetContent
          side="top"
          showCloseButton
          className="border-primary bg-primary text-background [&_[data-slot=sheet-close]_button]:text-background [&_[data-slot=sheet-close]_button]:hover:bg-background/15"
        >
          <SheetHeader>
            <SheetTitle className="sr-only">Menu</SheetTitle>
          </SheetHeader>
          <div className="flex flex-col space-y-6 px-4 py-2 text-lg">
            <nav aria-label="Player pages">
              <div className="mb-4 space-y-2 font-semibold">
                <p>Player Pages</p>
                <Separator className="bg-background/30" />
              </div>
              <ul className="list-none space-y-2">
                {playerLinks.map(({ label, to }) => (
                  <li key={to}>
                    <ReactLink to={to} className="hover:underline">
                      {label}
                    </ReactLink>
                  </li>
                ))}
              </ul>
            </nav>
            <nav aria-label="Admin pages">
              <div className="mb-4 space-y-2 font-semibold">
                <p>Admin Pages</p>
                <Separator className="bg-background/30" />
              </div>
              <ul className="list-none space-y-2">
                {adminLinks.map(({ label, to }) => (
                  <li key={to}>
                    <ReactLink to={to} className="hover:underline">
                      {label}
                    </ReactLink>
                  </li>
                ))}
                <li>
                  <button
                    type="button"
                    className="cursor-pointer hover:underline"
                  >
                    Logout
                  </button>
                </li>
                <li>
                  <ReactLink to="/admin/login" className="hover:underline">
                    Login
                  </ReactLink>
                </li>
              </ul>
            </nav>
          </div>
          <DevRoomAccessToggle />
          <DevAdminAccessToggle />
          <SheetFooter />
        </SheetContent>
      </Sheet>
    </RoomAccessProvider>
  )
}
