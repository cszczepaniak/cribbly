package divisions

import (
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type Handler struct {
	DivisionService divisions.Service
	TeamService     teams.Service
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	ds, err := h.DivisionService.GetAll(r.Context())
	if err != nil {
		return err
	}

	return Index(ds).Render(r.Context(), w)
}

func (h Handler) GetDivisions(w http.ResponseWriter, r *http.Request) error {
	divisionID := r.PathValue("id")
	d, err := h.DivisionService.Get(r.Context(), divisionID)
	if err != nil {
		return err
	}

	ts, err := h.TeamService.GetForDivision(r.Context(), divisionID)
	if err != nil {
		return err
	}

	return Division(d, ts).Render(r.Context(), w)
}
