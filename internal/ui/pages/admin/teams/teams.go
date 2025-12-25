package teams

import (
	"net/http"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	"github.com/cszczepaniak/cribbly/internal/ui/pages/admin/admincomponents"
)

type TeamsHandler struct {
	PlayerService players.Service
	TeamService   teams.Service
}

func (h TeamsHandler) Index(w http.ResponseWriter, r *http.Request) error {
	teams, err := h.TeamService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return index(teams).Render(r.Context(), w)
}

func (h TeamsHandler) Edit(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	team, err := h.TeamService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	availablePlayers, err := h.PlayerService.GetFreeAgents(r.Context())
	if err != nil {
		return err
	}

	onThisTeam, err := h.PlayerService.GetForTeam(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(
		admincomponents.EditTeamOrDivisionModal(team, onThisTeam, availablePlayers),
		datastar.WithSelectorID(admincomponents.MainContentID[teams.Team]()),
		datastar.WithModeAppend(),
	)
}

func (h TeamsHandler) CancelEdit(w http.ResponseWriter, r *http.Request) error {
	sse := datastar.NewSSE(w, r)
	return resetEdit(sse)
}

func (h TeamsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	_, err := h.TeamService.Create(r.Context())
	if err != nil {
		return err
	}

	teams, err := h.TeamService.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(admincomponents.TeamOrDivisionTable(teams))
}

func (h TeamsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	err := h.TeamService.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.RemoveElementByID(admincomponents.TableRowID(id))
}

type signals struct {
	Name string `json:"name"`
}

func (h TeamsHandler) Save(w http.ResponseWriter, r *http.Request) error {
	teamID := r.PathValue("id")

	a := admincomponents.GetAssignmentForEdit(r)
	if a != (admincomponents.Assignment{}) {
		var err error
		if a.Assign != "" {
			err = h.PlayerService.AssignToTeam(r.Context(), a.Assign, teamID)
		} else {
			err = h.PlayerService.UnassignFromTeam(r.Context(), a.Unassign)
		}
		if err != nil {
			return err
		}

		onThisTeam, err := h.PlayerService.GetForTeam(r.Context(), teamID)
		if err != nil {
			return err
		}

		available, err := h.PlayerService.GetFreeAgents(r.Context())
		if err != nil {
			return err
		}

		// If we're unassigning a team, we'll keep the modal open (by not redirecting).
		sse := datastar.NewSSE(w, r)
		return sse.PatchElementTempl(admincomponents.ItemsListing[teams.Team](teamID, available, onThisTeam))
	}

	var signals signals
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	if signals.Name != "" {
		err := h.TeamService.Rename(r.Context(), teamID, signals.Name)
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(w, r)

		err = sse.PatchElementTempl(admincomponents.TeamOrDivisionName(teamID, signals.Name))
		if err != nil {
			return err
		}
		return resetEdit(sse)
	}

	return nil
}

func resetEdit(sse *datastar.ServerSentEventGenerator) error {
	err := sse.MarshalAndPatchSignals(signals{Name: ""})
	if err != nil {
		return err
	}
	return sse.RemoveElementByID(admincomponents.EditModalID)
}
