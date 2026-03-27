# Frontend in React

Until now, Go+datastar has served me quite well. However there are a few key areas it's falling
short:
1. It's hard to test UI interactions at all. Doing so basically requires headless chrome. Can't unit
   test it at all.
2. The developer tooling isn't good. We have live reload _kind of_, but it's slow and buggy, and
   tailwind live reload has never worked.
3. Some of the things I want to do are fairly interactive and I just think I may need more
   javascript (richer animations, e.g.)

## New Architecture

### Frontend

- I want to use React for client-side ONLY (no SSR, no SSG). It's important to me that routing is
  easy, but I don't have a particular React framework in mind. Maybe vanilla React is enough.
- I want to continue using tailwind
- I want to use vite
- I want to use shadcn for components (https://ui.shadcn.com/)
- To keep things simple, I'd like to bundle and serve the frontend as static files from the Go
  backend. This means I need to build developer tooling/scripts/etc. to bundle things up and then
  I'll use //go:embed to embed them into the Go binary.

It's also important to me that we add tests to the frontend interactions as we go!

### Backend

- Aside from UI and API, this will stay the same.
- For backend RPCs (to be consumed by the frontend), I want to use Connect RPC (https://connectrpc.com/).
- I'll need a replacement for the server sent events. Ideally I could use a similar proto-based IDL
  to define the type of messages that the server sends back to the client and then we could generate
  the server handler and the client library in some form. I might need to write a protoc plugin for
  this!

### Migration

Importantly, I would like the migration to be gradual. Ideally I could serve any given page
side-by-side with the legacy Go/Templ/Datastar version and flip over via a query parameter. This way
I can make sure we have good parity before committing to the React version.
