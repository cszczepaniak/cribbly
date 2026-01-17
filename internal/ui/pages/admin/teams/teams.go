package teams

import (
	"net/http"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/moreiter"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	teamservice "github.com/cszczepaniak/cribbly/internal/service/teams"
)

type TeamsHandler struct {
	PlayerRepo  players.Repository
	TeamRepo    teams.Repository
	TeamService teamservice.Service
}

func (h TeamsHandler) Index(w http.ResponseWriter, r *http.Request) error {
	teams, err := h.TeamRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	playersByTeam := make(map[string][]players.Player, len(teams))
	for _, team := range teams {
		players, err := h.PlayerRepo.GetForTeam(r.Context(), team.ID)
		if err != nil {
			return err
		}
		playersByTeam[team.ID] = players
	}

	return index(teams, playersByTeam).Render(r.Context(), w)
}

func (h TeamsHandler) Edit(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	team, err := h.TeamRepo.Get(r.Context(), id)
	if err != nil {
		return err
	}

	availablePlayers, err := h.PlayerRepo.GetFreeAgents(r.Context())
	if err != nil {
		return err
	}

	onThisTeam, err := h.PlayerRepo.GetForTeam(r.Context(), id)
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

func (h TeamsHandler) ConfirmDelete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	team, err := h.TeamRepo.Get(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	err = sse.PatchElementTempl(confirmDeleteTitle(team.Name))
	if err != nil {
		return err
	}

	return sse.PatchElementTempl(confirmDeleteButton(team.ID))
}

func (h TeamsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	_, err := h.TeamRepo.Create(r.Context(), "Unnamed Team")
	if err != nil {
		return err
	}

	teams, err := h.TeamRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	playersByTeam := make(map[string][]players.Player, len(teams))
	for _, team := range teams {
		players, err := h.PlayerRepo.GetForTeam(r.Context(), team.ID)
		if err != nil {
			return err
		}
		playersByTeam[team.ID] = players
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(teamGrid(teams, playersByTeam))
}

func (h TeamsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	err := h.TeamService.DeleteTeam(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.RemoveElementByID(teamItemID(id))
}

func (h TeamsHandler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	players, err := h.PlayerRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	for _, p := range players {
		if p.TeamID != "" {
			err := h.PlayerRepo.UnassignFromTeam(r.Context(), p.ID, moreiter.Of(p.TeamID))
			if err != nil {
				return err
			}
		}
	}

	err = h.TeamRepo.DeleteAll(r.Context())
	if err != nil {
		return err
	}

	return datastar.NewSSE(w, r).Redirect("/admin/teams")
}

func (h TeamsHandler) Generate(w http.ResponseWriter, r *http.Request) error {
	players, err := h.PlayerRepo.GetFreeAgents(r.Context())
	if err != nil {
		return err
	}

	for len(players) >= 2 {
		team, err := h.TeamService.CreateTeam(r.Context())
		if err != nil {
			return err
		}

		for i := range 2 {
			_, err = h.TeamService.AssignPlayerToTeam(r.Context(), players[i].ID, team.ID)
			if err != nil {
				return err
			}
		}

		players = players[2:]
	}

	return datastar.NewSSE(w, r).Redirect("/admin/teams")
}

func (h TeamsHandler) Save(w http.ResponseWriter, r *http.Request) error {
	teamID := r.PathValue("id")

	assign := r.FormValue("assign")
	unassign := r.FormValue("unassign")
	if assign != "" || unassign != "" {
		var team teamservice.Team
		var err error
		if assign != "" {
			team, err = h.TeamService.AssignPlayerToTeam(r.Context(), assign, teamID)
		} else {
			team, err = h.TeamService.UnassignPlayerFromTeam(r.Context(), unassign, teamID)
		}
		if err != nil {
			return err
		}

		onThisTeam, err := h.PlayerRepo.GetForTeam(r.Context(), teamID)
		if err != nil {
			return err
		}

		available, err := h.PlayerRepo.GetFreeAgents(r.Context())
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(w, r)
		err = sse.PatchElementTempl(teamName(teamID, team.Name))
		if err != nil {
			return err
		}
		err = sse.PatchElementTempl(teamPlayersList(teamID, onThisTeam))
		if err != nil {
			return err
		}
		return sse.PatchElementTempl(teamListing(teamID, available, onThisTeam))
	}

	var sigs struct {
		Name string `json:"name"`
	}

	err := datastar.ReadSignals(r, &sigs)
	if err != nil {
		return err
	}

	if sigs.Name != "" {
		err := h.TeamRepo.Rename(r.Context(), teamID, sigs.Name)
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(w, r)

		err = sse.PatchElementTempl(teamName(teamID, sigs.Name))
		if err != nil {
			return err
		}
		sigs.Name = ""
		return sse.MarshalAndPatchSignals(sigs)
	}

	return nil
}

func teamItemID(teamID string) string {
	return "team-" + teamID
}
