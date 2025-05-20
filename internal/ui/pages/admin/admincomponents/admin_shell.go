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
	default:
		panic(fmt.Sprintf("unexpected admincomponents.Route: %#v", r))
	}
}
