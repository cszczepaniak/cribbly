package admincomponents

import (
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type Assignment struct {
	Assign   string
	Unassign string
}

func GetAssignmentForEdit(r *http.Request) Assignment {
	return Assignment{
		Assign:   r.FormValue("assign"),
		Unassign: r.FormValue("unassign"),
	}
}

func TableRowID(id string) string {
	return "table-row-" + id
}

func tableRowNameID(id string) string {
	return "table-row-name-" + id
}

func MainContentID[T teamOrDiv]() string {
	return ifTeam[T]("teams", "divisions")
}

const EditModalID = "edit-modal"

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

func editURL[T teamOrDiv](val T) string {
	switch t := any(val).(type) {
	case teams.Team:
		return "/admin/teams/edit/" + t.ID
	case divisions.Division:
		return "/admin/divisions/edit/" + t.ID
	default:
		panic("unreachable")
	}
}

func deleteURL[T teamOrDiv](val T) string {
	switch t := any(val).(type) {
	case teams.Team:
		return "/admin/teams/" + t.ID
	case divisions.Division:
		return "/admin/divisions/" + t.ID
	default:
		panic("unreachable")
	}
}

func assignURL[U playerOrTeam](teamOrDivID string, item U) string {
	switch t := any(item).(type) {
	case players.Player:
		return "/admin/teams/" + teamOrDivID + "?assign=" + t.ID
	case teams.Team:
		return "/admin/divisions/" + teamOrDivID + "?assign=" + t.ID
	default:
		panic("unreachable")
	}

}

func unassignURL[U playerOrTeam](teamOrDivID string, item U) string {
	switch t := any(item).(type) {
	case players.Player:
		return "/admin/teams/" + teamOrDivID + "?unassign=" + t.ID
	case teams.Team:
		return "/admin/divisions/" + teamOrDivID + "?unassign=" + t.ID
	default:
		panic("unreachable")
	}
}

func saveURL[T teamOrDiv](val T) string {
	switch t := any(val).(type) {
	case teams.Team:
		return "/admin/teams/" + t.ID
	case divisions.Division:
		return "/admin/divisions/" + t.ID
	default:
		panic("unreachable")
	}
}

func itemName[U playerOrTeam](item U) string {
	switch t := any(item).(type) {
	case players.Player:
		return t.Name
	case teams.Team:
		return t.Name
	default:
		panic("unreachable")
	}
}

func itemOrDivID[T teamOrDiv](item T) string {
	switch t := any(item).(type) {
	case teams.Team:
		return t.ID
	case divisions.Division:
		return t.ID
	default:
		panic("unreachable")
	}
}

func availableTitle[U playerOrTeam]() string {
	switch any(*new(U)).(type) {
	case players.Player:
		return "Available Players"
	case teams.Team:
		return "Available Teams"
	default:
		panic("unreachable")
	}
}
