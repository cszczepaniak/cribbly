package divisions

import (
	"cmp"
	"net/http"
	"slices"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type Handler struct {
	DivisionRepo divisions.Repository
	TeamRepo     teams.Repository
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	ds, err := h.DivisionRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	slices.SortFunc(ds, func(a, b divisions.Division) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return Index(ds).Render(r.Context(), w)
}

func (h Handler) GetDivisions(w http.ResponseWriter, r *http.Request) error {
	divisionID := r.PathValue("id")
	d, err := h.DivisionRepo.Get(r.Context(), divisionID)
	if err != nil {
		return err
	}

	ts, err := h.TeamRepo.GetForDivision(r.Context(), divisionID)
	if err != nil {
		return err
	}

	slices.SortFunc(ts, func(a, b teams.Team) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return Division(d, ts).Render(r.Context(), w)
}
