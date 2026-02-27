package teams

import (
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

func newService(t *testing.T) (Service, players.Repository, teams.Repository) {
	t.Helper()

	db := database.NewInMemory(t)
	txer := database.NewTransactor(db)
	pr := players.NewRepository(db)
	tr := teams.NewRepository(db)

	assert.NoError(t, pr.Init(t.Context()))
	assert.NoError(t, tr.Init(t.Context()))

	return New(txer, pr, tr), pr, tr
}

func TestGetTeamAggregatesPlayers(t *testing.T) {
	svc, pr, tr := newService(t)
	ctx := t.Context()

	team, err := tr.Create(ctx, "team")
	assert.NoError(t, err)

	p1ID, err := pr.Create(ctx, "Alice", "Alpha")
	assert.NoError(t, err)
	p2ID, err := pr.Create(ctx, "Bob", "Beta")
	assert.NoError(t, err)

	assert.NoError(t, pr.AssignToTeam(ctx, p1ID, team.ID))
	assert.NoError(t, pr.AssignToTeam(ctx, p2ID, team.ID))

	got, err := svc.GetTeam(ctx, team.ID)
	assert.NoError(t, err)

	assert.Equal(t, team.ID, got.ID)
	assert.Equal(t, team.Name, got.Name)
	assert.SliceLen(t, got.Players, 2)
}

func TestCreateTeamCreatesUnnamedTeam(t *testing.T) {
	svc, _, tr := newService(t)
	ctx := t.Context()

	team, err := svc.CreateTeam(ctx)
	assert.NoError(t, err)

	got, err := tr.Get(ctx, team.ID)
	assert.NoError(t, err)

	assert.Equal(t, "Unnamed Team", got.Name)
}

func TestAssignPlayerToTeamUpdatesNameAndPlayers(t *testing.T) {
	svc, pr, tr := newService(t)
	ctx := t.Context()

	team, err := tr.Create(ctx, "team")
	assert.NoError(t, err)

	p1ID, err := pr.Create(ctx, "Alice", "Alpha")
	assert.NoError(t, err)

	p2ID, err := pr.Create(ctx, "Bob", "Beta")
	assert.NoError(t, err)

	teamWithOne, err := svc.AssignPlayerToTeam(ctx, p1ID, team.ID)
	assert.NoError(t, err)
	assert.Equal(t, "A Alpha", teamWithOne.Name)
	assert.SliceLen(t, teamWithOne.Players, 1)

	teamWithTwo, err := svc.AssignPlayerToTeam(ctx, p2ID, team.ID)
	assert.NoError(t, err)
	assert.Equal(t, "A Alpha / B Beta", teamWithTwo.Name)
	assert.SliceLen(t, teamWithTwo.Players, 2)
}

func TestUnassignPlayerFromTeamUpdatesName(t *testing.T) {
	svc, pr, tr := newService(t)
	ctx := t.Context()

	team, err := tr.Create(ctx, "team")
	assert.NoError(t, err)

	p1ID, err := pr.Create(ctx, "Alice", "Alpha")
	assert.NoError(t, err)
	p2ID, err := pr.Create(ctx, "Bob", "Beta")
	assert.NoError(t, err)

	_, err = svc.AssignPlayerToTeam(ctx, p1ID, team.ID)
	assert.NoError(t, err)
	_, err = svc.AssignPlayerToTeam(ctx, p2ID, team.ID)
	assert.NoError(t, err)

	teamAfterUnassign, err := svc.UnassignPlayerFromTeam(ctx, p2ID, team.ID)
	assert.NoError(t, err)
	assert.Equal(t, "A Alpha", teamAfterUnassign.Name)
	assert.SliceLen(t, teamAfterUnassign.Players, 1)

	teamAfterAllUnassigned, err := svc.UnassignPlayerFromTeam(ctx, p1ID, team.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Unnamed Team", teamAfterAllUnassigned.Name)
	assert.SliceLen(t, teamAfterAllUnassigned.Players, 0)
}
