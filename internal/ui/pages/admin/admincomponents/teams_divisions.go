package admincomponents

import (
	"github.com/a-h/templ"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type teamOrDiv interface {
	teams.Team | divisions.Division
}

func getID[T teamOrDiv](val T) string {
	switch t := any(val).(type) {
	case teams.Team:
		return t.ID
	case divisions.Division:
		return t.ID
	default:
		panic("unreachable")
	}
}

func getName[T teamOrDiv](val T) string {
	switch t := any(val).(type) {
	case teams.Team:
		return t.Name
	case divisions.Division:
		return t.Name
	default:
		panic("unreachable")
	}
}

func getHXTarget[T teamOrDiv]() string {
	return ifTeam[T]("#teams", "#divisions")
}

func isTeam[T teamOrDiv]() bool {
	_, ok := any(*new(T)).(teams.Team)
	return ok
}

func ifTeam[T teamOrDiv, U any](yes, el U) U {
	if isTeam[T]() {
		return yes
	}
	return el
}

type playerOrTeam interface {
	players.Player | teams.Team
}

type teamOrDivisionRowProps[T teamOrDiv] struct {
	val T
}

func (ps teamOrDivisionRowProps[T]) editURL() templ.SafeURL {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return templ.SafeURL("/admin/teams?edit=" + t.ID)
	case divisions.Division:
		return templ.SafeURL("/admin/divisions?edit=" + t.ID)
	default:
		panic("unreachable")
	}
}

func (ps teamOrDivisionRowProps[T]) deleteURL() string {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return "/admin/teams/" + t.ID
	case divisions.Division:
		return "/admin/divisions/" + t.ID
	default:
		panic("unreachable")
	}
}

type editTeamOrDivisionModalProps[T teamOrDiv, U playerOrTeam] struct {
	val                  T
	itemsInThisTeamOrDiv []U
	availableItems       []U
}

func (ps editTeamOrDivisionModalProps[T, U]) isTeam() bool {
	_, ok := any(ps.val).(teams.Team)
	return ok
}

func (ps editTeamOrDivisionModalProps[T, U]) itemID(item U) string {
	switch t := any(item).(type) {
	case players.Player:
		return t.ID
	case teams.Team:
		return t.ID
	default:
		panic("unreachable")
	}
}

func (ps editTeamOrDivisionModalProps[T, U]) assignURL(item U) string {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return "/admin/teams/" + t.ID + "?assignPlayer=" + ps.itemID(item)
	case divisions.Division:
		return "/admin/divisions/" + t.ID + "?assignTeam=" + ps.itemID(item)
	default:
		panic("unreachable")
	}

}

func (ps editTeamOrDivisionModalProps[T, U]) unassignURL(item U) string {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return "/admin/teams/" + t.ID + "?unassignPlayer=" + ps.itemID(item)
	case divisions.Division:
		return "/admin/divisions/" + t.ID + "?unassignTeam=" + ps.itemID(item)
	default:
		panic("unreachable")
	}
}

func (ps editTeamOrDivisionModalProps[T, U]) saveURL() string {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return "/admin/teams/" + t.ID
	case divisions.Division:
		return "/admin/divisions/" + t.ID
	default:
		panic("unreachable")
	}
}

func (ps editTeamOrDivisionModalProps[T, U]) name() string {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return t.Name
	case divisions.Division:
		return t.Name
	default:
		panic("unreachable")
	}
}

func (ps editTeamOrDivisionModalProps[T, U]) itemName(item U) string {
	switch t := any(item).(type) {
	case players.Player:
		return t.Name
	case teams.Team:
		return t.Name
	default:
		panic("unreachable")
	}
}

func (ps editTeamOrDivisionModalProps[T, U]) availableTitle() string {
	switch any(*new(U)).(type) {
	case players.Player:
		return "Available Players:"
	case teams.Team:
		return "Available Teams:"
	default:
		panic("unreachable")
	}
}
