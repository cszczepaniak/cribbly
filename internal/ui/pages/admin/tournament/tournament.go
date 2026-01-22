package tournament

import (
	"net/http"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/games"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

type Handler struct {
	DivisionRepo divisions.Repository
	TeamRepo     teams.Repository
	GameRepo     games.Repository
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	return index().Render(r.Context(), w)
}
