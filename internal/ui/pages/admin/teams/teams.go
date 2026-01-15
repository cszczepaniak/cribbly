package teams

import (
	"net/http"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
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
	err = sse.PatchElementTempl(editNameInput(team.Name))
	if err != nil {
		return err
	}
	err = sse.PatchElementTempl(editSaveButton(id))
	if err != nil {
		return err
	}
	return sse.PatchElementTempl(teamListing(id, availablePlayers, onThisTeam))
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
	return sse.PatchElementTempl(teamTable(teams))
}

func (h TeamsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	err := h.TeamService.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.RemoveElementByID(tableRowID(id))
}

type signals struct {
	Name string `json:"name"`
}

func (h TeamsHandler) Save(w http.ResponseWriter, r *http.Request) error {
	teamID := r.PathValue("id")

	assign := r.FormValue("assign")
	unassign := r.FormValue("unassign")
	if assign != "" || unassign != "" {
		var err error
		if assign != "" {
			err = h.PlayerService.AssignToTeam(r.Context(), assign, teamID)
		} else {
			err = h.PlayerService.UnassignFromTeam(r.Context(), unassign)
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
		return sse.PatchElementTempl(teamListing(teamID, available, onThisTeam))
	}

	var sigs signals
	err := datastar.ReadSignals(r, &sigs)
	if err != nil {
		return err
	}

	if sigs.Name != "" {
		err := h.TeamService.Rename(r.Context(), teamID, sigs.Name)
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(w, r)

		err = sse.PatchElementTempl(teamName(teamID, sigs.Name))
		if err != nil {
			return err
		}
		return sse.MarshalAndPatchSignals(signals{Name: ""})
	}

	return nil
}

func tableRowID(teamID string) string {
	return "table-row-" + teamID
}
