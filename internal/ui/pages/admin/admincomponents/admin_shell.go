package admincomponents

import (
	"fmt"

	"github.com/a-h/templ"
)

type Route string

const (
	Home      Route = ""
	Players   Route = "Players"
	Teams     Route = "Teams"
	Divisions Route = "Divisions"
	Games     Route = "Games"
	Users     Route = "Users"
)

func (r Route) ToSafeURL() templ.SafeURL {
	switch r {
	case Divisions:
		return "/admin/divisions"
	case Home:
		return "/admin"
	case Players:
		return "/admin/players"
	case Teams:
		return "/admin/teams"
	case Games:
		return "/admin/games"
	case Users:
		return "/admin/users"
	default:
		panic(fmt.Sprintf("unexpected admincomponents.Route: %#v", r))
	}
}
