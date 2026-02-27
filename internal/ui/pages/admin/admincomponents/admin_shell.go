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
	RoomCodes Route = "Room Codes"
	Users     Route = "Users"
	Profile   Route = "My Profile"
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
	case RoomCodes:
		return "/admin/room-codes"
	case Users:
		return "/admin/users"
	case Profile:
		return "/admin/profile"
	default:
		panic(fmt.Sprintf("unexpected admincomponents.Route: %#v", r))
	}
}
