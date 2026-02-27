package divisions

import (
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

func newDivisionService(t *testing.T) (Service, divisions.Repository, teams.Repository) {
	t.Helper()

	db := database.NewInMemory(t)
	txer := database.NewTransactor(db)
	dr := divisions.NewRepository(db)
	tr := teams.NewRepository(db)

	assert.NoError(t, dr.Init(t.Context()))
	assert.NoError(t, tr.Init(t.Context()))

	return New(txer, tr, dr), dr, tr
}

func TestGetDivisionAggregatesTeams(t *testing.T) {
	svc, dr, tr := newDivisionService(t)
	ctx := t.Context()

	div, err := dr.Create(ctx)
	assert.NoError(t, err)

	team1, err := tr.Create(ctx, "team1")
	assert.NoError(t, err)
	team2, err := tr.Create(ctx, "team2")
	assert.NoError(t, err)

	assert.NoError(t, tr.AssignToDivision(ctx, team1.ID, div.ID))
	assert.NoError(t, tr.AssignToDivision(ctx, team2.ID, div.ID))

	got, err := svc.Get(ctx, div.ID)
	assert.NoError(t, err)

	assert.Equal(t, div.ID, got.ID)
	assert.Equal(t, div.Name, got.Name)
	assert.Equal(t, div.Size, got.Size)
	assert.SliceLen(t, got.Teams, 2)
}
