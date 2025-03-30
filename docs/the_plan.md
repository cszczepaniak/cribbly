# The Plan

## Minimum Viable Product
There are lots of things we could do with Cribbly to make it awesome. To begin with, let's outline
the minimum set of things it would need to do to be usable during a cribbage tournament.

As of this writing, this document is living and nothing here is final.

### Requirements

#### 0. Users

There are two primary kinds of users for Cribbly: admins and players.

Admins should be required to be authenticated somehow to perform admin duties. Players should not be
required to do so for simplicity.

#### 1. Player Registration

An admin should be able to register players for the tournament by name.

#### 2. Team Formation

An admin should be able to pair players together to form teams.

#### 3. Division Formation

An admin should be able to manually group teams together to form divisions.

Questions:
- Are divisions always 4 teams? Do we need to support 5- or 6-team divisions?

#### 4. Prelim Scheduling

The application should automatically schedule 3 preliminary games for each division such that each
team plays each other team once.

Any player should be able to view the schedule so they know who to play next.

Questions:
- If 5- or 6-team divisions are supported, how do we handle this?

#### 5. Prelim Score Reporting

A player on a team should be able to enter the score of their prelim game once it's over.

Any player should be able to view the current standings of any given division while games are
in-progress.

#### 6. Tournament Seeding

Based on prelim game scores, the application should automatically select 16 teams and seed them in a
single-elimination tournament.

#### 7. Tournament Win Reporting

A tournament player should be able to enter whether their team won or lost a tournament game, which
the application should use to advance the winning team.

#### 8. Tournament View

Any player should be able to view the state of the tournament bracket. This will probably be
projected on a screen as well.

#### 9. Admin Escape Hatch

Admins should have some set of tools to allow manually setting any of the above describe state in
case things go wrong and something needs to be corrected.

### Technologies

#### Language

Server options:
- Go: this is pretty much already decided as we're both comfortable with Go. Are there others we
should consider?

Assuming Go, more options:
  - Router options:
    - stdlib: it now supports [path parameters and method specification](https://pkg.go.dev/net/http#hdr-Patterns-ServeMux) for routes. Would this suffice, or do we need to reach for a third party library?
    - [fiber](https://gofiber.io/)

Frontend options:
- Markup options
  - HTMX from the server
    - [Go's builtin template package](https://pkg.go.dev/html/template)
    - [Templ](https://templ.guide/)
  - Something javascript-y like [Svelte](https://svelte.dev/) or [React](https://react.dev/)
- Styling: Tailwind (decided) ✅

#### Infrastructure
We'd like to keep costs as low as possible while in development.

Server options:
  - [Digital Ocean](https://www.digitalocean.com/)
  - [Railway](https://railway.com/)
  - [fly.io](https://fly.io/)
  - Some sort of cloud function option like AWS Lambda

Persistence options:
  - SQLite, [Turso](https://turso.tech/)) as an initial deployment option ✅

#### Other Stuff

Authentication:
  - OAuth2 with an identity provider like Google
    - I don't really know what this takes, but should be possible
  - Identity as a service, like [Auth0](https://auth0.com/)
    - Might have more bells and whistles than we need
