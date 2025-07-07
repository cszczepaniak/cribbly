package admincomponents

import (
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type teamOrDiv interface {
	teams.Team | divisions.Division
}

type playerOrTeam interface {
	players.Player | teams.Team
}

type EditTeamOrDivisionModalProps[T teamOrDiv, U playerOrTeam] struct {
	val                  T
	itemsInThisTeamOrDiv []U
	availableItems       []U
}

func (ps EditTeamOrDivisionModalProps[T, U]) isTeam() bool {
	_, ok := any(ps.val).(teams.Team)
	return ok
}

func (ps EditTeamOrDivisionModalProps[T, U]) editTitle() string {
	if ps.isTeam() {
		return "Edit Team"
	}
	return "Edit Division"
}

func (ps EditTeamOrDivisionModalProps[T, U]) nameTitle() string {
	if ps.isTeam() {
		return "Team Name"
	}
	return "Division Name"
}

func (ps EditTeamOrDivisionModalProps[T, U]) listTitle() string {
	if ps.isTeam() {
		return "Teams:"
	}
	return "Divisions:"
}
func (ps EditTeamOrDivisionModalProps[T, U]) itemID(item U) string {
	switch t := any(item).(type) {
	case players.Player:
		return t.ID
	case teams.Team:
		return t.ID
	default:
		panic("unreachable")
	}
}

func (ps EditTeamOrDivisionModalProps[T, U]) assignURL(item U) string {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return "/admin/teams/" + t.ID + "?assignPlayer=" + ps.itemID(item)
	case divisions.Division:
		return "/admin/divisions/" + t.ID + "?assignTeam=" + ps.itemID(item)
	default:
		panic("unreachable")
	}

}

func (ps EditTeamOrDivisionModalProps[T, U]) unassignURL(item U) string {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return "/admin/teams/" + t.ID + "?unassignPlayer=" + ps.itemID(item)
	case divisions.Division:
		return "/admin/divisions/" + t.ID + "?unassignTeam=" + ps.itemID(item)
	default:
		panic("unreachable")
	}
}

func (ps EditTeamOrDivisionModalProps[T, U]) saveURL() string {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return "/admin/teams/" + t.ID
	case divisions.Division:
		return "/admin/divisions/" + t.ID
	default:
		panic("unreachable")
	}
}

func (ps EditTeamOrDivisionModalProps[T, U]) hxTarget() string {
	if ps.isTeam() {
		return "#teams"
	}
	return "#divisions"
}

func (ps EditTeamOrDivisionModalProps[T, U]) name() string {
	switch t := any(ps.val).(type) {
	case teams.Team:
		return t.Name
	case divisions.Division:
		return t.Name
	default:
		panic("unreachable")
	}
}

func (ps EditTeamOrDivisionModalProps[T, U]) itemName(item U) string {
	switch t := any(item).(type) {
	case players.Player:
		return t.Name
	case teams.Team:
		return t.Name
	default:
		panic("unreachable")
	}
}
